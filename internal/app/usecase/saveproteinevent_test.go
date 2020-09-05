package usecase

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"proteinreminder/internal/app/apprule"
	"proteinreminder/internal/app/enterpriserule"
	"proteinreminder/internal/pkg/testutil"
	"testing"
	"time"
)

func TestSaveProteinEvent_saveProteinEventValue(t *testing.T) {
	t.Run("OK: save new, to drink", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		userId := "abc"
		m := apprule.NewMockRepository(ctrl)
		m.EXPECT().FindProteinEvent(gomock.Eq(ctx), gomock.Eq(userId)).
			Return(nil, nil)

		now := time.Now()
		event, _ := enterpriserule.NewProteinEvent(userId)
		event.UtcTimeToDrink = now
		m.EXPECT().SaveProteinEvent(gomock.Eq(ctx), gomock.Eq([]*enterpriserule.ProteinEvent{event})).
			Return(nil, nil)

		s, _ := NewSaveProteinEvent(m)
		if err := s.saveProteinEventValue(context.TODO(), userId, &now, nil); err != SaveProteinEventNoError {
			t.Errorf("failed to save")
		}
	})

	t.Run("OK: update, to drink", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		event, _ := enterpriserule.NewProteinEvent("abc")
		m := apprule.NewMockRepository(ctrl)
		m.EXPECT().FindProteinEvent(gomock.Eq(ctx), gomock.Eq(event.UserId)).
			Return(event, nil)
		m.EXPECT().SaveProteinEvent(gomock.Eq(ctx), gomock.Eq([]*enterpriserule.ProteinEvent{event})).
			Return(nil, nil)

		s, _ := NewSaveProteinEvent(m)

		now := time.Now()
		if err := s.saveProteinEventValue(context.TODO(), event.UserId, &now, nil); err != SaveProteinEventNoError {
			t.Errorf("failed to save")
		}
	})

	t.Run("OK: interval", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		userId := "abc"
		m := apprule.NewMockRepository(ctrl)
		m.EXPECT().FindProteinEvent(gomock.Eq(ctx), gomock.Eq(userId)).
			Return(nil, nil)

		interval := time.Duration(1)
		event, _ := enterpriserule.NewProteinEvent(userId)
		event.DrinkTimeIntervalSec = interval
		m.EXPECT().SaveProteinEvent(gomock.Eq(ctx), gomock.Eq([]*enterpriserule.ProteinEvent{event})).
			Return(nil, nil)

		s, _ := NewSaveProteinEvent(m)
		if err := s.saveProteinEventValue(context.TODO(), userId, nil, &interval); err != SaveProteinEventNoError {
			t.Errorf("failed to save")
		}
	})

	t.Run("NG: new event", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		userId := ""
		m := apprule.NewMockRepository(ctrl)
		m.EXPECT().FindProteinEvent(gomock.Eq(ctx), gomock.Eq(userId)).
			Return(nil, nil)

		s, _ := NewSaveProteinEvent(m)
		dummy := time.Now()
		if err := s.saveProteinEventValue(context.TODO(), userId, &dummy, nil); err != SaveProteinEventErrorCreate {
			t.Error(testutil.MakeTestMessageWithGotWant(err, SaveProteinEventErrorCreate))
		}
	})

	t.Run("NG: find event", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		userId := ""
		m := apprule.NewMockRepository(ctrl)
		m.EXPECT().FindProteinEvent(gomock.Eq(ctx), gomock.Eq(userId)).
			Return(nil, errors.New("error"))

		s, _ := NewSaveProteinEvent(m)
		dummy := time.Now()
		if err := s.saveProteinEventValue(context.TODO(), userId, &dummy, nil); err != SaveProteinEventErrorFind {
			t.Error(testutil.MakeTestMessageWithGotWant(err, SaveProteinEventErrorFind))
		}
	})
}
