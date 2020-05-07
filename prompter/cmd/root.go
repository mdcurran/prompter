package cmd

import (
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/mdcurran/prompter/internal/pkg/server"
)

var rootCmd = &cobra.Command{
	Use:   "prompter",
	Short: "prompter helps you learn to word good!",
	Run: func(cmd *cobra.Command, args []string) {
		err := logger()
		if err != nil {
			panic(err)
		}

		s, err := server.New()
		if err != nil {
			zap.S().Errorf("unable to instantiate Server: %s", err)
			os.Exit(1)
		}

		zap.S().Info("starting HTTP server")
		err = http.ListenAndServe(":8080", s.Router)
		if err != nil {
			zap.S().Errorf("unable to initalise HTTP server: %s", err)
			os.Exit(1)
		}
	},
}

// logger instantiates an instance of the Zap logger with recommended defaults and replaces the
// global logger.
func logger() error {
	l, err := zap.NewProduction()
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(l)
	return nil
}

// Execute is the main application entrypoint.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
