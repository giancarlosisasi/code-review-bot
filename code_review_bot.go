package code_review_bot

import (
	"github.com/giancarlosisasi/code-review-bot/config"
	"github.com/giancarlosisasi/code-review-bot/database"
	"github.com/giancarlosisasi/code-review-bot/gitlab_client"
	"github.com/giancarlosisasi/code-review-bot/repository"
	"github.com/giancarlosisasi/code-review-bot/server"
	slackbot "github.com/giancarlosisasi/code-review-bot/slackbot"
	"github.com/rs/zerolog/log"
)

type App struct {
	config *config.Config
	// dbConn *pgxpool.Pool
}

func NewApp(
	config *config.Config,
	// dbConn *pgxpool.Pool,
) *App {
	return &App{
		config: config,
		// dbConn: dbConn,
	}
}

func RunApplication() error {
	// config
	config := config.NewConfig()

	// for now we are not going to to use a db
	// dbpool, err := database.NewDBConn(config.DBUrl)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("Failed to connect to the database")
	// }

	// in memory database
	inMemoryDatabase := database.NewInMemoryDatabase()

	// app := NewApp(config)

	// Channel to wait for both services
	done := make(chan error, 2)

	gitlabClient, err := gitlab_client.CreateGitlabClient(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Error to create the gitlab client")
	}

	// repositories
	teamMemberRepository := repository.NewTeamMembersInMemoryRepository(
		inMemoryDatabase,
	)
	reviewRepository := repository.NewReviewInMemoryRepository(inMemoryDatabase)

	go func() {
		server := server.NewServer(inMemoryDatabase, config)
		log.Info().Msgf("Starting web server on http://localhost:%d", config.Port)
		err := server.Run()
		if err != nil {
			log.Error().Err(err).Msg("failed to run the web server")
		}
		done <- err
	}()

	go func() {
		// setup slack bot
		slackBot := slackbot.CreateSlackBot(config, gitlabClient, teamMemberRepository, reviewRepository)
		log.Info().Msg("Starting Slack bot")
		err := slackBot.Run()
		if err != nil {
			log.Error().Err(err).Msg("Error to `run` the slack client module")
		}
		done <- err
	}()

	// Wait for both services to complete (or fail)
	for range 2 {
		if err := <-done; err != nil {
			return err
		}
	}

	return nil
}
