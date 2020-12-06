// Package slackhandler provides features that control slack API.
package slackhandler

import (
	"slacktimer/internal/app/driver/slack"
	"slacktimer/internal/app/util/di"
)

type SlackHandler struct {
	api slack.SlackApi
}

func NewSlackHandler() *SlackHandler {
	s := &SlackHandler{
		api: di.Get("slack.SlackApi").(slack.SlackApi),
	}
	return s
}

func (s SlackHandler) Notify(userId string, message string) error {
	// Need to open DM channel to send DM.
	channelId, err := s.api.ConversationsOpen(userId)
	if err != nil {
		return err
	}

	// Send DM.
	err = s.api.ChatPostMessage(channelId, message)
	if err != nil {
		return err
	}

	return nil
}
