package slackbot

import (
	"fmt"
	originalLog "log"
	"os"

	"github.com/giancarlosisasi/code-review-bot/config"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type SlackBot struct {
	SlackClient *socketmode.Client
}

func CreateSlackBot(config *config.Config) *SlackBot {
	slackApi := slack.New(
		config.SlackBotOauthToken,
		slack.OptionDebug(true),
		slack.OptionLog(originalLog.New(os.Stdout, "slack-api: ", originalLog.Lshortfile|originalLog.LstdFlags)),
		slack.OptionAppLevelToken(config.SlackSocketModeToken),
	)

	slackClient := socketmode.New(
		slackApi,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(originalLog.New(os.Stdout, "slack-socketmode: ", originalLog.Lshortfile|originalLog.LstdFlags)),
	)

	go func() {
		for evt := range slackClient.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				log.Info().Msg("Connecting to slack with socket mode....")
			case socketmode.EventTypeConnectionError:
				log.Error().Msg("Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				log.Info().Msg("Connected to slack with Socket mode.")
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					log.Info().Msgf("Ignored %+v", evt)
					continue
				}

				log.Debug().Msgf("Event received: %+v", eventsAPIEvent)

				slackClient.Ack(*evt.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent

					switch ev := innerEvent.Data.(type) {
					case *slackevents.AppMentionEvent:
						_, _, err := slackClient.PostMessage(ev.Channel, slack.MsgOptionText("ok, working on it...", false))
						if err != nil {
							log.Err(err).Msg("failed posting message")
						}
					case *slackevents.MemberJoinedChannelEvent:
						log.Info().Msgf("User %q joined to channel %q", ev.User, ev.Channel)
					}
				default:
					slackClient.Debugf("unsupported Events API event received")
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

					slackClient.Debugf("button clicked!")
				case slack.InteractionTypeShortcut:
				case slack.InteractionTypeViewSubmission:
					// See https://api.slack.com/apis/connections/socket-implement#modal
				case slack.InteractionTypeDialogSubmission:
				default:

				}

				slackClient.Ack(*evt.Request, payload)
			case socketmode.EventTypeSlashCommand:
				cmd, ok := evt.Data.(slack.SlashCommand)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)

					continue
				}

				slackClient.Debugf("Slash command received: %+v", cmd)

				_, _, err := slackClient.PostMessage(cmd.ChannelID, slack.MsgOptionText("ok, working on it2...", false))
				if err != nil {
					log.Err(err).Msg("failed posting message")
				}

				// payload := map[string]interface{}{
				// 	"blocks": []slack.Block{
				// 		slack.NewSectionBlock(
				// 			&slack.TextBlockObject{
				// 				Type: slack.MarkdownType,
				// 				Text: "foo",
				// 			},
				// 			nil,
				// 			slack.NewAccessory(
				// 				slack.NewButtonBlockElement(
				// 					"",
				// 					"somevalue",
				// 					&slack.TextBlockObject{
				// 						Type: slack.PlainTextType,
				// 						Text: "bar",
				// 					},
				// 				),
				// 			),
				// 		),
				// 	},
				// }
				payload := map[string]interface{}{
					"text": "Command received!",
				}
				// its important to provide a proper response for this
				slackClient.Ack(*evt.Request, payload)

			case socketmode.EventTypeHello:
				slackClient.Debugf("Hello received!")
			default:
				fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)

			}
		}
	}()

	return &SlackBot{
		SlackClient: slackClient,
	}
}

func (s *SlackBot) Run() error {
	err := s.SlackClient.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("Error to `run` the slack client module")
	}

	return nil
}
