package pathrange

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathRange(t *testing.T) {
	cases := []struct {
		from, to string   // input
		w        []string // want
		we       bool     // want error
	}{
		{ // rel path
			"a/b", "a/b/c/d/e",
			[]string{"a/b/c", "a/b/c/d", "a/b/c/d/e"}, false,
		},
		{
			"a/b", "a/b/c/d",
			[]string{"a/b/c", "a/b/c/d"}, false,
		},
		{ // abs path
			"/a/b", "/a/b/c/d",
			[]string{"/a/b/c", "/a/b/c/d"}, false,
		},
		{ // the same (rel) path
			"a/b", "a/b",
			nil, false,
		},
		{ // the same (abs) path
			"/a/b", "/a/b",
			nil, false,
		},
		{ // not normalized path
			"a/../a/b/", "a/b/d/../../b/c/d/.",
			[]string{"a/b/c", "a/b/c/d"}, false,
		},
		{ // "to" is the parent of "from"
			"a/b/c/d", "a/b",
			[]string{"a/b/c", "a/b"}, false,
		},
		{ // abs to rel
			"/a/b", "a/b/c/d",
			nil, true,
		},
		{ // rel to abs
			"a/b", "/a/b/c/d",
			nil, true,
		},
	}

	for _, c := range cases {
		pr, err := PathRange(c.from, c.to)
		if c.we {
			assert.Error(t, err)
			assert.Equal(t, c.w, pr)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, c.w, pr)
		}
	}
}

func ExamplePathRange() {
	cases := []struct {
		from, to string
	}{
		{"a/b", "a/b/c/d/e"},
		{"/a/b/c/d/e", "/a/b"},
		{"a/b", "a/b"},
		{"/a/b", "a/b"},
	}

	for _, c := range cases {
		pr, err := PathRange(c.from, c.to)

		fmt.Printf("%v, %v\n", pr, err != nil)
	}

	// Output:
	// [a/b/c a/b/c/d a/b/c/d/e], false
	// [/a/b/c/d /a/b/c /a/b], false
	// [], false
	// [], true
}
