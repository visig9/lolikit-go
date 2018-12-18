// Package pathrange offer a function to build a range of paths between
// two paths.
package pathrange

import (
	"os"
	"path/filepath"
	"strings"
)

// PathRange build a range of paths between two paths.
//
// The return slice would not contain "from" path. If build a range is
// impossiable, err != nil.
func PathRange(from, to string) ([]string, error) {
	from = filepath.Clean(from)

	diffPath, err := filepath.Rel(from, to)
	if err != nil {
		return nil, err
	}
	if diffPath == "." { // no difference
		return nil, nil
	}

	diffParts := strings.Split(diffPath, string(os.PathSeparator))

	results := make([]string, 0, len(diffParts))
	for i := 1; i < len(diffParts)+1; i++ {
		results = append(
			results,
			filepath.Join(
				from, filepath.Join(diffParts[0:i]...),
			),
		)
	}

	return results, nil
}
