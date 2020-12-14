// Package slack is Slack API adapter.
package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/config"
	"slacktimer/internal/app/util/log"
	"slacktimer/internal/pkg/helper"
)

// API lists the Slack API using in the app.
type API interface {
	ConversationsOpen(ac appcontext.AppContext, userID string) (string, error)
	ChatPostMessage(ac appcontext.AppContext, body ChatPostMessageRequestBody) error
}

// APIDriver implements API interface.
type APIDriver struct {
}

var _ API = (*APIDriver)(nil)

// NewAPIDriver creates new struct.
func NewAPIDriver() *APIDriver {
	s := &APIDriver{}
	return s
}

// ConversationsOpenRequestBody defines request body of conversations.open API.
// Ref: https://api.slack.com/methods/conversations.open
type ConversationsOpenRequestBody struct {
	// Token is set to Bearer Header.
	Users string `json:"users"`
}

// ConversationsOpenResponseBody defines response body of conversations.open API.
type ConversationsOpenResponseBody struct {
	Ok      bool                                 `json:"ok"`
	Channel ConversationsOpenResponseBodyChannel `json:"channel,omitempty"`
	// Be set if the response is error
	Error string `json:"error,omitempty"`
}

// ConversationsOpenResponseBodyChannel defines a element of ConversationsOpenResponseBody.
type ConversationsOpenResponseBodyChannel struct {
	ID string `json:"id"`
}

// ConversationsOpen opens DM Slack channel.
// Ref: https://api.slack.com/methods/conversations.open
func (s APIDriver) ConversationsOpen(ac appcontext.AppContext, userID string) (string, error) {
	body := &ConversationsOpenRequestBody{
		userID,
	}
	url := config.MustGet("SLACK_API_URL_CONVERSATIONSOPEN")
	resp, err := postJSON(url, body)
	if err != nil {
		return "", err
	}

	ok := resp.StatusCode == http.StatusOK
	if !ok {
		return "", fmt.Errorf("request error slack conversations.open user_id=%s: %w", userID, err)
	}

	respBuf, err := helper.GetResponseBody(resp)
	if err != nil {
		return "", fmt.Errorf("response reading error slack conversations.open user_id=%s: %w", userID, err)
	}

	log.DebugWithContext(ac, "response body slack conversations.open", string(respBuf))

	var respBody ConversationsOpenResponseBody
	err = json.Unmarshal(respBuf, &respBody)
	if err != nil {
		return "", fmt.Errorf("unmarshal error slack conversations.open user_id=%s: %w", userID, err)
	}

	// It must be returned one ID because of sending one user ID.
	if !respBody.Ok {
		return "", fmt.Errorf("response NG slack conversations.open user_id=%s body=%v", userID, respBody)
	}

	return respBody.Channel.ID, nil
}

// ChatPostMessageRequestBody defines request body of chat.postMessage API.
// Ref: https://api.slack.com/methods/chat.postMessage
type ChatPostMessageRequestBody struct {
	// Token is set to Bearer Header.
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// ChatPostMessageResponseBody defines response body of chat.postMessage API.
// Ref: https://api.slack.com/methods/chat.postMessage
type ChatPostMessageResponseBody struct {
	// Define only need
	Ok bool `json:"ok"`
	// Be set if the response is error
	Error string `json:"error,omitempty"`
}

// ChatPostMessage send message to DM Slack channel.
// Ref: https://api.slack.com/methods/chat.postMessage
func (s APIDriver) ChatPostMessage(ac appcontext.AppContext, body ChatPostMessageRequestBody) error {
	url := config.MustGet("SLACK_API_URL_CHATPOSTMESSAGE")
	resp, err := postJSON(url, body)
	if err != nil {
		return err
	}

	ok := resp.StatusCode == http.StatusOK
	if !ok {
		return fmt.Errorf("request error slack chat.postMessage channel_id=%s message=%s: %w", body.Channel, body.Text, err)
	}

	respBuf, err := helper.GetResponseBody(resp)
	if err != nil {
		return fmt.Errorf("response reading error slack chat.postMessage channel_id=%s message=%s: %w", body.Channel, body.Text, err)
	}

	log.DebugWithContext(ac, "response body slack chat.postMessage", string(respBuf))

	var respBody ChatPostMessageResponseBody
	err = json.Unmarshal(respBuf, &respBody)
	if err != nil {
		return fmt.Errorf("unmarshal error slack chat.postMessage channel_id=%s message=%s", body.Channel, body.Text)
	}

	if !respBody.Ok {
		return fmt.Errorf("response NG slack chat.postMessage channel_id=%s message=%s body=%v", body.Channel, body.Text, respBody)
	}

	return nil
}

// Post to API
func postJSON(url string, body interface{}) (*http.Response, error) {
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.MustGet("SLACK_API_BOT_TOKEN")))

	client := http.Client{}
	return client.Do(req)
}
