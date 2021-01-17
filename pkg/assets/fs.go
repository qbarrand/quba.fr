package assets

import (
	"errors"
	"net/http"
	"os"
)

type FS map[string]fileGetter

func (fs FS) AddFile(path string, fg fileGetter) error {
	if fs[path] != nil {
		return errors.New("this path already exists")
	}

	fs[path] = fg

	return nil
}

func (fs FS) Open(path string) (http.File, error) {
	f := fs[path]
	if f == nil {
		return nil, os.ErrNotExist
	}

	return f.getFile()
}
