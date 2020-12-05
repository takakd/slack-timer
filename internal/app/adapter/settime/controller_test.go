// Package slackcontroller provides the slack Event API callback handler.
// Ref: https://api.slack.com/events-api#the-events-api__receiving-events
package settime

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
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

func TestMessageEvent_isSetEvent(t *testing.T) {
	cases := []struct {
		name      string
		eventType string
		text      string
		isSet     bool
	}{
		{"ok", "message", "set 10", true},
		{"ng:wrong type", "wrong", "set 10", false},
		{"ng:wrong body", "messazge", "set a", false},
		{"ng:empty body", "messazge", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := &MessageEvent{
				Type: c.eventType,
				Text: c.text,
			}
			got := d.isSetEvent()
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
