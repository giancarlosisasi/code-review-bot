package main

import (
	code_review_bot "github.com/giancarlosisasi/code-review-bot"
	"github.com/rs/zerolog/log"
)

func main() {
	err := code_review_bot.RunApplication()
	if err != nil {
		log.Fatal().Err(err).Msg("Error to run the main application.")
	}
}
