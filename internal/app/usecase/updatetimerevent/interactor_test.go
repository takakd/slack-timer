package updatetimerevent

import (
	"errors"
	"fmt"
	"slacktimer/internal/app/enterpriserule"
	"testing"
	"time"

	"slacktimer/internal/app/util/appcontext"

	"slacktimer/internal/app/util/di"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewInteractor(t *testing.T) {
	assert.NotPanics(t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mr := NewMockRepository(ctrl)
		mp := NewMockReplier(ctrl)
		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("updatetimerevent.Repository").Return(mr)
		md.EXPECT().Get("updatetimerevent.Replier").Return(mp)
		di.SetDi(md)

		NewInteractor()
	})

	assert.Panics(t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("updatetimerevent.Repository").Return(nil)
		di.SetDi(md)

		NewInteractor()
	})
}

func TestInteractor_SaveIntervalMin(t *testing.T) {
	t.Run("ok:create", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()
		userID := "abc"
		caseTime := time.Now().UTC()

		caseEvent, _ := enterpriserule.NewTimerEvent(userID)
		caseEvent.IntervalMin = 10
		caseEvent.NotificationTime = caseTime
		caseEvent.Text = "Hi!"

		caseInput := SaveEventInputData{
			UserID:      caseEvent.UserID(),
			CurrentTime: caseEvent.NotificationTime,
			Minutes:     caseEvent.IntervalMin,
			Text:        caseEvent.Text,
		}

		want, _ := enterpriserule.NewTimerEvent(caseEvent.UserID())
		want.Text = caseEvent.Text
		want.IntervalMin = caseEvent.IntervalMin
		want.NotificationTime = caseTime.Add(time.Duration(want.IntervalMin) * time.Minute)
		wantOutput := OutputData{
			SavedEvent: want,
		}

		mr := NewMockRepository(ctrl)
		mr.EXPECT().FindTimerEvent(userID).Return(nil, nil)
		mr.EXPECT().SaveTimerEvent(gomock.Any()).DoAndReturn(func(event *enterpriserule.TimerEvent) (*enterpriserule.TimerEvent, error) {
			return event, nil
		})

		mp := NewMockReplier(ctrl)
		mp.EXPECT().SendMessage(ac, userID, ReplySuccess).Return(nil)

		mo := NewMockOutputPort(ctrl)
		mo.EXPECT().Output(ac, wantOutput)

		interactor := &Interactor{
			repository: mr,
			replier:    mp,
		}

		interactor.SaveIntervalMin(ac, caseInput, mo)
	})

	t.Run("ok:next notify", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()
		userID := "abc"
		caseTime := time.Now().UTC()

		caseEvent, _ := enterpriserule.NewTimerEvent(userID)
		caseEvent.NotificationTime = caseTime
		caseEvent.IntervalMin = 0
		caseEvent.Text = "Hi!"

		caseInput := SaveEventInputData{
			UserID:      caseEvent.UserID(),
			CurrentTime: caseEvent.NotificationTime,
			Minutes:     caseEvent.IntervalMin,
			Text:        caseEvent.Text,
		}

		want, _ := enterpriserule.NewTimerEvent(caseEvent.UserID())
		want.Text = "Hi!"
		want.IntervalMin = caseEvent.IntervalMin
		want.NotificationTime = caseTime.Add(time.Duration(want.IntervalMin) * time.Minute)
		wantOutput := OutputData{
			SavedEvent: want,
		}

		mr := NewMockRepository(ctrl)
		mr.EXPECT().FindTimerEvent(caseEvent.UserID()).
			Return(caseEvent, nil)
		mr.EXPECT().SaveTimerEvent(caseEvent).Return(nil, nil)

		mp := NewMockReplier(ctrl)
		mp.EXPECT().SendMessage(ac, userID, ReplySuccess).Return(nil)

		mo := NewMockOutputPort(ctrl)
		mo.EXPECT().Output(ac, wantOutput)

		interactor := &Interactor{
			repository: mr,
			replier:    mp,
		}

		interactor.SaveIntervalMin(ac, caseInput, mo)
	})

	t.Run("ng:find", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()
		caseUserID := "abc"
		caseError := errors.New("find error")

		caseInput := SaveEventInputData{
			UserID: caseUserID,
		}

		mr := NewMockRepository(ctrl)
		mr.EXPECT().FindTimerEvent(caseUserID).Return(nil, caseError)

		mp := NewMockReplier(ctrl)
		mp.EXPECT().SendMessage(ac, caseUserID, ReplyFailure).Return(nil)

		wantOutput := OutputData{
			Result: fmt.Errorf("finding timer event error userID=%v: %w", caseUserID, caseError),
		}

		mo := NewMockOutputPort(ctrl)
		mo.EXPECT().Output(ac, wantOutput)

		interactor := &Interactor{
			repository: mr,
			replier:    mp,
		}

		interactor.SaveIntervalMin(ac, caseInput, mo)
	})

	t.Run("ng:create", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()
		caseTime := time.Now().UTC()
		caseError := errors.New("must set userID")
		caseUserID := ""

		caseInput := SaveEventInputData{
			UserID:      caseUserID,
			CurrentTime: caseTime,
			Minutes:     0,
			Text:        "Hi!",
		}

		wantOutput := OutputData{
			Result: fmt.Errorf("creating timer event error userID=%v: %w", caseInput.UserID, caseError),
		}

		mr := NewMockRepository(ctrl)
		mr.EXPECT().FindTimerEvent(caseUserID).
			Return(nil, nil)

		mo := NewMockOutputPort(ctrl)
		mo.EXPECT().Output(ac, wantOutput)

		interactor := &Interactor{
			repository: mr,
			replier:    nil,
		}

		interactor.SaveIntervalMin(ac, caseInput, mo)
	})

	t.Run("ng:update", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()
		caseUserID := "abc"
		caseTime := time.Now().UTC()
		caseError := errors.New("save error")

		caseEvent, _ := enterpriserule.NewTimerEvent(caseUserID)
		caseEvent.NotificationTime = caseTime
		caseEvent.IntervalMin = 0
		caseEvent.Text = "Hi!"

		caseInput := SaveEventInputData{
			UserID:      caseEvent.UserID(),
			CurrentTime: caseEvent.NotificationTime,
			Minutes:     caseEvent.IntervalMin,
			Text:        caseEvent.Text,
		}

		wantOutput := OutputData{
			Result: fmt.Errorf("saving timer event error userID=%v: %w", caseEvent.UserID(), caseError),
		}

		mr := NewMockRepository(ctrl)
		mr.EXPECT().FindTimerEvent(caseEvent.UserID()).
			Return(caseEvent, nil)
		mr.EXPECT().SaveTimerEvent(caseEvent).Return(nil, caseError)

		mp := NewMockReplier(ctrl)
		mp.EXPECT().SendMessage(ac, caseUserID, ReplyFailure).Return(nil)

		mo := NewMockOutputPort(ctrl)
		mo.EXPECT().Output(ac, wantOutput)

		interactor := &Interactor{
			repository: mr,
			replier:    mp,
		}

		interactor.SaveIntervalMin(ac, caseInput, mo)
	})

	t.Run("ng:reply", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()
		caseTime := time.Now().UTC()

		caseEvent, _ := enterpriserule.NewTimerEvent("abc")
		caseEvent.IntervalMin = 10
		caseEvent.NotificationTime = caseTime
		caseEvent.Text = "Hi!"

		caseInput := SaveEventInputData{
			UserID:      caseEvent.UserID(),
			CurrentTime: caseEvent.NotificationTime,
			Minutes:     caseEvent.IntervalMin,
			Text:        caseEvent.Text,
		}

		caseError := errors.New("error")

		mr := NewMockRepository(ctrl)
		mr.EXPECT().FindTimerEvent(caseEvent.UserID()).Return(nil, nil)
		mr.EXPECT().SaveTimerEvent(gomock.Any()).DoAndReturn(func(event *enterpriserule.TimerEvent) (*enterpriserule.TimerEvent, error) {
			return event, nil
		})

		mp := NewMockReplier(ctrl)
		mp.EXPECT().SendMessage(ac, caseEvent.UserID(), ReplySuccess).Return(caseError)

		want, _ := enterpriserule.NewTimerEvent(caseEvent.UserID())
		want.Text = caseEvent.Text
		want.IntervalMin = caseEvent.IntervalMin
		want.NotificationTime = caseTime.Add(time.Duration(want.IntervalMin) * time.Minute)
		wantOutput := OutputData{
			SavedEvent: want,
			Result:     fmt.Errorf("reply error userID=%v: %w", caseInput.UserID, caseError),
		}

		mo := NewMockOutputPort(ctrl)
		mo.EXPECT().Output(ac, wantOutput)

		interactor := &Interactor{
			repository: mr,
			replier:    mp,
		}

		interactor.SaveIntervalMin(ac, caseInput, mo)
	})
}
