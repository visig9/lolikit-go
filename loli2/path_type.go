package loli2

import (
	"path/filepath"
	"strings"

	"github.com/visig9/lolikit-go/fstest"
)

// IsRepo detect a path is a Repo or not.
func IsRepo(path string) bool {
	cd := getLolinoteDir(path)

	return fstest.IsDir(cd)
}

// IsSimpleNote check a path is a simple note.
//
// It would not consider this path may shadowed by some noise or
// attachment dir.
func IsSimpleNote(path string) bool {
	if !fstest.IsRegular(path) {
		return false
	}

	// start with .
	base := filepath.Base(path)
	if string(base[0]) == "." {
		return false
	}

	// without any "." (no chance has a ext)
	baseslice := strings.Split(base, ".")
	if len(baseslice) == 1 {
		return false
	}

	// without ext (content type)
	if ext := baseslice[len(baseslice)-1]; ext == "" {
		return false
	}

	// main filename != index
	baseWithoutExt := strings.Join(baseslice[:len(baseslice)-1], ".")
	if baseWithoutExt == "index" {
		return false
	}

	return true
}

// IsComplexNote check a path is a complex note.
//
// It would not consider this path may shadowed by some noise dir.
func IsComplexNote(path string) bool {
	if !fstest.IsDir(path) {
		return false
	}

	_, err := getComplexNoteIndexPath(path)

	if err != nil {
		return false
	}

	return true
}

// IsNoise test the path is a noise in lolinote.
//
// This function doesn't consider it may be shadowed by other noise
// or attachment directory.
func IsNoise(path string) bool {
	base := filepath.Base(path)

	// start with "."
	if string(base[0]) == "." {
		return true
	}

	// a file without filename ext
	extWithDot := filepath.Ext(base)
	if fstest.IsRegular(path) &&
		(extWithDot == "." || extWithDot == "") {
		return true
	}

	// is a sub lolinote repository
	if IsRepo(path) {
		return true
	}

	// contain multiple index
	_, err := getComplexNoteIndexPath(path)
	if err != nil && err.(*complexNoteIndexError).tooMuch {
		return true
	}

	return false
}
