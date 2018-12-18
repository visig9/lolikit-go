package nps

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/visig/lolikit-go/loli2"
	"gitlab.com/visig/pager"
)

func TestIsCN(t *testing.T) {
	cases := []struct {
		v interface{}
		e bool // expected
	}{
		// {1, false},
		// {"2", false},
		{&loli2.SimpleNote{}, false},
		{&loli2.ComplexNote{}, true},
	}

	for _, c := range cases {
		if c.e {
			assert.True(t, isCN(c.v))
		} else {
			assert.False(t, isCN(c.v))
		}
	}
}

func TestInverseItemPrint(t *testing.T) {
	itemStringer := func(pi pager.PageItem) string {
		return pi.Data().(string)
	}

	cases := []struct {
		p pager.Page
		e string // expected
	}{
		{
			pager.Pager{
				Items:    []interface{}{"1", "2", "3"},
				PageSize: 3,
			}.Page(1),
			"3\n2\n1",
		},
		{
			pager.Pager{
				Items:    []interface{}{"a", "b", "c"},
				PageSize: 3,
			}.Page(1),
			"c\nb\na",
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.e, inverseItemPrint(c.p, itemStringer))
	}
}
