// Package jsoncmd print the all note in json format.
package jsoncmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/visig/lolikit-go/cmd/config"
	"gitlab.com/visig/lolikit-go/cmd/env"
	"gitlab.com/visig/lolikit-go/logger"
	"gitlab.com/visig/lolikit-go/loli2"
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
