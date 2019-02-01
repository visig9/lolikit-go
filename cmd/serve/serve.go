// Package serve active a http server for the Lolinote repo.
package serve

import (
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/visig9/lolikit-go/cmd/config"
	"github.com/visig9/lolikit-go/cmd/env"
	"github.com/visig9/lolikit-go/httpfs"
	"github.com/visig9/lolikit-go/logger"
)

func fsFilter(name string) bool {
	if strings.HasPrefix(name, ".") {
		return false
	}

	return true
}

func getFSHandler(root string) http.Handler {
	fs := http.FileServer(&httpfs.CondNameFileSystem{
		Root: root, Filter: fsFilter, Recu: true,
	})

	return fs
}

// ServeCmd offer a command for initialize a lolinote repo.
var ServeCmd = &cobra.Command{
	Use:   "serve [path]",
	Short: "Start a HTTP server for remote access",
	Long: `Start a HTTP server for remote access

  It will startup a HTTP web server to serve the assigned repository as
  a file server.

  All directories and files in repository will be showed, including the
  "noise" in Lolinote specification. But all hidden entries (file name
  starts with ".") and those sub-entries will be filtered out for
  security reason.

  NOTE: Existing path structures may change in the future.
`,
	Run: func(cmd *cobra.Command, args []string) {
		inRepoPath, _ := cmd.Flags().GetString("repo")
		cfg := config.New(env.Env, inRepoPath, true)

		repo := cfg.Repo()
		http.Handle("/", getFSHandler(repo.Path()))

		inAddr, _ := cmd.Flags().GetString("addr")
		addr := cfg.ServeAddr(inAddr)

		logger.Std.Printf("Lolikit listening on %v ...", addr)
		logger.Err.Fatal(http.ListenAndServe(addr, nil))
	},
}

func init() {
	ServeCmd.Flags().StringP(
		"repo", "r", "",
		"assign a Lolinote repository",
	)
	ServeCmd.Flags().StringP(
		"addr", "a", "",
		"http server's bind address",
	)
}
