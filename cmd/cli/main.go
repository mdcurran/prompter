package main

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/mdcurran/prompter/internal/pkg/redis"
)

var (
	nouns      int64
	verbs      int64
	adjectives int64
)

var cliCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		prompt := promptui.Prompt{
			Label: "prompter",
		}

		for {
			w, err := fetch()
			if err != nil {
				fmt.Printf("error getting prompts: %s", err)
				os.Exit(1)
			}
			fmt.Println(w)

			s, err := prompt.Run()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = redis.Record(s)
			if err != nil {
				fmt.Printf("error recording sentence: %s", err)
				os.Exit(1)
			}
		}
	},
}

func fetch() (map[string][]string, error) {
	w := make(map[string][]string, 1)

	n, err := redis.Get("noun", nouns)
	if err != nil {
		return nil, err
	}

	v, err := redis.Get("verb", verbs)
	if err != nil {
		return nil, err
	}

	a, err := redis.Get("adjective", adjectives)
	if err != nil {
		return nil, err
	}

	w["noun"] = n
	w["verb"] = v
	w["adjective"] = a

	return w, nil
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
