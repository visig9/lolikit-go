// Package newnote can add a new simple note in buffer directory.
package newnote

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	shellquote "github.com/kballard/go-shellquote"
	"github.com/spf13/cobra"
	"github.com/visig9/lolikit-go/cmd/config"
	"github.com/visig9/lolikit-go/cmd/env"
	"github.com/visig9/lolikit-go/logger"
)

func verifyArgs(args []string) {
	if len(args) <= 1 {
		logger.Err.Fatal("arguments not enough")
	}

	noteSig := args[0]
	if strings.Contains(noteSig, string(os.PathSeparator)) {
		logger.Err.Fatalf(
			"%v contain invalid characters",
			noteSig,
		)
	}

	splitNoteSig := strings.Split(noteSig, ".")
	contentType := splitNoteSig[len(splitNoteSig)-1]
	title := splitNoteSig[0]
	if len(splitNoteSig) <= 1 || title == "" || contentType == "" {
		logger.Err.Fatalf(
			"%v not a correct note signature",
			noteSig,
		)
	}
}

func prepareBufferArea(bufferPath string) {
	switch fi, err := os.Stat(bufferPath); {
	case err != nil && os.IsNotExist(err):
		if err := os.MkdirAll(bufferPath, 0755); err != nil {
			logger.Err.Fatal(err)
		}
	case fi.IsDir() == false:
		logger.Err.Fatalf("%v not a directory", bufferPath)
	}
}

func verifyBufferAreaSize(
	bufferPath string, inForce bool, bufferSize int,
) {
	f, err := os.Open(bufferPath)
	if err != nil {
		logger.Err.Fatal(err)
	}

	switch names, err := f.Readdirnames(0); {
	case err != nil:
		logger.Err.Fatal(err)
	case !inForce && len(names) >= bufferSize:
		logger.Err.Fatalf(
			"buffer directory %v has too many entries to be classified",
			bufferPath,
		)
	}
}

func verifyNotePath(notePath string) {
	if _, err := os.Stat(notePath); !os.IsNotExist(err) {
		logger.Err.Fatalf("%v already exists", notePath)
	}
}

func runRunner(runner []string) {
	logger.Std.Print("new: ", shellquote.Join(runner...))

	execCmd := exec.Command(runner[0], runner[1:]...)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	err := execCmd.Run()
	if err != nil {
		logger.Err.Fatal(err)
	}
}

// NewCmd offer a command for create a new note in buffer directory.
var NewCmd = &cobra.Command{
	Use:   "new <title>.<type> <runner>",
	Short: "Do whatever swiftly on a non-existing Simple Note",
	Long: `Do whatever swiftly on a non-existing Simple Note.

  This command offer a quick method to operate a non-existing Simple Note
  without assign a path explicitly.

  This command will make sure the note has a valid note filename, and
  really not exists, or else runner will not be run.

  If run, the real file path will be auto assigned in the "buffer"
  directory that can be configured by configuration file.

  If buffer directory out of limit, this command will also reject to run.
  User should clear up the buffer directory (rearrange data to the other
  places) first by themself.

  Basically, this command was designed for create a simple note quickly &
  safely. Especially user won't to consider the classification right now.

  Example:
    $ loli new game-play-log.txt touch
    $ loli new how-to-cast-magic.md vim
`,
	Run: func(cmd *cobra.Command, args []string) {
		verifyArgs(args)

		inRepoPath, _ := cmd.Flags().GetString("repo")
		cfg := config.New(env.Env, inRepoPath, true)
		bPath := filepath.Join(cfg.Repo().Path(), cfg.NewBuffer())
		prepareBufferArea(bPath)

		inForce, _ := cmd.Flags().GetBool("force")
		bufferSize := cfg.NewBufferSize()
		verifyBufferAreaSize(bPath, inForce, bufferSize)

		nPath := filepath.Join(bPath, args[0])
		verifyNotePath(nPath)

		runner := args[1:]
		runner = append(runner, nPath)
		runRunner(runner)
	},
}

func init() {
	NewCmd.Flags().StringP(
		"repo", "r", "",
		"assign a Lolinote repository",
	)
	NewCmd.Flags().BoolP(
		"force", "f", false,
		"force add a note even exceeding buffer size limit",
	)
}
