package usecase

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"proteinreminder/internal/app/apprule"
	"proteinreminder/internal/app/enterpriserule"
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
		m.EXPECT().SaveProteinEvent(gomock.Eq(ctx), gomock.Len(1)).Return(nil, nil)

		s, _ := NewSaveProteinEvent(m)
		err := s.saveProteinEventValue(context.TODO(), userId, 0)
		assert.NoError(t, err)
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

		err := s.saveProteinEventValue(context.TODO(), event.UserId, 0)
		assert.NoError(t, err)
	})

	t.Run("OK: interval", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		userId := "abc"
		m := apprule.NewMockRepository(ctrl)
		m.EXPECT().FindProteinEvent(gomock.Eq(ctx), gomock.Eq(userId)).
			Return(nil, nil)

		interval := 1
		event, _ := enterpriserule.NewProteinEvent(userId)
		event.DrinkTimeIntervalMin = interval
		m.EXPECT().SaveProteinEvent(gomock.Eq(ctx), gomock.Len(1)).
			Return(nil, nil)

		s, _ := NewSaveProteinEvent(m)
		err := s.saveProteinEventValue(context.TODO(), userId, interval)
		assert.NoError(t, err)
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
		err := s.saveProteinEventValue(context.TODO(), userId, 0)
		assert.True(t, errors.Is(err, ErrCreate))
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
		err := s.saveProteinEventValue(context.TODO(), userId, 0)
		assert.True(t, errors.Is(err, ErrFind))
	})
}
