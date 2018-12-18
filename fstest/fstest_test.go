package fstest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDir(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{"testdata/exist-dir", true},
		{"testdata/non-exist-dir", false},
		{"testdata/exist-dir/exist.txt", false},
	}

	for _, c := range cases {
		assert.Equal(t, c.want, IsDir(c.path))
	}
}

func TestIsRegular(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{"testdata/exist-dir", false},
		{"testdata/non-exist-dir", false},
		{"testdata/exist-dir/exist.txt", true},
	}

	for _, c := range cases {
		assert.Equal(t, c.want, IsRegular(c.path))
	}
}

func TestIsExist(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{"testdata/exist-dir", true},
		{"testdata/non-exist-dir", false},
		{"testdata/exist-dir/exist.txt", true},
	}

	for _, c := range cases {
		assert.Equal(t, c.want, IsExist(c.path))
	}
}
