// Package slackcontroller provides the slack Event API callback handler.
// Ref: https://api.slack.com/events-api#the-events-api__receiving-events
package settime

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventCallbackData_isVerificationEvent(t *testing.T) {
	cases := []struct {
		name         string
		dataType     string
		verification bool
	}{
		{"ok", "url_verification", true},
		{"ng", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := &EventCallbackData{
				Type: c.dataType,
			}
			assert.Equal(t, d.isVerificationEvent(), c.verification)
		})
	}
}

func TestMessageEvent_isSetTimeEvent(t *testing.T) {
	cases := []struct {
		name        string
		eventType   string
		channelType string
		text        string
		isSet       bool
	}{
		{"ok:with text", "message", "im", "set 10 Hi!", true},
		{"ng:without text", "message", "im", "set 10", false},
		{"ng:wrong event type", "wrong", "im", "set 10", false},
		{"ng:wrong channel type", "wrong", "channel", "set 10", false},
		{"ng:wrong body", "message", "im", "set a", false},
		{"ng:empty body", "message", "im", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := &MessageEvent{
				Type:        c.eventType,
				ChannelType: c.channelType,
				Text:        c.text,
			}
			got := d.isSetTimeEvent()
			assert.Equal(t, c.isSet, got)
		})
	}
}

func TestMessageEvent_isOnEvent(t *testing.T) {
	cases := []struct {
		name        string
		eventType   string
		channelType string
		text        string
		isSet       bool
	}{
		{"ok", "message", "im", "on", true},
		{"ng:wrong event type", "wrong", "im", "on", false},
		{"ng:wrong channel type", "message", "channel", "on", false},
		{"ng:wrong body", "message", "im", "onwrong", false},
		{"ng:empty body", "message", "im", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := &MessageEvent{
				Type:        c.eventType,
				ChannelType: c.channelType,
				Text:        c.text,
			}
			got := d.isOnEvent()
			assert.Equal(t, c.isSet, got)
		})
	}
}

func TestMessageEvent_isOffEvent(t *testing.T) {
	cases := []struct {
		name        string
		eventType   string
		channelType string
		text        string
		isSet       bool
	}{
		{"ok", "message", "im", "off", true},
		{"ng:wrong event type", "wrong", "im", "off", false},
		{"ng:wrong channel type", "message", "channel", "off", false},
		{"ng:wrong body", "message", "im", "offwrong", false},
		{"ng:empty body", "message", "im", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := &MessageEvent{
				Type:        c.eventType,
				ChannelType: c.channelType,
				Text:        c.text,
			}
			got := d.isOffEvent()
			assert.Equal(t, c.isSet, got)
		})
	}
}

func TestMessageEvent_eventUnixTimeStamp(t *testing.T) {
	cases := []struct {
		name    string
		ts      string
		tsNano  string
		success bool
	}{
		{"ok", "1607165903", "000010", true},
		{"ok:empty nano", "1607165903", "", true},
		{"ng:invalid format", "abc", "", false},
		{"ng:empty", "", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := &MessageEvent{
				EventTs: fmt.Sprintf("%s.%s", c.ts, c.tsNano),
			}
			got, err := d.eventUnixTimeStamp()

			if c.success {
				assert.NoError(t, err)
				want, err := strconv.ParseInt(c.ts, 10, 64)
				require.NoError(t, err)
				assert.Equal(t, want, got)

			} else {
				assert.Error(t, err)
			}
		})
	}
}
