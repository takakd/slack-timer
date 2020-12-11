// Package enterpriserule provides entities.
package enterpriserule

import (
	"errors"
	"reflect"
	"time"
)

var (
	// ErrMustSetUserID returns if UserID parameter is empty.
	ErrMustSetUserID = errors.New("must set userID")
	// ErrMustSetText returns if Text parameter is empty.
	ErrMustSetText = errors.New("must set text")
)

// TimerEvent holds notification properties and notification state.
type TimerEvent struct {
	userID           string
	Text             string
	NotificationTime time.Time
	IntervalMin      int
	State            TimerEventState
}

// TimerEventState represents the type of Queueing state.
// The details of states are below:
// 		TimerEventStateWait: wait for enqueueing.
// 		TimerEventStateQueued: wait for notifying.
// 		TimerEventStateDisabled: disabled, not enqueueing and notifying.
type TimerEventState string

const (
	// TimerEventStateWait represents the state of waiting for queueing.
	TimerEventStateWait TimerEventState = "wait"
	// TimerEventStateQueued represents the state enqueued.
	TimerEventStateQueued TimerEventState = "queued"
	// TimerEventStateDisabled represents the disabled state.
	TimerEventStateDisabled TimerEventState = "disabled"
)

// NewTimerEvent creates new struct.
func NewTimerEvent(userID string) (*TimerEvent, error) {
	if userID == "" {
		return nil, ErrMustSetUserID
	}

	p := &TimerEvent{
		userID: userID,
		State:  TimerEventStateWait,
	}
	return p, nil
}

// Equal returns whether p is equal to another.
func (p TimerEvent) Equal(another TimerEvent) bool {
	return reflect.DeepEqual(p, another)
}

// IncrementNotificationTime increment a notification time with IntervalMin.
func (p *TimerEvent) IncrementNotificationTime() {
	p.NotificationTime = p.NotificationTime.Add(time.Duration(p.IntervalMin) * time.Minute)
}

// UserID returns the ID of having this event.
func (p TimerEvent) UserID() string {
	return p.userID
}
