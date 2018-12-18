package parse

import (
	"strconv"
)

// Int convert string to int.
//
// If parse fail, print error message and exit 1.
func Int(str string) (int, error) {
	i64, err := strconv.ParseInt(str, 10, 64)

	return int(i64), err
}
