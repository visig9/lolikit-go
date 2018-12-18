package config

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
	"gitlab.com/visig/lolikit-go/fstest"
	"gitlab.com/visig/lolikit-go/loli2"
)

//go:generate mockery -name iViper -testonly -inpkg

type iViper interface {
	GetString(string) string
	GetStringSlice(string) []string
	GetInt(string) int
}

//go:generate mockery -name iUtil -testonly -inpkg

// iUtil is internal utility level 1 interface.
type iUtil interface {
	getViper(confPath string) (iViper, error)
	getUserConfPath(xdgConfigHome, home string) (confPath string)
	getRepoConfPath(repoPath string) (confPath string)
	getDefaultRepoPath(uv iViper) (defaultRepoPath string)
	getRepoPath(flagRepoPath, cwdPath, defaultPath string) (
		repoPath string, err error,
	)
}

type util struct{}

func (util) getViper(confpath string) (iViper, error) {
	v := viper.New()

	if fstest.IsExist(confpath) { // found a config file
		v.SetConfigFile(confpath)
		if err := v.ReadInConfig(); err != nil {
			err := fmt.Errorf(
				"config file %q reading failed: %v",
				confpath, err,
			)
			return nil, err
		}
	}

	return v, nil
}

func (util) getUserConfPath(xdgConfigHome, home string) (confPath string) {
	var frags = []string{"", "lolikit", "config.toml"}

	switch {
	case len(xdgConfigHome) > 0:
		frags[0] = xdgConfigHome
	case len(home) > 0:
		frags[0] = filepath.Join(home, ".config")
	default:
		panic("both xdgConfigHome and home not given")
	}

	return filepath.Join(frags...)
}

func (util) getRepoConfPath(repoPath string) (confPath string) {
	if len(repoPath) == 0 {
		panic("repoPath not given")
	}

	return filepath.Join(
		repoPath, ".lolinote", "lolikit", "config.toml",
	)
}

func (util) getDefaultRepoPath(uv iViper) (defaultRepoPath string) {
	return uv.GetString("default-repo")
}

func (util) getRepoPath(flagRepoPath, cwdPath, defaultPath string) (
	repoPath string, err error,
) {
	notRepoErr := func(path string) error {
		return errors.New("Not a Lolinote repository: " + path)
	}

	// Force using flagRepoPath if given
	if flagRepoPath != "" {
		repo, err := loli2.NewRepo(flagRepoPath, false)
		if err != nil {
			return "", notRepoErr(flagRepoPath)
		}

		return repo.Path(), nil
	}

	// Check working directory is valid repo. If it's, use it.
	cwdRepo, err := loli2.NewRepo(cwdPath, false)
	if err == nil {
		return cwdRepo.Path(), nil
	}

	// If working directory not valid, try default-repo.
	if defaultPath != "" {
		defaultRepo, err := loli2.NewRepo(defaultPath, false)
		if err != nil {
			return "", notRepoErr(defaultPath)
		}

		return defaultRepo.Path(), nil
	}

	// default-repo not given
	return "", errors.New("Lolinote repository not found")
}
