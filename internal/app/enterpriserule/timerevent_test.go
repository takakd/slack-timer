package enterpriserule

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewTimerEvent(t *testing.T) {
	cases := []struct {
		name   string
		userId string
		want   *TimerEvent
	}{
		{name: "ok", userId: "id1234", want: &TimerEvent{UserId: "id1234", State: timerEventStateWait}},
		{name: "ng", userId: "", want: nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := NewTimerEvent(c.userId)
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
		UserId:           "id1",
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
		{name: "ng:nil", lhs: event, rhs: nil, result: false},
		{name: "ng:user_id", lhs: event, rhs: &TimerEvent{
			UserId: "id2", NotificationTime: now, IntervalMin: sec}, result: false},
		{name: "ng:utc_time", lhs: event, rhs: &TimerEvent{
			UserId: "id1", NotificationTime: now.Add(time.Second * 1), IntervalMin: sec}, result: false},
		{name: "ng:time_interval_sec", lhs: event, rhs: &TimerEvent{
			UserId: "id1", NotificationTime: now, IntervalMin: sec + 1}, result: false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// TODO
			assert.Equal(t, c.lhs.Equal(c.rhs), c.result)
			//if c.lhs.Equal(c.rhs) != c.result {
			//	t.Error(testutil.MakeTestMessageWithGotWant(c.lhs, c.rhs))
			//}
		})
	}
}

// TODO: add unit test
