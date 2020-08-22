package entity

import (
	"context"
	"proteinreminder/internal/testutil"
	"testing"
	"time"
)

// --------------------------------------------------------
// Repository Role

func testEvents() []*ProteinEvent {
	return []*ProteinEvent{
		{
			"id1", time.Now().UTC(), 0,
		},
		{
			"id2", time.Now().UTC(), 0,
		},
	}
}

func TestGetProteinEvent(t *testing.T) {
	testEvents := testEvents()
	err, _ := SaveProteinEvent(context.TODO(), testEvents)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("ok", func(t *testing.T) {
		for _, event := range testEvents {
			got, err := GetProteinEvent(context.TODO(), event.UserId)
			if err != nil {
				t.Error(err)
			}

			if !got.Equal(event) {
				t.Error(testutil.MakeTestMessageWithGotWant(got, event))
			}
		}
	})

	t.Run("ng", func(t *testing.T) {
		_, err := GetProteinEvent(context.TODO(), "disable id")
		if err == nil {
			t.Error("should not exist")
		}
	})

}

func TestFindProteinEventByTime(t *testing.T) {
	now := time.Now().UTC()
	events := []*ProteinEvent{
		{
			"id1", now, 2,
		},
		{
			"id2", now, 2,
		},
		{
			"id3", now, 2,
		},
	}

	from := time.Now().UTC()
	to := from.AddDate(0, 0, 1)

	events[0].UtcTimeToDrink = from
	events[1].UtcTimeToDrink = to
	events[2].UtcTimeToDrink = to.AddDate(0, 0, 1)

	SaveProteinEvent(context.TODO(), events)

	got, err := FindProteinEventByTime(context.TODO(), from, to)
	if err != nil {
		t.Error(err)
	}
	wants := events[:2]
	if len(got) != 2 || !wants[0].Equal(got[0]) || !wants[1].Equal(got[1]) {
		t.Error(testutil.MakeTestMessageWithGotWant(got[0], wants[0]))
		t.Error(testutil.MakeTestMessageWithGotWant(got[1], wants[1]))
	}
}

func TestSaveProteinEvent(t *testing.T) {
	testEvents := testEvents()
	err, savedEvents := SaveProteinEvent(context.TODO(), testEvents)
	if err != nil {
		t.Error(err)
		return
	}

	for i, savedEvent := range savedEvents {
		if !testEvents[i].Equal(savedEvent) {
			t.Error(testutil.MakeTestMessageWithGotWant(savedEvent, testEvents[i]))
		}
	}
}

// --------------------------------------------------------
// Entity

func TestNewProteinEvent(t *testing.T) {
	cases := []struct {
		name string
		in   string
		out  *ProteinEvent
	}{
		{name: "ok", in: "id1234", out: &ProteinEvent{UserId: "id1234"}},
		{name: "ng", in: "", out: nil},
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
		{name: "ok", lhs: event, rhs: event, result: true},
		{name: "ng:nil", lhs: event, rhs: nil, result: false},
		{name: "ng:user_id", lhs: event, rhs: &ProteinEvent{"id2", now, sec}, result: false},
		{name: "ng:utc_time_to_drink", lhs: event, rhs: &ProteinEvent{"id1", now.Add(time.Second * 1), sec}, result: false},
		{name: "ng:drink_time_interval_sec", lhs: event, rhs: &ProteinEvent{"id1", now, sec + 1}, result: false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.lhs.Equal(c.rhs) != c.result {
				t.Error(testutil.MakeTestMessageWithGotWant(c.lhs, c.rhs))
			}
		})
	}
}
