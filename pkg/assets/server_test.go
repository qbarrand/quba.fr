package assets

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	hashedResourcePath = "/path/to/resource." + testFileHash + ".txt"
	resourcePath       = "/path/to/resource.txt"
	testFileContents   = "text"
	testFileHash       = "b12bfa38"
)

func TestServer_AddTemplate(t *testing.T) {
	logger, _ := test.NewNullLogger()

	t.Run("original path already exists", func(t *testing.T) {
		const path = "/test/path"

		s := NewServer(logger)

		s.hashedPaths[path] = "/some/other/path"

		err := s.AddTemplate("", path, false)

		assert.True(t, errors.Is(err, ErrPathTranslationConflict))
	})

	t.Run("refers to a key that does not exist", func(t *testing.T) {
		s := NewServer(logger)

		err := s.AddTemplate("{{ .someKey }}", resourcePath, false)

		assert.Error(t, err)
	})

	t.Run("template refers to a dependency that does not exist", func(t *testing.T) {
		s := NewServer(logger)

		err := s.AddTemplate(`{{ getDependency "/some/dependency" }}`, resourcePath, false)

		assert.Error(t, err)
	})

	t.Run("works as expected", func(t *testing.T) {
		fd, err := ioutil.TempFile(t.TempDir(), "")

		require.NoError(t, err)

		defer fd.Close()

		_, err = fd.Write([]byte(testFileContents))

		require.NoError(t, err)

		const (
			depPath            = "/dependency/path.txt"
			resourceHash       = "7c0e4615"
			hashedResourcePath = "/path/to/resource." + resourceHash + ".txt"
		)

		cases := []struct {
			hashPath             bool
			expectedResourcePath string
		}{
			{
				hashPath:             false,
				expectedResourcePath: resourcePath,
			},
			{
				hashPath:             true,
				expectedResourcePath: hashedResourcePath,
			},
		}

		for _, c := range cases {
			t.Run(fmt.Sprintf("hashPath: %t", c.hashPath), func(t *testing.T) {
				s := NewServer(logger)

				s.hashedPaths[depPath] = depPath

				err = s.AddTemplate(
					fmt.Sprintf(`{{ getDependency "%s" }}`, depPath),
					resourcePath,
					c.hashPath)

				require.NoError(t, err)
				assert.Equal(t, c.expectedResourcePath, s.hashedPaths[resourcePath])

				req := httptest.NewRequest(http.MethodGet, resourcePath, nil)
				w := httptest.NewRecorder()

				handler := s.elements[c.expectedResourcePath]

				require.NotNil(t, handler)

				handler.ServeHTTP(w, req)

				assert.Equal(t, resourceHash, w.Header().Get("ETag"))
				assert.Equal(t, depPath, w.Body.String())
			})
		}
	})
}

func TestServer_AddHashedStaticFile(t *testing.T) {
	logger, _ := test.NewNullLogger()

	t.Run("original path already exists", func(t *testing.T) {
		const path = "/test/path"

		s := NewServer(logger)

		s.hashedPaths[path] = "/some/other/path"

		err := s.AddStaticFile("", path, false)

		assert.True(t, errors.Is(err, ErrPathTranslationConflict))
	})

	t.Run("file does not exists", func(t *testing.T) {
		s := NewServer(logger)

		err := s.AddStaticFile("/non/existent/file", "/some/path", false)

		assert.Error(t, err)
	})

	t.Run("works as expected", func(t *testing.T) {

		fd, err := ioutil.TempFile(t.TempDir(), "")
		defer fd.Close() // t.TempDir() will remove the directory and its contents

		require.NoError(t, err)

		_, err = fd.Write([]byte(testFileContents))

		require.NoError(t, err)

		cases := []struct {
			hashPath             bool
			expectedResourcePath string
		}{
			{
				hashPath:             false,
				expectedResourcePath: resourcePath,
			},
			{
				hashPath:             true,
				expectedResourcePath: hashedResourcePath,
			},
		}

		for _, c := range cases {
			t.Run(fmt.Sprintf("hashPath: %t", c.hashPath), func(t *testing.T) {
				s := NewServer(logger)

				err = s.AddStaticFile(fd.Name(), resourcePath, c.hashPath)

				require.NoError(t, err)
				assert.Equal(t, c.expectedResourcePath, s.hashedPaths[resourcePath])

				req := httptest.NewRequest(http.MethodGet, resourcePath, nil)
				w := httptest.NewRecorder()

				handler := s.elements[c.expectedResourcePath]

				require.NotNil(t, handler)

				handler.ServeHTTP(w, req)

				assert.Equal(t, testFileHash, w.Header().Get("ETag"))
				assert.Equal(t, testFileContents, w.Body.String())

			})
		}
	})
}

func TestServer_ServeHTTP(t *testing.T) {
	logger, _ := test.NewNullLogger()

	s := NewServer(logger)

	t.Run("not found", func(t *testing.T) {
		assert.HTTPStatusCode(t, s.ServeHTTP, http.MethodGet, "/some/resource", nil, http.StatusNotFound)
	})

	t.Run("works as expected", func(t *testing.T) {
		const status = http.StatusTeapot

		s.elements[hashedResourcePath] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
		})

		req := httptest.NewRequest(http.MethodGet, hashedResourcePath, nil)
		w := httptest.NewRecorder()

		s.ServeHTTP(w, req)

		assert.Equal(t, status, w.Result().StatusCode)
	})
}
