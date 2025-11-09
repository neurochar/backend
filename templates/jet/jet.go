// Package jetset contains jet templates
package jetset

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"path"

	"github.com/CloudyKit/jet/v6"
	"github.com/govalues/decimal"
)

//go:embed files/**/*.jet
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

// Tpls - jet templates
var Tpls = jet.NewSet(
	&embedFSLoader{fs: templateFS, root: "files"},
)

func init() {
	Tpls.AddGlobal("sprintf", func(format string, v ...any) string {
		return fmt.Sprintf(format, v...)
	})

	Tpls.AddGlobal("decformat", func(d decimal.Decimal, digits int) string {
		s := d.Trim(digits).String()
		return s
	})
}
