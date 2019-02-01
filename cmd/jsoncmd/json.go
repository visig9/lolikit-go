// Package jsoncmd print the all note in json format.
package jsoncmd

import (
	"github.com/spf13/cobra"
	"github.com/visig9/lolikit-go/cmd/config"
	"github.com/visig9/lolikit-go/cmd/env"
	"github.com/visig9/lolikit-go/logger"
	"github.com/visig9/lolikit-go/loli2"
)

// JSONCmd offer a command for print all note data in json format.
var JSONCmd = &cobra.Command{
	Use:   "json",
	Short: "Print all information of notes in JSON format",
	Long:  "Print all information of notes in JSON format",
	Run: func(cmd *cobra.Command, args []string) {
		inRepoPath, _ := cmd.Flags().GetString("repo")
		cfg := config.New(env.Env, inRepoPath, true)

		cfg.Repo().Dir().Walk(func(e loli2.Entry) {
			if n, ok := e.(loli2.Note); ok {
				logger.Std.Print(string(n.JSON()))
			}
		})
	},
}

func init() {
	JSONCmd.Flags().StringP(
		"repo", "r", "",
		"assign a Lolinote repository",
	)
}
