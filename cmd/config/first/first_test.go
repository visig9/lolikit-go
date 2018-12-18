package first

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOne(t *testing.T) {
	cases := []struct {
		is []interface{} // int slice
		w  interface{}   // want
		wp bool          // want panic
	}{
		{[]interface{}{1, 2, 3, 4}, 1, false},
		{[]interface{}{1, 2, 3}, 1, false},
		{[]interface{}{0, 1, 2, 3, 4}, 1, false},
		{[]interface{}{0, 0, 4}, 4, false},
		{[]interface{}{0, 0, 0}, 0, false},

		{[]interface{}{"", 0, 1}, 1, false},
		{[]interface{}{"", "a"}, "a", false},
		{[]interface{}{"", ""}, "", false},
		{[]interface{}{[]string{"a"}}, []string{"a"}, false},
		{[]interface{}{[]string{}, []string{"a"}}, []string{"a"}, false},
		{[]interface{}{[]string{""}, []string{"a"}}, []string{""}, false},

		{[]interface{}{}, nil, true},
	}

	for _, c := range cases {
		funcP := func() { assert.Equal(t, c.w, One(c.is...)) }

		if c.wp {
			assert.Panics(t, funcP)
		} else {
			assert.NotPanics(t, funcP)
		}

	}
}
