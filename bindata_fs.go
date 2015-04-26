package main

import (
	"bytes"
	"net/http"
	"os"
)

type BindataFileSystem struct{}

func (f BindataFileSystem) Open(name string) (http.File, error) {
	path := "dist" + name

	data, err := Asset(path)
	if err != nil {
		return nil, err
	}

	return &BindataFile{bytes.NewReader(data), path}, nil
}

type BindataFile struct {
	*bytes.Reader
	path string
}

func (f *BindataFile) Close() error {
	return nil
}

func (f *BindataFile) Readdir(count int) ([]os.FileInfo, error) {
	return make([]os.FileInfo, 0), nil
}

func (f *BindataFile) Stat() (os.FileInfo, error) {
	return AssetInfo(f.path)
}
