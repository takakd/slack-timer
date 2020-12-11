package timeroffevent

import (
	"fmt"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/di"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewInteractor(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		assert.NotPanics(t, func() {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mr := NewMockRepository(ctrl)
			md := di.NewMockDI(ctrl)
			md.EXPECT().Get("timeroffevent.Repository").Return(mr)
			di.SetDi(md)

			NewInteractor()
		})
	})

	t.Run("ng", func(t *testing.T) {
		assert.Panics(t, func() {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			md := di.NewMockDI(ctrl)
			md.EXPECT().Get("timeroffevent.Repository").Return(nil)
			di.SetDi(md)

			NewInteractor()
		})
	})
}

func TestInteractor_SetEventOff(t *testing.T) {
	t.Run("ok:set", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()

		caseEvent, _ := enterpriserule.NewTimerEvent("abc", "Hi!")
		caseEvent.State = enterpriserule.TimerEventStateWait

		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEvent(caseEvent.UserID()).
			Return(caseEvent, nil)
		m.EXPECT().SaveTimerEvent(caseEvent).Return(nil, nil)

		want, _ := enterpriserule.NewTimerEvent(caseEvent.UserID(), caseEvent.Text())
		want.State = enterpriserule.TimerEventStateDisabled

		wantOutputData := OutputData{
			Result:     nil,
			SavedEvent: want,
		}

		mp := NewMockOutputPort(ctrl)
		mp.EXPECT().Output(ac, wantOutputData)

		interactor := &Interactor{
			repository: m,
		}

		caseInput := InputData{
			UserID: caseEvent.UserID(),
		}
		interactor.SetEventOff(ac, caseInput, mp)
	})

	t.Run("ng:find", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()

		caseError := errors.New("error")

		caseEvent, _ := enterpriserule.NewTimerEvent("abc", "Hi!")
		caseEvent.State = enterpriserule.TimerEventStateWait

		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEvent(caseEvent.UserID()).
			Return(nil, caseError)

		wantOutputData := OutputData{
			Result:     fmt.Errorf("finding timer event error userID=%v: %w", caseEvent.UserID(), caseError),
			SavedEvent: nil,
		}

		mp := NewMockOutputPort(ctrl)
		mp.EXPECT().Output(ac, wantOutputData)

		interactor := &Interactor{
			repository: m,
		}

		caseInput := InputData{
			UserID: caseEvent.UserID(),
		}
		interactor.SetEventOff(ac, caseInput, mp)
	})

	t.Run("ng:find", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()

		var caseError error
		caseUserID := "test"

		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEvent(caseUserID).
			Return(nil, nil)

		wantOutputData := OutputData{
			Result:     fmt.Errorf("finding timer event error userID=%v: %w", caseUserID, caseError),
			SavedEvent: nil,
		}

		mp := NewMockOutputPort(ctrl)
		mp.EXPECT().Output(ac, wantOutputData)

		interactor := &Interactor{
			repository: m,
		}

		caseInput := InputData{
			UserID: caseUserID,
		}
		interactor.SetEventOff(ac, caseInput, mp)
	})

	t.Run("ng:save", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()

		caseEvent, _ := enterpriserule.NewTimerEvent("abc", "Hi!")
		caseEvent.State = enterpriserule.TimerEventStateWait

		caseError := errors.New("error")

		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEvent(caseEvent.UserID()).
			Return(caseEvent, nil)
		m.EXPECT().SaveTimerEvent(caseEvent).Return(nil, caseError)

		wantOutputData := OutputData{
			Result:     fmt.Errorf("saving timer event error userID=%v: %w", caseEvent.UserID(), caseError),
			SavedEvent: nil,
		}

		mp := NewMockOutputPort(ctrl)
		mp.EXPECT().Output(ac, wantOutputData)

		interactor := &Interactor{
			repository: m,
		}

		caseInput := InputData{
			UserID: caseEvent.UserID(),
		}
		interactor.SetEventOff(ac, caseInput, mp)
	})
}
