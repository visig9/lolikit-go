package runnote

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/visig/lolikit-go/loli2"
	"gitlab.com/visig/pager"
)

func TestCheck(t *testing.T) {
	page := pager.Pager{
		Items:    []interface{}{1, 2, 3},
		PageSize: 3,
	}.Page(1)

	cases := []struct {
		rid int  // runid
		we  bool // wanted error
	}{
		{0, true},
		{1, false},
		{2, false},
		{3, false},
		{4, true},
	}

	for _, c := range cases {
		if c.we {
			assert.Error(t, check(page, c.rid))
		} else {
			assert.Nil(t, check(page, c.rid))
		}
	}
}

func TestGetNote(t *testing.T) {
	mn := new(mockINote)

	page := pager.Pager{
		Items:    []interface{}{mn, 2, 3},
		PageSize: 3,
	}.Page(1)

	cases := []struct {
		rid int        // runid
		wp  bool       // wanted panic
		w   loli2.Note // want
	}{
		{0, true, nil}, // out of range
		{1, false, mn},
		{2, true, nil}, // not loli2.Note
		{3, true, nil}, // not loli2.Note
		{4, true, nil}, // out of range
	}

	for _, c := range cases {
		fn := func() {
			assert.Equal(t, c.w, getNote(page, c.rid))
		}

		if c.wp {
			assert.Panics(t, fn)
		} else {
			assert.NotPanics(t, fn)
		}
	}
}

func TestBuildArgs(t *testing.T) {
	cases := []struct {
		rt Type     // run type
		we bool     // want err
		w  []string // want
	}{
		{
			RunContent,
			false, []string{"clcr", "cp"},
		},
		{
			RunDir,
			false, []string{"cldr", "dir"},
		},
		{ // using content runner due to dir/path is a fake path
			RunEntry,
			false, []string{"clcr", "dir/path"},
		},
		{
			RunAttachmentDir,
			true, []string{}, // err on not a complex note
		},
	}

	inRunner := []string{}

	for _, c := range cases {
		mc := &MockCfg{}
		mc.On("ListContentRunner", inRunner, "ct").Return([]string{"clcr"}, nil).Maybe()
		mc.On("ListRunner", inRunner).Return([]string{"clr"}, nil).Maybe()
		mc.On("ListDirRunner", inRunner).Return([]string{"cldr"}, nil).Maybe()

		mn := &mockINote{}
		mn.On("ContentType").Return("ct").Maybe()
		mn.On("Path").Return(filepath.Join("dir", "path")).Maybe()
		mn.On("ContentPath").Return("cp").Maybe()

		args, err := buildArgs(inRunner, c.rt, mn, mc)
		if c.we {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, c.w, args)
		}

		mn.AssertExpectations(t)
		mc.AssertExpectations(t)
	}
}
