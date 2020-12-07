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
)

// TimerEvent holds notification properties and notification state.
type TimerEvent struct {
	UserID           string          `dynamodbav:"UserId" db:"user_id" bson:"user_id"`
	NotificationTime time.Time       `dynamodbav:"NotificationTime" db:"notification_time_utc" bson:"notification_time_utc"`
	IntervalMin      int             `dynamodbav:"IntervalMin" db:"interval_min" bson:"interval_min"`
	State            TimerEventState `dynamodbav:"State" db:"state" bson:"state"`
}

// TimerEventState represents the type of Queueing state.
type TimerEventState string

const (
	_timerEventStateWait   TimerEventState = "wait"
	_timerEventStateQueued TimerEventState = "queued"
)

// NewTimerEvent create new struct.
func NewTimerEvent(userID string) (*TimerEvent, error) {
	if userID == "" {
		return nil, ErrMustSetUserID
	}

	p := &TimerEvent{
		UserID: userID,
		State:  _timerEventStateWait,
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

// Queued returns whether this event is queued.
func (p TimerEvent) Queued() bool {
	return p.State == _timerEventStateQueued
}

// SetQueued sets this event state to queued state.
func (p *TimerEvent) SetQueued() {
	p.State = _timerEventStateQueued
}

// SetWait sets this event state to waiting state.
func (p *TimerEvent) SetWait() {
	p.State = _timerEventStateWait
}
