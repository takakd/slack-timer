package timeroffevent

import (
	"fmt"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/di"
	"testing"

	"time"

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
			mp := NewMockReplier(ctrl)
			md := di.NewMockDI(ctrl)
			md.EXPECT().Get("timeroffevent.Repository").Return(mr)
			md.EXPECT().Get("timeroffevent.Replier").Return(mp)
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

		caseEvent, _ := enterpriserule.NewTimerEvent("abc")
		caseEvent.Text = "Hi!"
		caseEvent.State = enterpriserule.TimerEventStateWait

		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEvent(caseEvent.UserID()).
			Return(caseEvent, nil)
		m.EXPECT().SaveTimerEvent(caseEvent).Return(nil, nil)

		want, _ := enterpriserule.NewTimerEvent(caseEvent.UserID())
		want.Text = caseEvent.Text
		want.NotificationTime = caseEvent.NotificationTime
		want.IntervalMin = caseEvent.IntervalMin
		want.State = enterpriserule.TimerEventStateDisabled

		wantOutputData := OutputData{
			Result:     nil,
			SavedEvent: want,
		}

		mp := NewMockReplier(ctrl)
		mp.EXPECT().SendMessage(ac, caseEvent.UserID(), ReplySuccess).Return(nil)

		mo := NewMockOutputPort(ctrl)
		mo.EXPECT().Output(ac, wantOutputData)

		interactor := &Interactor{
			repository: m,
			replier:    mp,
		}

		caseInput := InputData{
			UserID: caseEvent.UserID(),
		}
		interactor.SetEventOff(ac, caseInput, mo)
	})

	t.Run("ng:find error", func(t *testing.T) {
		cases := []struct {
			name   string
			userID string
		}{
			{"find error", "test"},
			{"find nil", ""},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				ac := appcontext.TODO()

				var caseError error
				if c.userID != "" {
					caseError = errors.New("error")
				}

				m := NewMockRepository(ctrl)
				m.EXPECT().FindTimerEvent(c.userID).
					Return(nil, caseError)

				wantOutputData := OutputData{
					Result:     fmt.Errorf("finding timer event error userID=%v: %w", c.userID, caseError),
					SavedEvent: nil,
				}

				mp := NewMockReplier(ctrl)
				if c.userID != "" {
					mp.EXPECT().SendMessage(ac, c.userID, ReplyFailure).Return(nil)
				}

				mo := NewMockOutputPort(ctrl)
				mo.EXPECT().Output(ac, wantOutputData)

				interactor := &Interactor{
					repository: m,
					replier:    mp,
				}

				caseInput := InputData{
					UserID: c.userID,
				}
				interactor.SetEventOff(ac, caseInput, mo)
			})
		}
	})

	t.Run("ng:save", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()

		caseEvent, _ := enterpriserule.NewTimerEvent("abc")
		caseEvent.NotificationTime = time.Now().UTC()
		caseEvent.IntervalMin = 10
		caseEvent.State = enterpriserule.TimerEventStateDisabled

		caseError := errors.New("error")

		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEvent(caseEvent.UserID()).
			Return(caseEvent, nil)
		m.EXPECT().SaveTimerEvent(caseEvent).Return(nil, caseError)

		want, _ := enterpriserule.NewTimerEvent(caseEvent.UserID())
		want.NotificationTime = caseEvent.NotificationTime
		want.IntervalMin = caseEvent.IntervalMin
		want.State = enterpriserule.TimerEventStateWait

		wantOutputData := OutputData{
			Result:     fmt.Errorf("saving timer event error userID=%v: %w", caseEvent.UserID(), caseError),
			SavedEvent: nil,
		}

		mp := NewMockReplier(ctrl)
		mp.EXPECT().SendMessage(ac, caseEvent.UserID(), ReplyFailure)

		mo := NewMockOutputPort(ctrl)
		mo.EXPECT().Output(ac, wantOutputData)

		interactor := &Interactor{
			repository: m,
			replier:    mp,
		}

		caseInput := InputData{
			UserID: caseEvent.UserID(),
		}
		interactor.SetEventOff(ac, caseInput, mo)
	})

	t.Run("ng:reply", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()

		caseUserID := "test"
		caseError := errors.New("error")

		mr := NewMockRepository(ctrl)
		mr.EXPECT().FindTimerEvent(caseUserID).Return(nil, nil)

		mp := NewMockReplier(ctrl)
		mp.EXPECT().SendMessage(ac, caseUserID, ReplyFailure).Return(caseError)

		wantOutputData := OutputData{
			SavedEvent: nil,
			Result:     fmt.Errorf("reply error userID=%v: %w", caseUserID, caseError),
		}
		mo := NewMockOutputPort(ctrl)
		mo.EXPECT().Output(ac, wantOutputData)

		interactor := &Interactor{
			repository: mr,
			replier:    mp,
		}

		caseInput := InputData{
			UserID: caseUserID,
		}
		interactor.SetEventOff(ac, caseInput, mo)
	})
}
