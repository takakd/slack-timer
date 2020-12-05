// Package slackcontroller provides the slack Event API callback handler.
// Ref: https://api.slack.com/events-api#the-events-api__receiving-events
package settime

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Controller interface {
	Handle(ctx context.Context, input HandleInput) *Response
}

type Response struct {
	Error      error
	StatusCode int
	Body       interface{}
}

type HandleInput struct {
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

func (m MessageEvent) isSetEvent() bool {
	if m.Type != "message" {
		return false
	}

	// e.g. set 10
	re := regexp.MustCompile(`^([^\s]*)\s*`)
	matches := re.FindStringSubmatch(m.Text)
	if matches == nil || matches[1] != CmdSet {
		return false
	}

	return true
}
