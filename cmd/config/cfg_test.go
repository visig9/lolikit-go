package config

import (
	"testing"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/stretchr/testify/assert"
)

func TestCfgRepo(t *testing.T) {
	cases := []struct {
		cfg    Cfg
		epanic bool   // expected panic
		epath  string // expected path
	}{
		{Cfg{repoPath: "testdata/r/r1"}, false, "testdata/r/r1"},
		{Cfg{repoPath: "testdata/r"}, true, ""},
	}

	for _, c := range cases {
		funcP := func() {
			assert.Equal(t, c.epath, c.cfg.Repo().Path())
		}

		if c.epanic {
			assert.Panics(t, funcP)
		} else {
			assert.NotPanics(t, funcP)
		}
	}
}

func TestCfgListPageSize(t *testing.T) {
	cases := []struct {
		inps       int // input's page size
		rvps, uvps int // uv & rv's page size
		e          int // expected
	}{
		{0, 0, 0, 10},
		{0, 0, 1, 1},
		{0, 2, 1, 2},
		{1, 2, 3, 1},
	}

	for _, c := range cases {
		uv := new(mockIViper)
		uv.On("GetInt", "list.page-size").Return(c.uvps)

		rv := new(mockIViper)
		rv.On("GetInt", "list.page-size").Return(c.rvps)

		cfg := Cfg{uv: uv, rv: rv}

		assert.Equal(t, c.e, cfg.ListPageSize(c.inps), "%v", c)
	}
}

func TestCfgTextTypes(t *testing.T) {
	defaultValue := hashset.New("txt", "md")

	cases := []struct {
		rvv, uvv []string     // uv & rv's value
		e        *hashset.Set // expected
	}{
		{[]string{}, []string{}, defaultValue},
		{[]string{"txt"}, []string{}, hashset.New("txt")},
		{[]string{}, []string{"md"}, hashset.New("md")},
		{[]string{"txt"}, []string{"md"}, hashset.New("txt")},
		{[]string{"txt"}, []string{"txt", "md"}, hashset.New("txt")},
		{[]string{"txt", "md"}, []string{"txt"}, hashset.New("txt", "md")},
	}

	for _, c := range cases {
		uv := new(mockIViper)
		uv.On("GetStringSlice", "text-types").Return(c.uvv)

		rv := new(mockIViper)
		rv.On("GetStringSlice", "text-types").Return(c.rvv)

		cfg := Cfg{uv: uv, rv: rv}

		assert.Equal(t, c.e, cfg.TextTypes(), "%v", c)
	}
}

func TestCfgListRunner(t *testing.T) {
	cases := []struct {
		ct           string   // content type
		inr          []string // inRunner
		rvctr, uvctr string   // rv & uv's content type runner
		rvr, uvr     string   // rv & uv's default runner
		w            []string // want
		we           bool     // want err
		wc           []string // want ListContentRunner
		wce          bool     // want ListContentRunner error
	}{
		{ // should err by no usable value
			"txt", []string{},
			"", "",
			"", "",
			[]string{}, true,
			[]string{}, true,
		},
		{ // should err by bad config string
			"txt", []string{},
			"rvctr\"", "",
			"", "",
			[]string{}, true,
			[]string{}, true,
		},
		{
			"txt", []string{"ok"},
			"", "",
			"", "",
			[]string{"ok"}, false,
			[]string{"ok"}, false,
		},
		{
			"txt", []string{"ok"},
			"rvcrt", "",
			"", "",
			[]string{"ok"}, false,
			[]string{"ok"}, false,
		},
		{
			"txt", []string{},
			"rvctr", "",
			"", "",
			[]string{}, true,
			[]string{"rvctr"}, false,
		},
	}

	for _, c := range cases {
		rv := new(mockIViper)
		rv.On("GetString", "list.runners."+c.ct).Return(c.rvctr)
		rv.On("GetString", "list.runner").Return(c.rvr)

		uv := new(mockIViper)
		uv.On("GetString", "list.runners."+c.ct).Return(c.uvctr)
		uv.On("GetString", "list.runner").Return(c.uvr)

		cfg := Cfg{uv: uv, rv: rv}

		if runner, err := cfg.ListRunner(c.inr); c.we {
			assert.Error(t, err, "%v", c)
		} else {
			assert.Nil(t, err, "%v", c)
			assert.Equal(t, c.w, runner, "%v", c)
		}

		if runner, err := cfg.ListContentRunner(c.inr, c.ct); c.wce {
			assert.Error(t, err, "%v", c)
		} else {
			assert.Nil(t, err, "%v", c)
			assert.Equal(t, c.wc, runner, "%v", c)
		}
	}
}

func TestCfgListDirRunner(t *testing.T) {
	cases := []struct {
		inr        []string // inRunner
		rvdr, uvdr string   // rv & uv's dir runner
		rvr, uvr   string   // rv & uv's default runner
		w          []string // want
		we         bool     // want error
	}{
		{ // should err by no usable value
			[]string{},
			"", "",
			"", "",
			[]string{}, true,
		},
		{ // should err by bad config string
			[]string{},
			"rvdr\"", "",
			"", "",
			[]string{}, true,
		},
		{
			[]string{"ok"},
			"", "",
			"", "",
			[]string{"ok"}, false,
		},
		{
			[]string{"ok"},
			"rvdr", "",
			"", "",
			[]string{"ok"}, false,
		},
		{
			[]string{},
			"rvdr", "",
			"", "",
			[]string{"rvdr"}, false,
		},
	}

	for _, c := range cases {
		rv := new(mockIViper)
		rv.On("GetString", "list.dir-runner").Return(c.rvdr)
		rv.On("GetString", "list.runner").Return(c.rvr)

		uv := new(mockIViper)
		uv.On("GetString", "list.dir-runner").Return(c.uvdr)
		uv.On("GetString", "list.runner").Return(c.uvr)

		cfg := Cfg{uv: uv, rv: rv}

		if runner, err := cfg.ListDirRunner(c.inr); c.we {
			assert.Error(t, err, "%v", c)
		} else {
			assert.Nil(t, err, "%v", c)
			assert.Equal(t, c.w, runner, "%v", c)
		}
	}
}
