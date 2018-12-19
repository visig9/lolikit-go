// Package serve active a http server for the Lolinote repo.
package serve

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"gitlab.com/visig/lolikit-go/cmd/config"
	"gitlab.com/visig/lolikit-go/cmd/env"
	"gitlab.com/visig/lolikit-go/loli2"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>ok</h1><div>Hi there, I love %s!</div>", r.URL.Path[1:])
}

func routing(repo loli2.Repo) {
	http.HandleFunc("/", handler)

	http.Handle(
		"/raw/",
		http.StripPrefix(
			"/raw/",
			http.FileServer(&repoFileSystem{repo}),
		),
	)

	http.Handle(
		"/view/",
		http.StripPrefix(
			"/view/",
			newViewHandler(repo),
		),
	)
}

// ServeCmd offer a command for initialize a lolinote repo.
var ServeCmd = &cobra.Command{
	Use:   "serve [path]",
	Short: "serve a Lolinote repository",
	Long:  `serve a Lolinote repository`,
	Run: func(cmd *cobra.Command, args []string) {
		inRepoPath, _ := cmd.Flags().GetString("repo")
		cfg := config.New(env.Env, inRepoPath, true)

		repo := cfg.Repo()

		routing(repo)
		log.Fatal(http.ListenAndServe(":8080", nil))
	},
}

func init() {
	ServeCmd.Flags().StringP(
		"repo", "r", "",
		"assign a Lolinote repository",
	)
}
