package enterpriserule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTimerEvent(t *testing.T) {
	cases := []struct {
		name   string
		userID string
		text   string
		want   *TimerEvent
	}{
		{name: "ok", userID: "id1234", text: "Hi!", want: &TimerEvent{userID: "id1234", State: TimerEventStateWait, text: "Hi!"}},
		{name: "ng", userID: "", text: "Hi!", want: nil},
		{name: "ng", userID: "id1234", text: "", want: nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := NewTimerEvent(c.userID, c.text)
			assert.Equal(t, c.want, got)
			if err != nil {
				assert.Nil(t, c.want)
			} else {
				assert.NotNil(t, c.want)
			}
		})
	}
}

func TestTimerEvent_Equal(t *testing.T) {
	now := time.Now().UTC()
	sec := 10
	event := &TimerEvent{
		userID:           "id1",
		text:             "Hi!",
		NotificationTime: now,
		IntervalMin:      sec,
	}
	cases := []struct {
		name   string
		lhs    *TimerEvent
		rhs    *TimerEvent
		result bool
	}{
		{name: "ok", lhs: event, rhs: event, result: true},
		{name: "ng:user_id", lhs: event, rhs: &TimerEvent{
			userID: "id2", text: "Hi!", NotificationTime: now, IntervalMin: sec}, result: false},
		{name: "ng:utc_time", lhs: event, rhs: &TimerEvent{
			userID: "id1", text: "Hi!", NotificationTime: now.Add(time.Second * 1), IntervalMin: sec}, result: false},
		{name: "ng:time_interval_sec", lhs: event, rhs: &TimerEvent{
			userID: "id1", text: "Hi!", NotificationTime: now, IntervalMin: sec + 1}, result: false},
		{name: "ng:text", lhs: event, rhs: &TimerEvent{
			userID: "id1", text: "Hi!!", NotificationTime: now, IntervalMin: sec}, result: false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.lhs.Equal(*c.rhs), c.result)
		})
	}
}

func TestTimerEvent_UserID(t *testing.T) {
	got, err := NewTimerEvent("test", "Hi!")
	require.NoError(t, err)
	assert.Equal(t, "test", got.UserID())
}

func TestTimerEvent_Text(t *testing.T) {
	got, err := NewTimerEvent("test", "Hi!")
	require.NoError(t, err)
	assert.Equal(t, "Hi!", got.Text())
}
