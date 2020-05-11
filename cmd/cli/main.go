package main

import (
	"github.com/spf13/cobra"

	"github.com/mdcurran/prompter/internal/pkg/cli"
)

var (
	nouns      int64
	verbs      int64
	adjectives int64
)

var cliCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		cli.Run(nouns, verbs, adjectives)
	},
}

func main() {
	cliCmd.Flags().Int64Var(&nouns, "nouns", 1, "Number of nouns in prompt")
	cliCmd.Flags().Int64Var(&verbs, "verbs", 1, "Number of verbs in prompt")
	cliCmd.Flags().Int64Var(&adjectives, "adjectives", 1, "Number of adjectives in prompt")

	err := cliCmd.Execute()
	if err != nil {
		panic(err)
	}
}
