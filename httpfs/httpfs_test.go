package httpfs

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func filter(name string) bool {
	if strings.HasPrefix(name, ".") {
		return false
	}

	return true
}

func TestCondNameFileSystem(t *testing.T) {
	cases := []struct {
		recu bool
		name string
		hit  bool
	}{
		{true, "dir", true},
		{true, ".hidden-dir", false},
		{true, ".hidden-dir/in-hidden.txt", false},
		{true, "dir/file.txt", true},
		{true, ".hidden-file.txt", false},
		{true, "file.txt", true},
		{true, "no-exist", false},
		{true, "", true},

		{false, "dir", true},
		{false, ".hidden-dir", false},
		{false, ".hidden-dir/in-hidden.txt", true},
		{false, "dir/file.txt", true},
		{false, ".hidden-file.txt", false},
		{false, "file.txt", true},
		{false, "no-exist", false},
		{false, "", true},
	}

	for _, c := range cases {
		fs := &CondNameFileSystem{
			Root: "testdata", Filter: filter, Recu: c.recu,
		}

		_, err := fs.Open(c.name)
		if c.hit {
			assert.Nil(t, err, "%v", c)
		} else {
			assert.Error(t, err, "%v", c)
		}
	}
}

func TestCondNameFileSystemReaddir(t *testing.T) {
	cases := []struct {
		name string
		rdns []string // f.Readdir()'s filename
	}{
		{"", []string{"dir", "file.txt"}},
		{"dir", []string{"file.txt"}},
	}

	for _, c := range cases {
		fs := &CondNameFileSystem{
			Root: "testdata", Filter: filter, Recu: false,
		}

		run := func(n int) {
			f, _ := fs.Open(c.name)

			var fis []os.FileInfo
			if n > 0 {
				for {
					buff, err := f.Readdir(n)
					fis = append(fis, buff...)
					if err == io.EOF {
						break
					}
				}
			} else {
				fis, _ = f.Readdir(n)
			}

			names := make([]string, 0, len(fis))
			for _, fi := range fis {
				names = append(names, fi.Name())
			}

			assert.ElementsMatch(t, c.rdns, names, "%v", c)
		}

		run(0)   // n <= 0
		run(1)   // n > 0
		run(100) // n > all element count
	}
}
