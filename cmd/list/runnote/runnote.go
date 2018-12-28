package runnote

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gitlab.com/visig/lolikit-go/fstest"
	"gitlab.com/visig/lolikit-go/logger"
	"gitlab.com/visig/lolikit-go/loli2"
	"gitlab.com/visig/pager"
)

// Type define how to run the openNote function (must choicing one)
type Type int

const (
	// RunContent open Path with content type.
	RunContent Type = iota
	// RunDir open filepath.Dir(Path) with dir runner.
	RunDir
	// RunAttachmentDir open Path with dir runner
	RunAttachmentDir
	// RunEntry open Path with default runner.
	RunEntry
)

//go:generate mockery -name iNote -testonly -inpkg

type iNote interface {
	loli2.Note
}

func check(page pager.Page, runid int) error {
	if runid > page.Size() || runid < 1 {
		return fmt.Errorf(
			"runid (%d) not found in current page", runid,
		)
	}

	return nil
}

func getNote(page pager.Page, runid int) iNote {
	v := page.Items()[runid-1]

	return v.Data().(iNote)
}

//go:generate mockery -name Cfg -testonly -inpkg

// Cfg is a config object interface.
type Cfg interface {
	ListRunner([]string) ([]string, error)
	ListContentRunner([]string, string) ([]string, error)
	ListDirRunner([]string) ([]string, error)
}

// ReplacePath find the placeholder "{}" in runner and replace it.
//
// If placeholder exists, returned ok == true, else false.
func replacePath(runner []string, path string) (
	newRunner []string, ok bool,
) {
	newRunner = make([]string, len(runner))

	for i, arg := range runner {
		if arg == "{}" {
			newRunner[i] = path
			ok = true
		} else {
			newRunner[i] = arg
		}
	}

	return
}

func buildArgs(
	inRunner []string, rt Type, n iNote, cfg Cfg,
) (args []string, err error) {
	var path string
	var runner []string

	switch rt {
	case RunContent:
		path = n.ContentPath()
		runner, err = cfg.ListContentRunner(
			inRunner, n.ContentType(),
		)
	case RunDir:
		path = filepath.Dir(n.Path())
		runner, err = cfg.ListDirRunner(inRunner)
	case RunAttachmentDir:
		if _, ok := n.(*loli2.ComplexNote); !ok {
			return nil, errors.New(n.Path() + " not a complex note")
		}
		path = n.Path()
		runner, err = cfg.ListDirRunner(inRunner)
	case RunEntry:
		path = n.Path()
		if fstest.IsDir(path) {
			runner, err = cfg.ListDirRunner(inRunner)
		} else {
			runner, err = cfg.ListContentRunner(
				inRunner, n.ContentType(),
			)
		}
	}

	if err != nil {
		return args, err
	}

	var ok bool
	args, ok = replacePath(runner, path)
	if !ok {
		args = append(runner, path)
	}

	return args, err
}

func runArgs(args []string) error {
	// logger.Std.Print("list: ", shellquote.Join(args...))

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Run can run a note by runner.
func Run(
	page pager.Page,
	runid int,
	inRunner []string,
	cfg Cfg,
	rt Type,
) {
	err := check(page, runid)
	if err != nil {
		logger.Err.Fatal(err)
	}

	n := getNote(page, runid)

	if args, err := buildArgs(inRunner, rt, n, cfg); err != nil {
		logger.Err.Fatal(err)
	} else {
		if err := runArgs(args); err != nil {
			logger.Err.Fatal(err)
		}
	}
}
