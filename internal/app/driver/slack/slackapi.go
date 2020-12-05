package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"slacktimer/internal/app/util/config"
	"slacktimer/internal/app/util/log"
	"slacktimer/internal/pkg/helper"
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
	// Token is set to Bearer Header.
	Users string `json:"users"`
}

type ConversationsOpenResponseBody struct {
	Ok      bool                                 `json:"ok"`
	Channel ConversationsOpenResponseBodyChannel `json:"channel,omitempty"`
	// Be set if the response is error
	Error string `json:"error,omitempty"`
}

type ConversationsOpenResponseBodyChannel struct {
	Id string `json:"id"`
}

// Ref: https://api.slack.com/methods/conversations.open
func (s *SlackApiDriver) ConversationsOpen(userId string) (string, error) {
	body := &ConversationsOpenRequestBody{
		userId,
	}
	url := config.MustGet("SLACK_API_URL_CONVERSATIONSOPEN")
	resp, err := postJson(url, body)
	if err != nil {
		return "", err
	}

	ok := resp.StatusCode == http.StatusOK
	if !ok {
		return "", fmt.Errorf("request error slack conversations.open user_id=%s: %w", userId, err)
	}

	respBuf, err := helper.GetResponseBody(resp)
	if err != nil {
		return "", fmt.Errorf("response reading error slack conversations.open user_id=%s: %w", userId, err)
	}

	log.Debug("response body slack conversations.open", respBuf)

	var respBody ConversationsOpenResponseBody
	err = json.Unmarshal(respBuf, &respBody)
	if err != nil {
		return "", fmt.Errorf("unmarshal error slack conversations.open user_id=%s: %w", userId, err)
	}

	// It must be returned one ID because of sending one user ID.
	if !respBody.Ok {
		return "", fmt.Errorf("response NG slack conversations.open user_id=%s body=%v", userId, respBody)
	}

	return respBody.Channel.Id, nil
}

// Ref: https://api.slack.com/methods/chat.postMessage
type ChatPostMessageRequestBody struct {
	// Token is set to Bearer Header.
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// Ref: https://api.slack.com/methods/chat.postMessage
type ChatPostMessageResponseBody struct {
	// Define only need
	Ok bool `json:"ok"`
	// Be set if the response is error
	Error string `json:"error,omitempty"`
}

// Ref: https://api.slack.com/methods/chat.postMessage
func (s *SlackApiDriver) ChatPostMessage(channelId string, message string) error {
	body := &ChatPostMessageRequestBody{
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
		return fmt.Errorf("request error slack chat.postMessage channel_id=%s message=%s: %w", channelId, message, err)
	}

	respBuf, err := helper.GetResponseBody(resp)
	if err != nil {
		return fmt.Errorf("response reading error slack chat.postMessage channel_id=%s message=%s: %w", channelId, message, err)
	}

	log.Debug("response body slack chat.postMessage", respBuf)

	var respBody ChatPostMessageResponseBody
	err = json.Unmarshal(respBuf, &respBody)
	if err != nil {
		return fmt.Errorf("unmarshal error slack chat.postMessage channel_id=%s message=%s", channelId, message)
	}

	if !respBody.Ok {
		return fmt.Errorf("response NG slack chat.postMessage channel_id=%s message=%s body=%v", channelId, message, respBody)
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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.MustGet("SLACK_API_BOT_TOKEN")))

	client := http.Client{}
	return client.Do(req)
}
