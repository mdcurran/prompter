package cmd

import (
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/mdcurran/prompter/internal/pkg/server"
)

var rootCmd = &cobra.Command{
	Use:   "prompter",
	Short: "prompter helps you learn to word good!",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := server.New()
		if err != nil {
			panic(err)
		}

		s.Logger.Info("starting HTTP server")
		err = http.ListenAndServe(":8080", s.Router)
		if err != nil {
			s.Logger.Errorf("unable to initalise HTTP server: %s", err)
			os.Exit(1)
		}
	},
}

// Execute is the main application entrypoint.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
