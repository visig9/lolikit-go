package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt(t *testing.T) {
	cases := []struct {
		s  string // string
		w  int    // want
		we bool   // want err?
	}{
		{"1", 1, false},
		{"25", 25, false},
		{"25.5", 0, true},
		{"二", 0, true},
		{"", 0, true},
	}

	for _, c := range cases {
		i, err := Int(c.s)
		if c.we {
			assert.Error(t, err)
		} else {
			if assert.Nil(t, err) {
				assert.Equal(t, c.w, i)
			}
		}
	}
}
