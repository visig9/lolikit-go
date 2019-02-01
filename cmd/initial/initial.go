// Package initial initialize a lolinote 2.0 repo in particular dirpath.
package initial

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/visig9/lolikit-go/cmd/env"
	"github.com/visig9/lolikit-go/logger"
	"github.com/visig9/lolikit-go/loli2"
)

func assertArgsNumber(args []string) {
	if len(args) > 1 {
		logger.Err.Fatal("too many arguments")
	}
}

func getPath(cwd string, args []string) string {
	path := cwd
	if len(args) > 0 {
		path = args[0]
	}

	return path
}

func assertNotAlreadyExists(path string) {
	if loli2.IsRepo(path) {
		logger.Err.Fatal(
			"Lolinote repository already exists: " + path,
		)
	}
}

func assertNotUnderARepo(path string) {
	absPath, _ := filepath.Abs(path)
	if upPath, found := loli2.FindUpperRepo(absPath); found {
		logger.Err.Printf(
			"%v has a parent repository: %v",
			absPath, upPath,
		)
		logger.Err.Fatal("use `-f` to force create it")
	}
}

func createRepo(path string) {
	repo, err := loli2.NewRepo(path, true)
	if err != nil {
		panic(err)
	} else {
		logger.Std.Printf(
			"Initialized a Lolinote repository in %v",
			repo.Path(),
		)
	}
}

// InitCmd offer a command for initialize a lolinote repo.
var InitCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize a Lolinote repository",
	Long: `Initialize a Lolinote repository

  By default, this command will initialize a lolinote repository in
  current working directory.

  If given a path, a directory will be created (if not exists) and
  initialized.
`,
	Run: func(cmd *cobra.Command, args []string) {
		inForce, _ := cmd.Flags().GetBool("force")

		assertArgsNumber(args)
		path := getPath(env.Env.CWD(), args)
		assertNotAlreadyExists(path)

		if !inForce {
			assertNotUnderARepo(path)
		}

		createRepo(path)
	},
}

func init() {
	InitCmd.Flags().BoolP(
		"force", "f", false,
		"force create a repository even parent repository exists",
	)
}
