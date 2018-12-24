package main

import (
	"github.com/spf13/cobra"

	"gitlab.com/visig/lolikit-go/cmd/initial"
	"gitlab.com/visig/lolikit-go/cmd/jsoncmd"
	"gitlab.com/visig/lolikit-go/cmd/list"
	"gitlab.com/visig/lolikit-go/cmd/serve"
)

var version string

var rootCmd = &cobra.Command{
	Short: "Lolikit - Lolinote 2.0 toolkit",
	Long: `Lolikit - Lolinote 2.0 toolkit

  The Lolinote 2.0 is a simple data specification for personal note-taking.

  Lolikit offer some extra conveniences to help users to manage their
  data for daily noting.

  Source Code: https://gitlab.com/visig/lolikit-go
  Lolinote Spec: https://gitlab.com/visig/lolinote-spec
	`,
	Version: version,
}

func init() {
	rootCmd.AddCommand(list.ListCmd)
	rootCmd.AddCommand(initial.InitCmd)
	rootCmd.AddCommand(serve.ServeCmd)
	rootCmd.AddCommand(jsoncmd.JSONCmd)
}

func main() {
	// defer func() { // panic processor
	// 	if err := recover(); err != nil {
	// 		logger.Err.Fatal(err)
	// 	}
	// }()
	rootCmd.Execute()
}
