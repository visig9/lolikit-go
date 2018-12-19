package serve

import (
	"path/filepath"
	"strings"

	"gitlab.com/visig/lolikit-go/pathrange"
)

// IsVisible check a fs path is visible from client or not.
func isVisible(repoPath, path string) (bool, error) {
	pr, err := pathrange.PathRange(repoPath, path)
	if err != nil {
		return false, err
	}

	for _, path := range pr {
		name := filepath.Base(path)
		if strings.HasPrefix(name, ".") {
			return false, nil
		}
	}

	return true, nil
}
