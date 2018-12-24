package loli2

import (
	"path/filepath"
	"strings"
	"time"
)

// ComplexNote is a complex note instance in lolinote.
type ComplexNote struct {
	path       string
	mtimeCache time.Time
}

// Path return note's entry.
func (n *ComplexNote) Path() string {
	return n.path
}

// Title return note's title.
func (n *ComplexNote) Title() string {
	return filepath.Base(n.path)
}

// ContentPath return note's content filepath.
func (n *ComplexNote) ContentPath() string {
	indexPath, err := getComplexNoteIndexPath(n.path)
	if err != nil { // should not happen
		panic(err)
	}

	return indexPath
}

// ContentType return note's conent type.
func (n *ComplexNote) ContentType() string {
	return contentType(n.ContentPath())
}

// MTime return note's modtime.
func (n *ComplexNote) MTime() time.Time {
	return *getMTime(&n.mtimeCache, n.ContentPath())
}

// JSON return the json string of this note.
func (n *ComplexNote) JSON() []byte {
	return getNoteJSON(n, "ComplexNote")
}

// The index page finder

type complexNoteIndexError struct {
	tooMuch  bool
	notFound bool
}

func (e complexNoteIndexError) Error() string {
	if e.tooMuch {
		return "count of index.* path number > 1"
	}

	return "count of index.* path number < 1"
}

// getComplexNoteIndexPath return the complex note index filepath
// by note's path.
//
// if index path number is not equal to 1, error != nil.
func getComplexNoteIndexPath(path string) (string, error) {
	matches, _ := filepath.Glob(path + "/index.*")
	filteredMatches := make([]string, 0, len(matches))
	for _, m := range matches {
		mbase := filepath.Base(m)

		if strings.Count(mbase, ".") == 1 &&
			mbase != "index." {
			filteredMatches = append(filteredMatches, m)

			if count := len(filteredMatches); count > 1 {
				return "", &complexNoteIndexError{
					tooMuch: true,
				}
			}
		}
	}

	if count := len(filteredMatches); count != 1 {
		return "", &complexNoteIndexError{notFound: true}
	}

	return filteredMatches[0], nil
}
