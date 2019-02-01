package loli2

import (
	"fmt"

	"github.com/visig9/lolikit-go/fstest"
)

// buildEntry build a entry instance from specific filesystem path.
func buildEntry(path string) (e Entry, err error) {
	switch {
	case IsSimpleNote(path):
		e = &SimpleNote{path: path}
	case IsComplexNote(path):
		e = &ComplexNote{path: path}
	case IsNoise(path):
		e = Noise{path: path}
	case fstest.IsDir(path):
		e = Dir{path: path}
	default:
		err = fmt.Errorf("%v not match any type of entry", path)
	}

	return
}
