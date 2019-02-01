// Package list show the list of note with particular order.
package list

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/spf13/cobra"

	"github.com/visig9/lolikit-go/cmd/config"
	"github.com/visig9/lolikit-go/cmd/env"
	"github.com/visig9/lolikit-go/cmd/list/nps"
	"github.com/visig9/lolikit-go/cmd/list/runnote"
	"github.com/visig9/lolikit-go/logger"
	"github.com/visig9/lolikit-go/loli2"
	"github.com/visig9/lolikit-go/parse"
	"github.com/visig9/pager"
)

func assertGoodInput(inAttMode, inDirMode, inEntryMode bool) {
	var count int
	if inAttMode {
		count++
	}
	if inDirMode {
		count++
	}
	if inEntryMode {
		count++
	}

	if count >= 2 {
		logger.Err.Fatal("flag -d, -a, -e are mutual exclusion")
	}
}

func isPrintMode(args []string) bool {
	if len(args) == 0 {
		return true
	}

	return false
}

func assertGoodPrintModeInput(inAttMode, inDirMode, inEntryMode bool) {
	if inAttMode {
		logger.Err.Fatal(
			"flag --attachment-dir are useless without runid",
		)
	}

	if inDirMode {
		logger.Err.Fatal(
			"flag --parent-dir are useless without runid",
		)
	}

	if inEntryMode {
		logger.Err.Fatal("flag --entry are useless without runid")
	}
}

func getRunModeData(
	args []string, inAttMode, inDirMode, inEntryMode bool,
) (
	runid int, inRunner []string, runType runnote.Type,
) {
	if len(args) >= 1 {
		var err error
		runid, err = parse.Int(args[0])
		if err != nil {
			logger.Err.Fatal(err)
		}
	}

	if len(args) >= 2 {
		inRunner = args[1:]
	}

	switch {
	case inAttMode:
		runType = runnote.RunAttachmentDir
	case inDirMode:
		runType = runnote.RunDir
	case inEntryMode:
		runType = runnote.RunEntry
	default:
		runType = runnote.RunContent
	}

	return
}

func getNotes(
	inTerms []string, cfg *config.Cfg, repo loli2.Repo,
) *arraylist.List {
	var notes *arraylist.List
	if len(inTerms) > 0 {
		textTypes := cfg.TextTypes()
		notes = repo.NotesByRel(inTerms, textTypes)
	} else {
		notes = repo.NotesByMod()
	}

	return notes
}

func getPage(
	inPageSize, inPageNum int, cfg *config.Cfg, items []interface{},
) pager.Page {
	return pager.Pager{
		Items:    items,
		PageSize: cfg.ListPageSize(inPageSize),
	}.Page(inPageNum)
}

func run(cmd *cobra.Command, args []string) {
	inRepoPath, _ := cmd.Flags().GetString("repo")
	cfg := config.New(env.Env, inRepoPath, true)

	repo := cfg.Repo()
	inTerms, _ := cmd.Flags().GetStringSlice("term")
	notes := getNotes(inTerms, cfg, repo)

	inPageSize, _ := cmd.Flags().GetInt("page-size")
	inPageNum, _ := cmd.Flags().GetInt("page")
	page := getPage(inPageSize, inPageNum, cfg, notes.Values())

	if isPrintMode(args) {
		if page.Size() > 0 {
			logger.Std.Print(nps.ToString(page))
		}
	} else {
		inAttMode, _ := cmd.Flags().GetBool("attachment-dir")
		inDirMode, _ := cmd.Flags().GetBool("parent-dir")
		inEntryMode, _ := cmd.Flags().GetBool("entry")
		runid, inRunner, runType := getRunModeData(
			args, inAttMode, inDirMode, inEntryMode,
		)

		runnote.Run(page, runid, inRunner, cfg, runType)
	}
}

// ListCmd is one of a sub command of lolikit-go
var ListCmd = &cobra.Command{
	Use:   "list [# runid] [runner]", // "\n   list --run-all [runner]",
	Short: "List notes in particular order",
	Long: `List notes in particular order

  === Listing ===

  By default, the list is ordered by note's modify time.
  If given the "-t" option, it will using the relevant order.

  Listing Example:
    $ loli list
    $ loli list -t <term>

  === Run Command ===

  After some listing, user can also choice the "runid" and "runner" to
  open (or processing) a particular note.

  Running Example:
    $ loli list 3 ls        # run "ls <note-3>"
    $ loli list 1 vim       # run "vim <note-1>"
    $ loli list -p 2 1 cat  # run cat with <the first note in page 2>

  User can omit the "runner" if already assigned in configuration file.

  By default, the *Run-Path* will be appended to the tail of runner. But
  it can also be assigned to particular location by using the placeholder
  "{}".

  === The Run-Path ===

  List command also offer multiple type of run-path.

  The run-path Type:
    1. Content Path: the filepath contain the note's content. Good for
       editing or viewing.
    2. Parent Directory: the dir path contain the note's entry. Good for
       files browsing.
    3. Attachment Directory: the dir path contain the complex-note's
       attachments. Good for tweak the attachments.
    4. Entry Path: the entry path of the notes. Good for deleting or
       moving.

  By default, list command using Content Path as run-path.
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		inAttMode, _ := cmd.Flags().GetBool("attachment-dir")
		inDirMode, _ := cmd.Flags().GetBool("parent-dir")
		inEntryMode, _ := cmd.Flags().GetBool("entry")

		assertGoodInput(inAttMode, inDirMode, inEntryMode)

		if isPrintMode(args) {
			assertGoodPrintModeInput(
				inAttMode, inDirMode, inEntryMode,
			)
		}
	},
	Run: run,
}

func init() {
	ListCmd.Flags().StringP(
		"repo", "r", "",
		"assign a Lolinote repository",
	)
	ListCmd.Flags().IntP(
		"page", "p", 1,
		"page number",
	)
	ListCmd.Flags().IntP(
		"page-size", "s", 0,
		"page size",
	)
	ListCmd.Flags().BoolP(
		"attachment-dir", "a", false,
		"using Attachment Directory as run-path and set\n"+
			"dir-runner as default runner",
	)
	ListCmd.Flags().BoolP(
		"parent-dir", "d", false,
		"using Parent Directory as run-path and set\n"+
			"dir-runner as default runner",
	)
	ListCmd.Flags().BoolP(
		"entry", "e", false,
		"using Entry Path as run-path.",
	)
	ListCmd.Flags().StringSliceP(
		"term", "t", []string{},
		"give some terms and enable relevance ordering",
	)
}
