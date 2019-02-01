// Package nps is Note Page Stringer
package nps

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/visig9/elign"
	"github.com/visig9/lolikit-go/loli2"
	"github.com/visig9/pager"
)

//go:generate mockery -name iNote -testonly -inpkg

type iNote interface { // only for internal test
	loli2.Note
}

func isCN(v interface{}) bool {
	switch v.(type) {
	case *loli2.ComplexNote:
		return true
	}

	return false
}

func parent(path string) string {
	return filepath.Base(filepath.Dir(path))
}

func inverseItemPrint(
	page pager.Page,
	itemStringer func(pager.PageItem) string,
) string {
	size := page.Size()
	strs := make([]string, size)

	for i, pi := range page.Items() {
		strs[size-1-i] = itemStringer(pi)
	}

	return strings.Join(strs, "\n")
}

type align func(string) string

func getAlign(page pager.Page) (tAlign align, ctAlign align) {
	et := elign.Default(0)
	ect := elign.Default(0)

	for _, pi := range page.Items() {
		n := pi.Data().(iNote)
		et.AdjustWidth(n.Title())
		ect.AdjustWidth("[" + n.ContentType() + "]")
	}

	return et.Left, ect.Left
}

func getItemToString(tAlign, ctAlign align, pageNumber int) func(pager.PageItem) string {
	return func(pi pager.PageItem) (out string) {
		if pageNumber != 1 {
			out += fmt.Sprintf("%2v} ", pi.GlobalIndex()+1)
		}

		out += fmt.Sprintf("%2v)", pi.InPageIndex()+1)

		if isCN(pi.Data()) {
			out += " + "
		} else {
			out += "   "
		}

		n := pi.Data().(iNote)
		out += fmt.Sprintf(
			"%v  %v  << %v",
			tAlign(n.Title()),
			ctAlign("["+n.ContentType()+"]"),
			parent(n.Path()),
		)

		return
	}
}

// ToString convert pager.Item (data = loli2.Note) to string.
func ToString(page pager.Page) (out string) {
	tAlign, ctAlign := getAlign(page)
	itemToString := getItemToString(tAlign, ctAlign, page.PageNumber())

	return inverseItemPrint(page, itemToString)
}
