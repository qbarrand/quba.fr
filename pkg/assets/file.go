package assets

import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"time"
)

const defaultMode = 0444

var modTime = time.Now()

type fileGetter interface {
	getFile() (http.File, error)
}

type staticFileGetter string

func (sfg staticFileGetter) getFile() (http.File, error) {
	return os.Open(string(sfg))
}

type templateFileGetter struct {
	buf  []byte
	name string
}

func (tfg *templateFileGetter) getFile() (http.File, error) {
	bf := bufferFile{
		Reader: bytes.NewReader(tfg.buf),

		modTime: modTime,
		name:    tfg.name,
		size:    len(tfg.buf),
	}

	return &bf, nil
}

// bufferFile implements both http.File and os.FileInfo.
type bufferFile struct {
	*bytes.Reader

	modTime time.Time
	name    string
	size    int
}

func (bf *bufferFile) Close() error {
	return nil
}

func (bf *bufferFile) IsDir() bool {
	return false
}

func (bf *bufferFile) Mode() os.FileMode {
	// TODO
	panic("implement me")
}

func (bf *bufferFile) ModTime() time.Time {
	return bf.modTime
}

func (bf *bufferFile) Name() string {
	return bf.name
}

func (bf *bufferFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, errors.New("not implemented")
}

func (bf *bufferFile) Stat() (os.FileInfo, error) {
	return bf, nil
}

func (bf *bufferFile) Sys() interface{} {
	return nil
}

// directoryFileGetter implements fileGetter, http.File and os.FileInfo.
type directoryFileGetter string

func (dfg directoryFileGetter) Close() error {
	return nil // noop
}

func (dfg directoryFileGetter) IsDir() bool {
	return true
}

func (dfg directoryFileGetter) Mode() os.FileMode {
	return os.ModeDir | defaultMode
}

func (dfg directoryFileGetter) ModTime() time.Time {
	return modTime
}

func (dfg directoryFileGetter) Name() string {
	return string(dfg)
}

func (dfg directoryFileGetter) Read(p []byte) (n int, err error) {
	return 0, errors.New("cannot Read() a directory")
}

func (dfg directoryFileGetter) Seek(_ int64, _ int) (int64, error) {
	return 0, errors.New("cannot Seek() a directory")
}

func (dfg directoryFileGetter) Size() int64 {
	return 0
}

func (dfg directoryFileGetter) Stat() (os.FileInfo, error) {
	return dfg, nil
}

func (dfg directoryFileGetter) Sys() interface{} {
	return nil
}

func (dfg directoryFileGetter) Readdir(_ int) ([]os.FileInfo, error) {
	// Never allow listing directories
	return nil, nil
}

func (dfg directoryFileGetter) getFile() (http.File, error) {
	return dfg, nil
}
