package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCfg(t *testing.T) {
	cases := []struct {
		nr     bool  // needRepo
		grperr error // getRepoPath return err
		we     bool  // want error
	}{
		{false, errors.New("err"), false},
		{false, nil, false},
		{true, errors.New("err"), true},
		{true, nil, false},
	}

	for _, c := range cases {
		uv := new(mockIViper)
		rv := new(mockIViper)

		mu := new(mockIUtil)

		// testing how internal work... it may be too over?
		mu.On("getUserConfPath", "xch", "home").Return("ucp")
		mu.On("getViper", "ucp").Return(uv, nil)
		mu.On("getDefaultRepoPath", uv).Return("drp")
		mu.On("getRepoPath", "frp", "cwd", "drp").Return("rp", c.grperr)

		var expected *Cfg
		if c.grperr != nil {
			expected = &Cfg{uv: uv, rv: rv}
			if !c.we {
				mu.On("getViper", "").Return(rv, nil)
			}
		} else {
			expected = &Cfg{uv: uv, rv: rv, repoPath: "rp"}

			mu.On("getRepoConfPath", "rp").Return("rcp")
			mu.On("getViper", "rcp").Return(rv, nil)
		}

		if cfg, err := createCfg("frp", "xch", "home", "cwd", c.nr, mu); c.we {
			assert.Error(t, err)
		} else {
			if assert.Nil(t, err) {
				assert.Exactly(t, expected, cfg)
			}
		}

		mu.AssertExpectations(t)
	}
}
