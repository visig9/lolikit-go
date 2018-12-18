package loli2

import (
	"path/filepath"
	"strings"
	"time"
)

// SimpleNote is a simple note instance in lolinote.
type SimpleNote struct {
	path       string
	mtimeCache time.Time
}

// Path return note's entry.
func (n SimpleNote) Path() string {
	return n.path
}

// Title return note's title.
func (n SimpleNote) Title() string {
	base := filepath.Base(n.path)
	slice := strings.Split(base, ".")

	return strings.Join(slice[:len(slice)-1], ".")
}

// ContentPath return note's content filepath.
func (n SimpleNote) ContentPath() string {
	return n.path
}

// ContentType return note's conent type.
func (n SimpleNote) ContentType() string {
	return contentType(n.ContentPath())
}

// MTime return note's modtime.
func (n *SimpleNote) MTime() time.Time {
	return *getMTime(&n.mtimeCache, n.ContentPath())
}
