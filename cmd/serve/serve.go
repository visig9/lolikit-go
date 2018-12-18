// Package serve active a http server for the Lolinote repo.
package serve

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>ok</h1><div>Hi there, I love %s!</div>", r.URL.Path[1:])
}

// ServeCmd offer a command for initialize a lolinote repo.
var ServeCmd = &cobra.Command{
	Use:   "serve [path]",
	Short: "serve a Lolinote repository",
	Long:  `serve a Lolinote repository`,
	Run: func(cmd *cobra.Command, args []string) {
		http.HandleFunc("/", handler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	},
}
