package serve

import (
	"net/http"
	"os"
	"path/filepath"

	"gitlab.com/visig/lolikit-go/fstest"
	"gitlab.com/visig/lolikit-go/logger"
	"gitlab.com/visig/lolikit-go/loli2"
)

type repoFileSystem struct {
	repo loli2.Repo
}

func (r *repoFileSystem) Open(name string) (http.File, error) {
	pfpath := filepath.FromSlash(name)
	path := filepath.Join(r.repo.Path(), pfpath)

	if fstest.IsDir(path) {
		return nil, os.ErrNotExist
	}

	visible, err := isVisible(r.repo.Path(), path)
	if err != nil {
		return nil, err
	}

	logger.Std.Print(visible, path)

	if visible {
		return os.Open(path)
	}

	return nil, os.ErrNotExist
}
