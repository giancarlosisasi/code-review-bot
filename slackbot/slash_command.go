package slackbot

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/giancarlosisasi/code-review-bot/utils"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

const MAX_NUMBER_OF_REVIEWERS = 2
const mergeRequestStateOpened = "opened"

var teamGuildFrontend = "frontend"
var teamGuildBackend = "backend"
var teamGuildDevops = "devops"
var teamGuildQa = "qa"

func (sb *SlackBot) handleCodeReviewSlackCommand(cmd slack.SlashCommand, event socketmode.Event) error {
	msg := strings.Split(cmd.Text, " ")
	var mrUrl string

	for _, value := range msg {
		if strings.Contains(value, "gitlab.com") && strings.Contains(value, sb.config.GitlabOrgSlug) {
			mrUrl = value
		}
	}

	_, _, err := sb.slackClient.PostMessage(cmd.ChannelID, slack.MsgOptionText(
		"Alright, i'm working on it :runner:",
		false,
	))
	if err != nil {
		log.Err(err).Msg("failed posting message")
		sb.slackClient.Ack(*event.Request, nil)
		return nil
	}

	// 1. Make a fetch to get the gitlab merge request data
	// The gitlab api expects two param ids
	// 	- id: repository path encoded, for example: nexus-team%2Ffrontend%2Ffrontend-super-project
	//	- merge_request_iid, usually the last number part of the url like 1790
	// we should get the mr iid, full url, author email and assignees emails
	// 	- the response doesn't return a user email, instead we should use the gitlab user.id and search in our in-memory database
	parsedUrl, err := utils.ParseGitlabURL(mrUrl)
	if err != nil {
		_, _, err := sb.slackClient.PostMessage(cmd.ChannelID, slack.MsgOptionText(
			err.Error(),
			false,
		))

		sb.slackClient.Ack(*event.Request, nil)

		if err != nil {
			log.Err(err).Msg("failed to post error parsed gitlab url")
			return err
		}

		return nil
	}

	mergeRequest, res, err := sb.gitlabClient.MergeRequests.GetMergeRequest(
		// parsedUrl.ProjectPath,
		"59603690",
		parsedUrl.MergeRequestIID,
		nil,
	)

	if err != nil {
		_, _, err = sb.slackClient.PostMessage(cmd.ChannelID, slack.MsgOptionText(
			fmt.Sprintf("Ups! there was an error when trying to access to the gitlab api: %s", err.Error()),
			false,
		))

		sb.slackClient.Ack(*event.Request, nil)

		if err != nil {
			log.Err(err).Msg("failed to post message")
			return err
		}

		return nil
	}

	if mergeRequest == nil || res.StatusCode == http.StatusNotFound {
		_, _, err = sb.slackClient.PostMessage(cmd.ChannelID, slack.MsgOptionText(
			fmt.Sprintf("The merge request with url `%s`. We are using the %s and %d", mrUrl, parsedUrl.ProjectPath, parsedUrl.MergeRequestIID),
			false,
		))

		sb.slackClient.Ack(*event.Request, nil)

		if err != nil {
			log.Err(err).Msg("failed to post message")
			return err
		}

		return nil
	}

	if mergeRequest.State != mergeRequestStateOpened {
		_, _, err = sb.slackClient.PostMessage(cmd.ChannelID, slack.MsgOptionText(
			fmt.Sprintf("The merge request with url `%s`. We are using the %s and %d", mrUrl, parsedUrl.ProjectPath, parsedUrl.MergeRequestIID),
			false,
		))

		sb.slackClient.Ack(*event.Request, nil)

		if err != nil {
			log.Err(err).Msg("failed to post message")
			return err
		}

		return nil
	}

	var reviewersGitlabMemberIDs []string
	for _, reviewer := range mergeRequest.Reviewers {
		reviewersGitlabMemberIDs = append(reviewersGitlabMemberIDs, strconv.Itoa(reviewer.ID))
	}

	// verify that the MR has not been yet assigned
	reviewerMembers := sb.teamMembersRepository.FindTeamMembersByGuild(
		teamGuildFrontend,
		reviewersGitlabMemberIDs,
	)

	// it means that the MR has already assigned members, in this case just notify who are the assigned to this MR
	if len(reviewerMembers) >= MAX_NUMBER_OF_REVIEWERS {
		slackUserTags := []string{}
		for _, m := range reviewerMembers {
			slackUserTags = append(slackUserTags, fmt.Sprintf("<@%s>", m.SlackMemberID))
		}

		_, _, err = sb.slackClient.PostMessage(cmd.ChannelID, slack.MsgOptionText(
			fmt.Sprintf("This MR already has %d reviewers. Ask to %s if you have questions.", MAX_NUMBER_OF_REVIEWERS, strings.Join(slackUserTags, ", ")),
			false,
		))
		if err != nil {
			log.Err(err).Msg("failed to post message")
			return err
		}

		sb.slackClient.Ack(*event.Request, nil)
		return nil
	}

	// 2. Run the algorithm to get two members to be assigned to the current mr excluding the author ID
	// membersToExcludeIDs := append(reviewersGitlabMemberIDs, strconv.Itoa(mergeRequest.Author.ID))
	// membersAvailable := sb.teamMembersRepository.GetTeamMembersExcludingMembers(membersToExcludeIDs, &teamGuildFrontend)

	// 3. Use the gitlab api to update the MR and assign the two members

	// 4. Answer or sent back a message to the same channel/thread notifying the slack user that the MR has been assigned to two persons
	// use the slack member ids to tag the reviewers so they can be notified too

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
		"text": "Ok, let me work on this! :runner:",
	}
	// its important to provide a proper response for this
	sb.slackClient.Ack(*event.Request, payload)

	return nil
}
