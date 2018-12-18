package loli2

import (
	"fmt"
	"os"
)

// Repo represent a Lolinote Repository.
type Repo struct {
	path string
}

// NewRepo create a new Repo object.
//
// If create is true, create a repo at the (exact) path in filesystem.
// If create is false, try to find the repo by path in filesystem.
//
// err != nil only when can not return a vaild repo.
func NewRepo(path string, create bool) (r Repo, err error) {
	if create {
		err = os.MkdirAll(getLolinoteDir(path), 0755)
		r = Repo{path}
	} else {
		if rpath, ok := FindUpperRepo(path); ok {
			r = Repo{rpath}
		} else {
			err = fmt.Errorf(
				"no Lolinote repository at %v",
				path,
			)
		}
	}

	return
}

// Path retuen the path of this repo.
func (r Repo) Path() string {
	return r.path
}

// Dir retuen the Dir of this repo's root.
func (r Repo) Dir() Dir {
	return Dir{r.path}
}
