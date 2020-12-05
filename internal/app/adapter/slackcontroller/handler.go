// Package slackcontroller provides the slack Event API callback handler.
// Ref: https://api.slack.com/events-api#the-events-api__receiving-events
package slackcontroller

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"strconv"
	"strings"
)

// Command types entered by users.
const (
	CmdSet = "set"
)

type Handler interface {
	Handler(ctx context.Context, input HandlerInput) *Response
}

type Response struct {
	Error      error
	StatusCode int
	Body       interface{}
}

type HandlerInput struct {
	EventData EventCallbackData
}

// Slack EventAPI Notification data
type EventCallbackData struct {
	Token  string `json:"token"`
	TeamId string `json:"team_id"`
	// Ref. https://api.slack.com/events
	MessageEvent MessageEvent `json:"event"`
	Type         string       `json:"type"`
	EventTime    int          `json:"event_time"`

	// This field is only included in URL Verification Event.
	// Ref: https://api.slack.com/events/url_verification
	Challenge string `json:"challenge"`
}

func (e *EventCallbackData) isVerificationEvent() bool {
	return e.Type == "url_verification"
}

type MessageEvent struct {
	Type    string `json:"type"`
	EventTs string `json:"event_ts"`
	User    string `json:"user"`
	Ts      string `json:"ts"`
	Text    string `json:"text"`
}

func (m MessageEvent) eventUnixTimeStamp() (ts int64, err error) {
	s := strings.Split(m.EventTs, ".")
	if len(s) < 1 {
		err = fmt.Errorf("invalid format %s", m.EventTs)
		return
	}

	ts, err = strconv.ParseInt(s[0], 10, 64)
	return
}

func (e *MessageEvent) isSetEvent() bool {
	if e.Type != "message" {
		return false
	}

	// e.g. set 10
	re := regexp.MustCompile(`^([^\s]*)\s*`)
	m := re.FindStringSubmatch(e.Text)
	if m == nil || m[1] != CmdSet {
		return false
	}

	return true
}

type SlackEventHandler struct {
}

func NewHandler() Handler {
	h := SlackEventHandler{}
	return h
}

func (h SlackEventHandler) Handler(ctx context.Context, input HandlerInput) *Response {
	// Create request struct corresponding to input.

	// URL verification event
	if input.EventData.isVerificationEvent() {
		log.Info("url verification event")
		rh := di.Get("slackcontroller.urlverificationhandler").(UrlVerificationRequestHandler)
		return rh.Handler(ctx, input.EventData)
	}

	// Set interval minutes event
	if !input.EventData.MessageEvent.isSetEvent() {
		return makeErrorHandlerResponse("invalid event", fmt.Sprintf("type=%s", input.EventData.MessageEvent.Type))
	}

	rh := di.Get("slackcontroller.setcontroller").(SetRequestHandler)
	return rh.Handler(ctx, input.EventData)
}

type ResponseErrorBody struct {
	Message string
	Detail  string
}

func makeErrorHandlerResponse(message string, detail string) *Response {
	body := &ResponseErrorBody{
		Message: message,
	}
	if detail != "" {
		body.Detail = detail
	}
	return &Response{
		StatusCode: http.StatusInternalServerError,
		Body:       body,
	}
}
