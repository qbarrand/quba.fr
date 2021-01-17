package assets

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/sirupsen/logrus"

	"github.com/qbarrand/quba.fr/pkg/fileutils"
)

var (
	ErrPathTranslationConflict = errors.New("this original path was already registered")
	ErrTemplateError           = errors.New("could not compile template")
)

type Server struct {
	fs          FS
	logger      logrus.FieldLogger
	hashedPaths map[string]string
}

func NewServer(logger logrus.FieldLogger) *Server {
	return &Server{
		fs:          make(FS),
		logger:      logger,
		hashedPaths: make(map[string]string),
	}
}

func (s *Server) AddDirectory(path string) error {
	return s.fs.AddFile(path, directoryFileGetter(path))
}

func (s *Server) AddStaticFile(filePath string, resourcePath string, hashPath bool) (string, error) {
	if s.hashedPaths[resourcePath] != "" {
		return "", ErrPathTranslationConflict
	}

	fd, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("could not open %q: %v", filePath, err)
	}
	defer fd.Close()

	h, err := fileutils.HashReader(fd)
	if err != nil {
		return "", fmt.Errorf("could not hash %q: %v", filePath, err)
	}

	newPath := resourcePath

	if hashPath {
		newPath = fileutils.PathWithHash(resourcePath, h)
	}

	s.hashedPaths[resourcePath] = newPath

	return newPath, s.fs.AddFile(newPath, staticFileGetter(filePath))
}

func (s *Server) AddTemplate(text, resourcePath string, hashPath bool) (string, error) {
	if s.hashedPaths[resourcePath] != "" {
		return "", ErrPathTranslationConflict
	}

	tpl := template.
		New(resourcePath).
		Option("missingkey=error").
		Funcs(template.FuncMap{
			"getDependency": func(path string) (string, error) {
				hashedPath := s.hashedPaths[path]

				if hashedPath == "" {
					return "", errors.New("dependency not found")
				}

				return hashedPath, nil
			},
		})

	if _, err := tpl.Parse(text); err != nil {
		return "", fmt.Errorf("could not parse the template: %v", err)
	}

	var buf bytes.Buffer

	if err := tpl.Execute(&buf, nil); err != nil {
		return "", fmt.Errorf("%w: %v", ErrTemplateError, err)
	}

	b := buf.Bytes()

	h, err := fileutils.HashReader(bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("could not hash the template output: %v", err)
	}

	newPath := resourcePath

	if hashPath {
		newPath = fileutils.PathWithHash(resourcePath, h)
	}

	s.hashedPaths[resourcePath] = newPath

	tfg := templateFileGetter{
		buf:  b,
		name: newPath,
	}

	return newPath, s.fs.AddFile(newPath, &tfg)
}

func (s *Server) Handler() http.Handler {
	return http.FileServer(s.fs)
}
