package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Port                 int
	AppEnv               string
	DBUrl                string
	GitlabAPIToken       string
	GitlabOrgSlug        string
	SlackSocketModeToken string
	SlackBotOauthToken   string
	SlackSigningSecret   string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Warn().Msg("Error to load the .env file.")
	}

	viper.AutomaticEnv()

	port := mustGetInt("PORT")
	dbUrl := mustGetString("DB_URL")
	appEnv := mustGetString("APP_ENV")
	gitlabApiToken := mustGetString("GITLAB_API_TOKEN")
	gitlabOrgSlug := mustGetString("GITLAB_ORG_SLUG")

	slackSocketToken := mustGetString("SLACK_SOCKET_MODE_TOKEN")
	if !strings.HasPrefix(slackSocketToken, "xapp-") {
		log.Fatal().Msg("SLACK_SOCKET_MODE_TOKEN must have the prefix \"xoxb-\"\n")
	}

	slackBotToken := mustGetString("SLACK_BOT_OAUTH_TOKEN")
	if !strings.HasPrefix(slackBotToken, "xoxb-") {
		log.Fatal().Msg("SLACK_BOT_OAUTH_TOKEN must have the prefix \"xoxb-\"\n")
	}

	slackSigninSecret := mustGetString("SLACK_SIGNING_SECRET")

	return &Config{
		Port:                 port,
		AppEnv:               appEnv,
		DBUrl:                dbUrl,
		GitlabAPIToken:       gitlabApiToken,
		GitlabOrgSlug:        gitlabOrgSlug,
		SlackSocketModeToken: slackSocketToken,
		SlackBotOauthToken:   slackBotToken,
		SlackSigningSecret:   slackSigninSecret,
	}
}

func mustGetInt(key string) int {
	v := viper.GetInt(key)
	if v == 0 {
		log.Fatal().Msgf("The env '%s' must be set.", key)
	}

	return v
}

func mustGetString(key string) string {
	v := viper.GetString(key)
	if v == "" {
		log.Fatal().Msgf("The env '%s' must be set.", key)
	}

	return v
}
