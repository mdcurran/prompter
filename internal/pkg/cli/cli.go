package cli

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"

	"github.com/mdcurran/prompter/internal/pkg/redis"
)

// state holds the number of tokens to be returned for each part of speech for prompts in the
// current session.
type state struct {
	nouns      int64
	verbs      int64
	adjectives int64
}

// Run begins the prompter command-line application. Prompts will be continuously generated unless
// of an error (such as storage becoming unavailable) or the program is terminated.
func Run(nouns, verbs, adjectives int64) {
	s := &state{
		nouns:      nouns,
		verbs:      verbs,
		adjectives: adjectives,
	}

	prompt := promptui.Prompt{
		Label: "prompter",
	}

	for {
		w, err := fetch(s)
		if err != nil {
			fmt.Printf("error getting prompts: %s", err)
			os.Exit(1)
		}
		fmt.Println(w)

		sentence, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = redis.Record(sentence)
		if err != nil {
			fmt.Printf("error recording sentence: %s", err)
			os.Exit(1)
		}
	}
}

// fetch retrieves a number of randomised tokens of varying parts of speech based on the
// application state.
func fetch(s *state) (map[string][]string, error) {
	w := make(map[string][]string)

	n, err := redis.Get("noun", s.nouns)
	if err != nil {
		return nil, err
	}

	v, err := redis.Get("verb", s.verbs)
	if err != nil {
		return nil, err
	}

	a, err := redis.Get("adjective", s.adjectives)
	if err != nil {
		return nil, err
	}

	w["noun"] = n
	w["verb"] = v
	w["adjective"] = a

	return w, nil
}
