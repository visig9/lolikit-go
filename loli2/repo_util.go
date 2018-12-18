package loli2

import (
	"path/filepath"
)

// GetLolinoteDir can return the .lolinote folder path by repoPath.
func getLolinoteDir(repoPath string) string {
	return filepath.Join(repoPath, ".lolinote")
}

// FindUpperRepo will detect all upper paths and return the nearest repo path.
func FindUpperRepo(path string) (repoPath string, ok bool) {
	for ; path != filepath.Dir(path); path = filepath.Dir(path) {
		if IsRepo(path) {
			return path, true
		}
	}

	return "", false
}
