package serve

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/visig/lolikit-go/fstest"
	"gitlab.com/visig/lolikit-go/loli2"
)

type noPeriodDir struct {
	*os.File
}

func filteredPeriod(fis []os.FileInfo) []os.FileInfo {
	out := make([]os.FileInfo, 0, len(fis))
	for _, fi := range fis {
		if !strings.HasPrefix(fi.Name(), ".") {
			out = append(out, fi)
		}
	}

	return out
}

func (npd noPeriodDir) Readdir(n int) ([]os.FileInfo, error) {
	if n <= 0 {
		buff, err := npd.File.Readdir(n)
		return filteredPeriod(buff), err
	}

	var out []os.FileInfo
	for {
		buff, err := npd.File.Readdir(len(out) - n)
		out = append(out, filteredPeriod(buff)...)

		if err != nil || len(out) == n {
			return out, err
		}
	}
}

type repoFileSystem struct {
	repo loli2.Repo
}

func (r *repoFileSystem) Open(name string) (http.File, error) {
	pfpath := filepath.FromSlash(name)
	path := filepath.Join(r.repo.Path(), pfpath)

	if fstest.IsDir(path) {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		return noPeriodDir{File: f}, err
	}

	visible, err := isVisible(r.repo.Path(), path)
	if err != nil {
		return nil, err
	}

	if visible {
		return os.Open(path)
	}

	return nil, os.ErrNotExist
}
