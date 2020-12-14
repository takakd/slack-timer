// Package settime provides features that set user's event notification time.
package settime

import (
	"fmt"
	"regexp"
	"slacktimer/internal/app/util/appcontext"
	"strconv"
	"strings"
)

// Command types entered by users.
const (
	_cmdSet = "set"
	_cmdOn  = "on"
	_cmdOff = "off"
)

// ControllerHandler is called by Lambda handler.
type ControllerHandler interface {
	Handle(ac appcontext.AppContext, input HandleInput) *Response
}

// Response is returns of Controller.
type Response struct {
	Error      error
	StatusCode int
	Body       interface{}
}

// HandleInput is input parameter of Controller.
type HandleInput struct {
	EventData EventCallbackData
}

// EventCallbackData represents Slack Event API payload.
// This contained in Slack EventAPI Request.
// Ref: https://api.slack.com/events-api#the-events-api__receiving-events
type EventCallbackData struct {
	Token     string `json:"token"`
	TeamID    string `json:"team_id"`
	Type      string `json:"type"`
	EventTime int    `json:"event_time"`

	// This field is only included in Message Event.
	// Ref: https://api.slack.com/events
	MessageEvent MessageEvent `json:"event"`

	// This field is only included in URL Verification Event.
	// Ref: https://api.slack.com/events/url_verification
	Challenge string `json:"challenge"`
}

func (e EventCallbackData) isVerificationEvent() bool {
	return e.Type == "url_verification"
}

// MessageEvent represents the message data in EventCallbackData.
// Ref: https://api.slack.com/events/message.im
type MessageEvent struct {
	Type        string `json:"type"`
	EventTs     string `json:"event_ts"`
	User        string `json:"user"`
	Ts          string `json:"ts"`
	Text        string `json:"text"`
	ChannelType string `json:"channel_type"`
	// This field is only included in message from the bot.
	BotID string `json:"bot_id,omitempty"`
}

// Extract second, because the format of timestamp sent by Slack has nano second.
func (m MessageEvent) eventUnixTimeStamp() (ts int64, err error) {
	s := strings.Split(m.EventTs, ".")
	if len(s) < 1 {
		err = fmt.Errorf("invalid format %s", m.EventTs)
		return
	}

	ts, err = strconv.ParseInt(s[0], 10, 64)
	return
}

func (m MessageEvent) isBotMessage() bool {
	return m.BotID != ""
}

func (m MessageEvent) isMatchCommand(pattern, cmd string) bool {
	// Only use in DM channel.
	if m.Type != "message" || m.ChannelType != "im" {
		return false
	}
	re := regexp.MustCompile(fmt.Sprintf(pattern, cmd))
	matches := re.FindStringSubmatch(m.Text)
	if matches == nil || len(matches) < 2 || matches[1] != cmd {
		return false
	}
	return true
}

func (m MessageEvent) isSetTimeEvent() bool {
	return m.isMatchCommand(`^(%s)\s+(\d+)\s+([\s\S]*)`, _cmdSet)
}

func (m MessageEvent) isOnEvent() bool {
	return m.isMatchCommand(`^(%s)$`, _cmdOn)
}

func (m MessageEvent) isOffEvent() bool {
	return m.isMatchCommand(`^(%s)$`, _cmdOff)
}
