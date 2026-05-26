package embed

import (
	"embed"
	"io"
	"io/fs"
	"path"
)

//go:embed files/*
var templateFS embed.FS

type embedFSLoader struct {
	fs   fs.FS
	root string
}

func (l *embedFSLoader) Exists(name string) bool {
	_, err := fs.Stat(l.fs, path.Join(l.root, name))
	return err == nil
}

func (l *embedFSLoader) Open(name string) (io.ReadCloser, error) {
	f, err := l.fs.Open(path.Join(l.root, name))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func NewLoader() *embedFSLoader {
	return &embedFSLoader{
		fs:   templateFS,
		root: "files",
	}
}
