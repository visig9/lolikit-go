package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUGetViper(t *testing.T) {
	u := new(util)

	cases := []struct {
		fp string // filepath
		wv string // expected field "v"'s value in cfg file
		we bool   // expected err
	}{
		{"testdata/cfg/ok1.toml", "1", false},
		{"testdata/cfg/ok2.toml", "2", false},
		{"testdata/cfg/not-exist.toml", "", false},
		{"testdata/cfg/bad.toml", "", true},
	}

	for _, c := range cases {
		if vip, err := u.getViper(c.fp); c.we {
			assert.Error(t, err)
		} else {
			if assert.Nil(t, err) {
				v := vip.GetString("v")
				assert.Equal(t, c.wv, v)
			}
		}
	}
}

func TestUGetUserConfPath(t *testing.T) {
	u := new(util)

	cases := []struct {
		xch      string
		home     string
		expected string
		ep       bool // expect panic
	}{
		{
			"", "home",
			filepath.Join("home", ".config", "lolikit", "config.toml"),
			false,
		},
		{
			"xch", "home",
			filepath.Join("xch", "lolikit", "config.toml"),
			false,
		},
		{"", "", "", true},
	}

	for _, c := range cases {
		if c.ep {
			assert.Panics(
				t,
				func() { u.getUserConfPath(c.xch, c.home) },
			)
		} else {
			assert.Equal(
				t,
				c.expected,
				u.getUserConfPath(c.xch, c.home),
			)
		}
	}
}

func TestUGetRepoConfPath(t *testing.T) {
	u := new(util)

	cases := []struct {
		rp string // repopath
		e  string // expected
		ep bool   // expected panic
	}{
		{
			"",
			filepath.Join(".lolinote", "lolikit", "config.toml"),
			true,
		},
		{
			"path",
			filepath.Join("path", ".lolinote", "lolikit", "config.toml"),
			false,
		},
	}

	for _, c := range cases {
		if c.ep {
			assert.Panics(t, func() { u.getRepoConfPath(c.rp) })
		} else {
			assert.Equal(t, c.e, u.getRepoConfPath(c.rp))
		}
	}
}

func TestUGetDefaultRepoPath(t *testing.T) {
	u := new(util)

	cases := []struct {
		dp string // default repo path
		e  string // expected
	}{
		{"a", "a"},
		{"b", "b"},
		{"v2", "v2"},
	}

	for _, c := range cases {
		mv := new(mockIViper)
		mv.On("GetString", "default-repo").Return(c.dp).Once()

		assert.Equal(t, c.e, u.getDefaultRepoPath(mv))
		mv.AssertExpectations(t)
	}
}

func TestUGetRepoPath(t *testing.T) {
	u := new(util)

	cases := []struct {
		ifp string // input flagpath
		icp string // input cwd path
		idp string // input default path
		e   string // expect repo path
		ee  bool   // expect err
	}{
		{"testdata/r/r1", "", "", "testdata/r/r1", false},
		{"", "testdata/r/r1", "", "testdata/r/r1", false},
		{"", "", "testdata/r/r1", "testdata/r/r1", false},
		{"", "", "", "", true},
		{
			"not-exist", "testdata/r/r1", "",
			"", true,
		},
		{
			"", "not-exist", "testdata/r/r1",
			"testdata/r/r1", false,
		},
		{
			"", "", "not-exist",
			"", true,
		},
		{
			"testdata/r/r1/sub-dir", "", "",
			"testdata/r/r1", false,
		},
	}

	for _, c := range cases {
		rp, err := u.getRepoPath(c.ifp, c.icp, c.idp)

		if c.ee {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}

		assert.Equal(t, c.e, rp)
	}
}
