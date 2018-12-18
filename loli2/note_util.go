package loli2

import (
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
