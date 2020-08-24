package interfaces

import (
	"context"
	"proteinreminder/internal/entity"
	"proteinreminder/internal/testutil"
	"testing"
	"time"
)

func testEvents() []*entity.ProteinEvent {
	return []*entity.ProteinEvent{
		{
			"id1", time.Now().UTC(), 0,
		},
		{
			"id2", time.Now().UTC(), 0,
		},
	}
}

func TestFindProteinEvent(t *testing.T) {
	testEvents := testEvents()

	repo := NewMongoDbRepository()
	_, err := repo.SaveProteinEvent(context.TODO(), testEvents)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("OK", func(t *testing.T) {
		for _, event := range testEvents {
			repo := NewMongoDbRepository()
			got, err := repo.FindProteinEvent(context.TODO(), event.UserId)
			if err != nil {
				t.Error(err)
			}

			if !got.Equal(event) {
				t.Error(testutil.MakeTestMessageWithGotWant(got, event))
			}
		}
	})

	t.Run("NG", func(t *testing.T) {
		repo := NewMongoDbRepository()
		event, err := repo.FindProteinEvent(context.TODO(), "disable id")
		if event != nil || err != nil {
			t.Error("should not exist")
		}
	})
}

func TestFindProteinEventByTime(t *testing.T) {
	now := time.Now().UTC()
	events := []*entity.ProteinEvent{
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

	repo := NewMongoDbRepository()

	repo.SaveProteinEvent(context.TODO(), events)

	got, err := repo.FindProteinEventByTime(context.TODO(), from, to)
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
	repo := NewMongoDbRepository()
	savedEvents, err := repo.SaveProteinEvent(context.TODO(), testEvents)
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
