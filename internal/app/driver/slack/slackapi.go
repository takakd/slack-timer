package slack

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"slacktimer/internal/pkg/config"
	"slacktimer/internal/pkg/httputil"
)

type SlackApi interface {
	ConversationsOpen(userId string) (string, error)
	ChatPostMessage(channelId string, message string) error
}

// Request Body
// Ref: https://api.slack.com/methods/chat.postMessage
type ChatPostMessageBody struct {
	Token   string `json:"token"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type SlackApiDriver struct {
}

func NewSlackApiDriver() SlackApi {
	s := &SlackApiDriver{}
	return s
}

type ConversationsOpenResponse struct {
	Ok      bool     `json:"ok"`
	Channel []string `json:"id,omitempty"`
	// Be set if the response is error
	Error string `json:"error,omitempty"`
}

// Ref: https://api.slack.com/methods/conversations.open
func (s *SlackApiDriver) ConversationsOpen(userId string) (string, error) {
	body := struct {
		Token string `json:"token"`
		Users string `json:"users"`
	}{
		config.MustGet("SLACK_API_BOT_TOKEN"),
		userId,
	}
	url := config.MustGet("SLACK_API_URL_CONVERSATIONSOPEN")
	resp, err := postJson(url, body)
	if err != nil {
		return "", err
	}

	ok := resp.StatusCode == http.StatusOK
	if !ok {
		return "", errors.New("request error Slack API:conversations.open")
	}

	respBuf, err := httputil.GetResponseBody(resp)
	if err != nil {
		return "", errors.New("failed to read response of Slack API:conversations.open")
	}

	var respBody ConversationsOpenResponse
	err = json.Unmarshal(respBuf, &respBody)
	if err != nil {
		return "", errors.New("unmarshal error Slack API:conversations.open")
	}

	// It must be returned one ID because of sending one user ID.
	if !respBody.Ok || len(respBody.Channel) != 1 {
		return "", errors.New("API response NG Slack API:conversations.open")
	}

	return respBody.Channel[0], nil
}

// Ref: https://api.slack.com/methods/chat.postMessage
func (s *SlackApiDriver) ChatPostMessage(channelId string, message string) error {
	body := &ChatPostMessageBody{
		Token:   config.Get("SLACK_API_BOT_TOKEN", ""),
		Channel: channelId,
		Text:    message,
	}
	url := config.Get("SLACK_API_URL_CHAPOSTMESSAGE", "")
	resp, err := postJson(url, body)
	if err != nil {
		return err
	}

	ok := resp.StatusCode == http.StatusOK
	if !ok {
		return errors.New("failed to request to slack api")
	}

	return nil
}

// Post to SlackAPI
func postJson(url string, body interface{}) (*http.Response, error) {
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	return client.Do(req)
}
