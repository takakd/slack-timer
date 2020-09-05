package enterpriserule

import (
	"proteinreminder/internal/pkg/testutil"
	"testing"
	"time"
)

func TestNewProteinEvent(t *testing.T) {
	cases := []struct {
		name string
		in   string
		out  *ProteinEvent
	}{
		{name: "OK", in: "id1234", out: &ProteinEvent{UserId: "id1234"}},
		{name: "NG", in: "", out: nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, _ := NewProteinEvent(c.in)
			if c.out == nil {
				if got != nil {
					t.Error(testutil.MakeTestMessageWithGotWant(got, c.out))
				}
			} else {
				if !got.Equal(c.out) {
					t.Error(testutil.MakeTestMessageWithGotWant(got, c.out))
				}
			}
		})
	}
}

func TestProteinEvent_Equal(t *testing.T) {
	now := time.Now().UTC()
	sec := time.Duration(10)
	event := &ProteinEvent{
		"id1",
		now,
		sec,
	}
	cases := []struct {
		name   string
		lhs    *ProteinEvent
		rhs    *ProteinEvent
		result bool
	}{
		{name: "OK", lhs: event, rhs: event, result: true},
		{name: "NG:nil", lhs: event, rhs: nil, result: false},
		{name: "NG:user_id", lhs: event, rhs: &ProteinEvent{"id2", now, sec}, result: false},
		{name: "NG:utc_time_to_drink", lhs: event, rhs: &ProteinEvent{"id1", now.Add(time.Second * 1), sec}, result: false},
		{name: "NG:drink_time_interval_sec", lhs: event, rhs: &ProteinEvent{"id1", now, sec + 1}, result: false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.lhs.Equal(c.rhs) != c.result {
				t.Error(testutil.MakeTestMessageWithGotWant(c.lhs, c.rhs))
			}
		})
	}
}
