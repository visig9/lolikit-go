package loli2

import (
	"time"
)

// Note is interface to represent a note in lolinote 2
type Note interface {
	Path() string
	Title() string
	ContentPath() string
	ContentType() string
	MTime() time.Time
}

// Entry is a interface reference to Lolinote fs entry.
type Entry interface {
	Path() string
}
