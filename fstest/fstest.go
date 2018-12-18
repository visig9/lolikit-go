// Package fstest is some utility function to test filesystem status.
package fstest

import (
	"os"
)

// IsDir test the path is a dir.
func IsDir(path string) bool {
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		return true
	}
	return false
}

// IsRegular test the path is a regular file.
func IsRegular(path string) bool {
	if fi, err := os.Stat(path); err == nil && fi.Mode().IsRegular() {
		return true
	}
	return false
}

// IsExist test the path is exists.
func IsExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
