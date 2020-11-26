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

type SlackApiDriver struct {
}

func NewSlackApiDriver() SlackApi {
	s := &SlackApiDriver{}
	return s
}

// Ref: https://api.slack.com/methods/conversations.open
type ConversationsOpenRequestBody struct {
	Token string `json:"token"`
	Users string `json:"users"`
}

type ConversationsOpenResponseBody struct {
	Ok      bool     `json:"ok"`
	Channel []string `json:"id,omitempty"`
	// Be set if the response is error
	Error string `json:"error,omitempty"`
}

// Ref: https://api.slack.com/methods/conversations.open
func (s *SlackApiDriver) ConversationsOpen(userId string) (string, error) {
	body := &ConversationsOpenRequestBody{
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
		return "", errors.New("response reading error Slack API:conversations.open")
	}

	var respBody ConversationsOpenResponseBody
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
type ChatPostMessageRequestBody struct {
	Token   string `json:"token"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// Ref: https://api.slack.com/methods/chat.postMessage
type ChatPostMessageResponseBody struct {
	// Define only need
	Ok bool `json:"ok"`
}

// Ref: https://api.slack.com/methods/chat.postMessage
func (s *SlackApiDriver) ChatPostMessage(channelId string, message string) error {
	body := &ChatPostMessageRequestBody{
		Token:   config.MustGet("SLACK_API_BOT_TOKEN"),
		Channel: channelId,
		Text:    message,
	}
	url := config.MustGet("SLACK_API_URL_CHATPOSTMESSAGE")
	resp, err := postJson(url, body)
	if err != nil {
		return err
	}

	ok := resp.StatusCode == http.StatusOK
	if !ok {
		return errors.New("request error Slack API:chat.postMessage")
	}

	respBuf, err := httputil.GetResponseBody(resp)
	if err != nil {
		return errors.New("response reading error Slack API:chat.postMessage")
	}

	var respBody ChatPostMessageResponseBody
	err = json.Unmarshal(respBuf, &respBody)
	if err != nil {
		return errors.New("unmarshal error Slack API:chat.postMessage")
	}

	if !respBody.Ok {
		return errors.New("API response NG Slack API:chat.postMessage")
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
