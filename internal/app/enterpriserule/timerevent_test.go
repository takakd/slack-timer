package enterpriserule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTimerEvent(t *testing.T) {
	cases := []struct {
		name   string
		userID string
		want   *TimerEvent
	}{
		{name: "ok", userID: "id1234", want: &TimerEvent{UserID: "id1234", State: _timerEventStateWait}},
		{name: "ng", userID: "", want: nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := NewTimerEvent(c.userID)
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
		UserID:           "id1",
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
			UserID: "id2", NotificationTime: now, IntervalMin: sec}, result: false},
		{name: "ng:utc_time", lhs: event, rhs: &TimerEvent{
			UserID: "id1", NotificationTime: now.Add(time.Second * 1), IntervalMin: sec}, result: false},
		{name: "ng:time_interval_sec", lhs: event, rhs: &TimerEvent{
			UserID: "id1", NotificationTime: now, IntervalMin: sec + 1}, result: false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// TODO
			assert.Equal(t, c.lhs.Equal(*c.rhs), c.result)
			//if c.lhs.Equal(c.rhs) != c.result {
			//	t.Error(helper.MakeTestMessageWithGotWant(c.lhs, c.rhs))
			//}
		})
	}
}

// TODO: add unit test
