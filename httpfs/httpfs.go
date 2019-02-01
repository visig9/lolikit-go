// Package httpfs offer some http.FileSystem implement.
package httpfs

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/visig9/lolikit-go/pathrange"
)

// func keepWithoutPrefixDot(name string) bool {
// 	if !strings.HasPrefix(name, ".") {
// 		return false
// 	}
//
// 	return true
// }

// NameFilter is using to check the file name (only base name)
// should keep or not.
type NameFilter func(name string) bool

type dir struct {
	*os.File
	filter NameFilter
}

func (d dir) filtered(fis []os.FileInfo) []os.FileInfo {
	out := make([]os.FileInfo, 0, len(fis))
	for _, fi := range fis {
		if d.filter(fi.Name()) {
			out = append(out, fi)
		}
	}

	return out
}

func (d dir) Readdir(n int) ([]os.FileInfo, error) {
	if n <= 0 {
		buff, err := d.File.Readdir(n)
		return d.filtered(buff), err
	}

	var out []os.FileInfo
	for {
		buff, err := d.File.Readdir(n - len(out))
		out = append(out, d.filtered(buff)...)
		fmt.Println(err, len(out), n)

		if err != nil || len(out) == n {
			return out, err
		}
	}
}

// CondNameFileSystem is a http.FileSystem support filter some file by
// file's name.
type CondNameFileSystem struct {
	Root   string     // the root path in filesystem
	Filter NameFilter // which name pattern allow to expose
	Recu   bool       // should Filter run through the whole path
}

// FilterAlongPath return true if all portion of path return true,
// else return false.
func filterAlongPath(from, to string, filter NameFilter) bool {
	pr, err := pathrange.PathRange(from, to)
	if err != nil { // It's impossible when using correctly, panic
		panic(err)
	}

	for _, path := range pr {
		name := filepath.Base(path)
		if !filter(name) { // any portion of path not matching
			return false
		}
	}

	return true
}

// Open file by given file name.
func (fs *CondNameFileSystem) Open(name string) (http.File, error) {
	name = filepath.FromSlash(name)
	path := filepath.Join(fs.Root, name)

	switch fs.Recu {
	case true:
		if !filterAlongPath(fs.Root, path, fs.Filter) {
			return nil, os.ErrNotExist
		}
	case false:
		base := filepath.Base(name)
		if base != "." {
			if !fs.Filter(base) {
				return nil, os.ErrNotExist
			}
		}
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return dir{File: f, filter: fs.Filter}, nil
	}

	return f, nil
}
