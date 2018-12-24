package loli2

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// contentType return the content type (filename ext without dot)
func contentType(contentPath string) string {
	base := filepath.Base(contentPath)

	slice := strings.SplitAfter(base, ".")

	return slice[len(slice)-1]
}

// GetMTime with a cache pointer
func getMTime(mtimeCachePtr *time.Time, contentPath string) *time.Time {
	if mtimeCachePtr.IsZero() {
		fi, err := os.Stat(contentPath)
		if err != nil { // it should not happen
			panic(err)
		}

		*mtimeCachePtr = fi.ModTime()
	}

	return mtimeCachePtr
}

type noteJSONData struct {
	Path        string    `json:"path"`
	Title       string    `json:"title"`
	ContentPath string    `json:"contentPath"`
	ContentType string    `json:"contentType"`
	Type        string    `json:"type"`
	ModTime     time.Time `json:"modTime"`
}

// GetNoteJSON return a json represent of a note.
func getNoteJSON(n Note, noteType string) []byte {
	data, err := json.Marshal(noteJSONData{
		Path:        n.Path(),
		Title:       n.Title(),
		ContentPath: n.ContentPath(),
		ContentType: n.ContentType(),
		ModTime:     n.MTime(),
		Type:        noteType,
	})
	if err != nil {
		panic(err)
	}

	return data
}
