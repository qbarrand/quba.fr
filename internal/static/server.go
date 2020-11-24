package static

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
)

var (
	ErrElementPathConflict     = errors.New("an element is already registered for this key")
	ErrNotFound                = errors.New("no dependency registered under this path")
	ErrPathTranslationConflict = errors.New("this original path was already registered")
	ErrTemplateError           = errors.New("could not compile template")
)

type Server struct {
	elements map[string]http.Handler
	logger   logrus.FieldLogger
	paths    map[string]string
}

func NewServer(logger logrus.FieldLogger) *Server {
	return &Server{
		elements: make(map[string]http.Handler),
		logger:   logger,
		paths:    make(map[string]string),
	}
}

func (s *Server) AddTemplate(text, resourcePath string, hashPath bool) error {
	if s.paths[resourcePath] != "" {
		return ErrPathTranslationConflict
	}

	tpl := template.
		New(resourcePath).
		Option("missingkey=error").
		Funcs(template.FuncMap{
			"getDependency": func(path string) (string, error) {
				hashedPath := s.paths[path]

				if hashedPath == "" {
					return "", errors.New("dependency not found")
				}

				return hashedPath, nil
			},
		})

	if _, err := tpl.Parse(text); err != nil {
		return fmt.Errorf("could not parse the template: %v", err)
	}

	var buf bytes.Buffer

	if err := tpl.Execute(&buf, nil); err != nil {
		return fmt.Errorf("%w: %v", ErrTemplateError, err)
	}

	b := buf.Bytes()

	h, err := hashReader(bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("could not hash the template output: %v", err)
	}

	newPath := resourcePath

	if hashPath {
		newPath = pathWithHash(resourcePath, h)
	}

	s.paths[resourcePath] = newPath

	if s.elements[newPath] != nil {
		return ErrElementPathConflict
	}

	s.elements[newPath] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", h)
		w.Write(b)
	})

	return nil
}

func (s *Server) AddStaticFile(filePath string, resourcePath string, hashPath bool) error {
	if s.paths[resourcePath] != "" {
		return ErrPathTranslationConflict
	}

	fd, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open %q: %v", filePath, err)
	}
	defer fd.Close()

	h, err := hashReader(fd)
	if err != nil {
		return fmt.Errorf("could not hash %q: %v", filePath, err)
	}

	newPath := resourcePath

	if hashPath {
		newPath = pathWithHash(resourcePath, h)
	}

	s.paths[resourcePath] = newPath

	if s.elements[newPath] != nil {
		return ErrElementPathConflict
	}

	s.elements[newPath] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", h)
		http.ServeFile(w, r, filePath)
	})

	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler := s.elements[req.URL.Path]

	if handler == nil {
		http.NotFound(w, req)
		return
	}

	handler.ServeHTTP(w, req)
}

func hashReader(r io.Reader) (string, error) {
	hasher := fnv.New32()

	if _, err := io.Copy(hasher, r); err != nil {
		return "", fmt.Errorf("could copy bytes into the hasher: %v", err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func pathWithHash(path, hash string) string {
	ext := filepath.Ext(path)

	// 4 times faster than fmt.Sprintf
	var sb strings.Builder

	sb.WriteString(path[:len(path)-len(ext)])
	sb.WriteRune('.')
	sb.WriteString(hash)
	sb.WriteString(ext)

	return sb.String()
}
