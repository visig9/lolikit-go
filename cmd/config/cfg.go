package config

import (
	"errors"

	"github.com/emirpasic/gods/sets/hashset"
	shellquote "github.com/kballard/go-shellquote"

	"github.com/visig9/lolikit-go/cmd/config/first"
	"github.com/visig9/lolikit-go/logger"
	"github.com/visig9/lolikit-go/loli2"
)

// Cfg is config object of whole project.
type Cfg struct {
	uv       iViper
	rv       iViper // if repo not found, a empty (valid) rv in here
	repoPath string
}

// Env is environment object of whole project.
type Env interface {
	XDGConfigHome() string
	Home() string
	CWD() string
}

// New create a Cfg object.
func New(env Env, inRepoPath string, needRepo bool) *Cfg {
	c, err := createCfg(
		inRepoPath,
		env.XDGConfigHome(), env.Home(), env.CWD(),
		needRepo,
		util{},
	)

	if err != nil {
		logger.Err.Fatal(err)
	}

	return c
}

// Repo return repo from Cfg.
func (c *Cfg) Repo() loli2.Repo {
	repo, err := loli2.NewRepo(c.repoPath, false)

	if err != nil {
		panic(err)
	}

	return repo
}

// ListPageSize return the list page size value.
func (c *Cfg) ListPageSize(inPageSize int) int {
	return first.One(
		inPageSize,
		c.rv.GetInt("list.page-size"),
		c.uv.GetInt("list.page-size"),
		10,
	).(int)
}

// ServeAddr return the list page size value.
func (c *Cfg) ServeAddr(inAddr string) string {
	return first.One(
		inAddr,
		c.rv.GetInt("serve.address"),
		c.uv.GetInt("serve.address"),
		":10204",
	).(string)
}

// TextTypes return the text-types.
func (c *Cfg) TextTypes() *hashset.Set {
	sli := first.One(
		c.rv.GetStringSlice("text-types"),
		c.uv.GetStringSlice("text-types"),
		[]string{"txt", "md"},
	).([]string)

	out := hashset.New()
	for _, v := range sli {
		out.Add(v)
	}

	return out
}

// ListContentRunner return the runner for the content type.
func (c *Cfg) ListContentRunner(
	inRunner []string, contentType string,
) (runner []string, err error) {
	if len(inRunner) > 0 {
		runner = inRunner
	} else {
		runner, err = shellquote.Split(first.One(
			c.rv.GetString("list.runners."+contentType),
			c.uv.GetString("list.runners."+contentType),
			c.rv.GetString("list.runner"),
			c.uv.GetString("list.runner"),
		).(string))

		if err == nil && len(runner) == 0 {
			err = errors.New("runner not found, please assign it first")
		}
	}

	return
}

// ListDirRunner return the dir runner.
func (c *Cfg) ListDirRunner(inRunner []string) (
	runner []string, err error,
) {
	if len(inRunner) > 0 {
		runner = inRunner
	} else {
		runner, err = shellquote.Split(first.One(
			c.rv.GetString("list.dir-runner"),
			c.uv.GetString("list.dir-runner"),
			c.rv.GetString("list.runner"),
			c.uv.GetString("list.runner"),
		).(string))

		if err == nil && len(runner) == 0 {
			err = errors.New("runner not found, please assign it first")
		}
	}

	return
}

// NewBuffer return the New buffer directory name.
func (c *Cfg) NewBuffer() string {
	return first.One(
		c.rv.GetString("new.buffer"),
		c.uv.GetString("new.buffer"),
		"to-be-classified",
	).(string)
}

// NewBufferSize return the add buffer-size.
func (c *Cfg) NewBufferSize() int {
	return first.One(
		c.rv.GetInt("new.buffer-size"),
		c.uv.GetInt("new.buffer-size"),
		10,
	).(int)
}
