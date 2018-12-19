// Package serve active a http server for the Lolinote repo.
package serve

import (
	"net/http"

	"github.com/spf13/cobra"
	"gitlab.com/visig/lolikit-go/cmd/config"
	"gitlab.com/visig/lolikit-go/cmd/env"
	"gitlab.com/visig/lolikit-go/logger"
)

// ServeCmd offer a command for initialize a lolinote repo.
var ServeCmd = &cobra.Command{
	Use:   "serve [path]",
	Short: "serve a Lolinote repository",
	Long:  `serve a Lolinote repository`,
	Run: func(cmd *cobra.Command, args []string) {
		inRepoPath, _ := cmd.Flags().GetString("repo")
		cfg := config.New(env.Env, inRepoPath, true)

		repo := cfg.Repo()

		inAddr, _ := cmd.Flags().GetString("addr")
		addr := cfg.ServeAddr(inAddr)

		http.Handle("/", http.FileServer(&repoFileSystem{repo}))

		logger.Std.Printf(
			"Lolikit listening on %v", addr,
		)
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
