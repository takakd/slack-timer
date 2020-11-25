// Gateway role
package slackhandler

import (
	"slacktimer/internal/app/driver/di"
	"slacktimer/internal/app/driver/slack"
)

type SlackHandler struct {
	api slack.SlackApi
}

func NewSlackHandler() *SlackHandler {
	s := &SlackHandler{
		api: di.Get("slackhandler.SlackApi").(slack.SlackApi),
	}
	return s
}

func (s *SlackHandler) Notify(userId string, message string) error {
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
