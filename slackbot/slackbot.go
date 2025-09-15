package slackbot

import (
	"fmt"
	originalLog "log"
	"os"

	"github.com/giancarlosisasi/code-review-bot/config"
	"github.com/giancarlosisasi/code-review-bot/repository"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type SlackBot struct {
	slackClient           *socketmode.Client
	config                *config.Config
	gitlabClient          *gitlab.Client
	teamMembersRepository repository.TeamMembersRepository
	reviewRepository      repository.ReviewRepository
}

func CreateSlackBot(config *config.Config, gitlabClient *gitlab.Client, tmRepo repository.TeamMembersRepository, rp repository.ReviewRepository) *SlackBot {
	slackApi := slack.New(
		config.SlackBotOauthToken,
		slack.OptionDebug(true),
		slack.OptionLog(originalLog.New(os.Stdout, "[SLACK API]: ", originalLog.Lshortfile|originalLog.LstdFlags)),
		slack.OptionAppLevelToken(config.SlackSocketModeToken),
	)

	slackClient := socketmode.New(
		slackApi,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(originalLog.New(os.Stdout, "[SLACK SOCKET MODE]: ", originalLog.Lshortfile|originalLog.LstdFlags)),
	)

	return &SlackBot{
		slackClient:           slackClient,
		config:                config,
		gitlabClient:          gitlabClient,
		teamMembersRepository: tmRepo,
		reviewRepository:      rp,
	}
}

func (sb *SlackBot) Run() error {
	go func() {
		for evt := range sb.slackClient.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				// log.Info().Msg("Connecting to slack with socket mode....")
				continue
			case socketmode.EventTypeConnectionError:
				log.Error().Msg("Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				// log.Info().Msg("Connected to slack with Socket mode.")
				continue
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					log.Info().Msgf("Ignored %+v", evt)
					continue
				}

				log.Debug().Msgf("Event received: %+v", eventsAPIEvent)

				sb.slackClient.Ack(*evt.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent

					switch ev := innerEvent.Data.(type) {
					case *slackevents.AppMentionEvent:
						_, _, err := sb.slackClient.PostMessage(ev.Channel, slack.MsgOptionText("ok, working on it...", false))
						if err != nil {
							log.Err(err).Msg("failed posting message")
						}
					case *slackevents.MemberJoinedChannelEvent:
						log.Info().Msgf("User %q joined to channel %q", ev.User, ev.Channel)
					}
				default:
					sb.slackClient.Debugf("unsupported Events API event received")
				}
			case socketmode.EventTypeInteractive:
				callback, ok := evt.Data.(slack.InteractionCallback)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)

					continue
				}

				fmt.Printf("Interaction received: %+v\n", callback)

				var payload interface{}

				switch callback.Type {
				case slack.InteractionTypeBlockActions:
					// See https://api.slack.com/apis/connections/socket-implement#button

					sb.slackClient.Debugf("button clicked!")
				case slack.InteractionTypeShortcut:
				case slack.InteractionTypeViewSubmission:
					// See https://api.slack.com/apis/connections/socket-implement#modal
				case slack.InteractionTypeDialogSubmission:
				default:

				}

				sb.slackClient.Ack(*evt.Request, payload)

			case socketmode.EventTypeSlashCommand:
				cmd, ok := evt.Data.(slack.SlashCommand)
				if !ok {
					log.Debug().Msgf("Slash event data ignored: %+v", evt)
					continue
				}

				log.Debug().Msgf("Command is: %s", cmd.Command)

				switch cmd.Command {
				case "/code-review":
					err := sb.handleCodeReviewSlackCommand(cmd, evt)
					if err != nil {
						continue
					}
				default:
					continue
				}

			case socketmode.EventTypeHello:
				sb.slackClient.Debugf("Hello received!")
			default:
				fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)

			}
		}
	}()

	err := sb.slackClient.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("Error to `run` the slack client module")
	}

	return nil
}
