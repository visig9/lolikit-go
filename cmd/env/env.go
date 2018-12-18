// Package env pre-load all environments and panic if necessary.
package env

import (
	"os"

	homedir "github.com/mitchellh/go-homedir"
)

type env struct{}

func (env) XDGConfigHome() string {
	return xdgConfigHome
}

func (env) Home() string {
	return home
}

func (env) CWD() string {
	return cwd
}

// Env is global env object.
var Env = env{}

// internal variable used by (global) Env.
var (
	xdgConfigHome string
	home          string
	cwd           string
)

func init() {
	var err error

	xdgConfigHome = os.Getenv("XDG_CONFIG_HOME")

	home, err = homedir.Dir()
	if err != nil {
		panic(err)
	}

	cwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}
