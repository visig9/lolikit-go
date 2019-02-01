package main

import (
	"github.com/spf13/cobra"

	"github.com/visig9/lolikit-go/cmd/initial"
	"github.com/visig9/lolikit-go/cmd/jsoncmd"
	"github.com/visig9/lolikit-go/cmd/list"
	"github.com/visig9/lolikit-go/cmd/newnote"
	"github.com/visig9/lolikit-go/cmd/serve"
)

var version string

var rootCmd = &cobra.Command{
	Short: "Lolikit - Lolinote 2.0 toolkit",
	Long: `Lolikit - Lolinote 2.0 toolkit

  The Lolinote 2.0 is a simple data specification for personal note-taking.

  Lolikit offer some extra conveniences to help users to manage their
  data for daily noting.

  Source Code: https://github.com/visig9/lolikit-go
  Lolinote Spec: https://github.com/visig9/lolinote-spec
	`,
	Version: version,
}

func init() {
	rootCmd.AddCommand(list.ListCmd)
	rootCmd.AddCommand(initial.InitCmd)
	rootCmd.AddCommand(serve.ServeCmd)
	rootCmd.AddCommand(jsoncmd.JSONCmd)
	rootCmd.AddCommand(newnote.NewCmd)
}

func main() {
	// defer func() { // panic processor
	// 	if err := recover(); err != nil {
	// 		logger.Err.Fatal(err)
	// 	}
	// }()
	rootCmd.Execute()
}
