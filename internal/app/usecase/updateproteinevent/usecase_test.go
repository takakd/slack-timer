package updateproteinevent

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

func TestInteractor_saveProteinEventValue(t *testing.T) {

	t.Run("ok:create", func(t *testing.T) {
		ctx := context.TODO()
		userId := "abc"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := apprule.NewMockRepository(ctrl)
		m.EXPECT().FindProteinEvent(gomock.Eq(ctx), gomock.Eq(userId)).
			Return(nil, nil)

		now := time.Now()
		event, _ := enterpriserule.NewProteinEvent(userId)
		event.UtcTimeToDrink = now
		m.EXPECT().SaveProteinEvent(gomock.Eq(ctx), gomock.Len(1)).Return(nil, nil)

		interactor := &Interactor{
			repository: m,
		}

		data := interactor.saveProteinEventValue(ctx, userId, 0)
		assert.NoError(t, data.Result)
	})

	t.Run("ok:next notify", func(t *testing.T) {
		ctx := context.TODO()
		event, _ := enterpriserule.NewProteinEvent("abc")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := apprule.NewMockRepository(ctrl)
		m.EXPECT().FindProteinEvent(gomock.Eq(ctx), gomock.Eq(event.UserId)).
			Return(event, nil)
		m.EXPECT().SaveProteinEvent(gomock.Eq(ctx), gomock.Eq([]*enterpriserule.ProteinEvent{event})).
			Return(nil, nil)

		interactor := &Interactor{
			repository: m,
		}

		data := interactor.saveProteinEventValue(context.TODO(), event.UserId, 0)
		assert.NoError(t, data.Result)
	})

	t.Run("ok:interval", func(t *testing.T) {
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

		interactor := &Interactor{
			repository: m,
		}

		data := interactor.saveProteinEventValue(context.TODO(), userId, interval)
		assert.NoError(t, data.Result)
	})

	t.Run("ng:create", func(t *testing.T) {
		ctx := context.TODO()
		userId := ""

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := apprule.NewMockRepository(ctrl)
		m.EXPECT().FindProteinEvent(gomock.Eq(ctx), gomock.Eq(userId)).
			Return(nil, nil)

		interactor := &Interactor{
			repository: m,
		}

		data := interactor.saveProteinEventValue(context.TODO(), userId, 0)
		assert.True(t, errors.Is(data.Result, ErrCreate))
	})

	t.Run("ng:update", func(t *testing.T) {
		ctx := context.TODO()
		userId := ""

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := apprule.NewMockRepository(ctrl)
		m.EXPECT().FindProteinEvent(gomock.Eq(ctx), gomock.Eq(userId)).
			Return(nil, errors.New("error"))

		interactor := &Interactor{
			repository: m,
		}

		data := interactor.saveProteinEventValue(context.TODO(), userId, 0)
		assert.True(t, errors.Is(data.Result, ErrFind))
	})
}
