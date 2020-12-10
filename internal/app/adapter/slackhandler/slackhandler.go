// Package slackhandler provides features that control slack API.
package slackhandler

import (
	"slacktimer/internal/app/driver/slack"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/di"
)

// SlackHandler serves Slack API handlers used by the app.
type SlackHandler struct {
	api slack.API
}

// NewSlackHandler create new struct.
func NewSlackHandler() *SlackHandler {
	s := &SlackHandler{
		api: di.Get("slack.API").(slack.API),
	}
	return s
}

// Notify notify message to user identified by userID.
func (s SlackHandler) Notify(ac appcontext.AppContext, userID string, message string) error {
	// Need to open DM channel to send DM.
	channelID, err := s.api.ConversationsOpen(ac, userID)
	if err != nil {
		return err
	}

	// Send DM.
	err = s.api.ChatPostMessage(ac, channelID, message)
	if err != nil {
		return err
	}

	return nil
}
