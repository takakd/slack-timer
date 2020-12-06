// Package enterpriserule provides entities.
package enterpriserule

import (
	"errors"
	"reflect"
	"time"
)

var (
	ErrMustSetUserId = errors.New("must set userId")
)

// Holds notification properties and notification state.
type TimerEvent struct {
	UserId           string          `dynamodbav:"UserId" db:"user_id" bson:"user_id"`
	NotificationTime time.Time       `dynamodbav:"NotificationTime" db:"notification_time_utc" bson:"notification_time_utc"`
	IntervalMin      int             `dynamodbav:"IntervalMin" db:"interval_min" bson:"interval_min"`
	State            TimerEventState `dynamodbav:"State" db:"state" bson:"state"`
}

type TimerEventState string

const (
	_timerEventStateWait   TimerEventState = "wait"
	_timerEventStateQueued TimerEventState = "queued"
)

func NewTimerEvent(userId string) (*TimerEvent, error) {
	if userId == "" {
		return nil, ErrMustSetUserId
	}

	p := &TimerEvent{
		UserId: userId,
		State:  _timerEventStateWait,
	}
	return p, nil
}

func (p TimerEvent) Equal(another TimerEvent) bool {
	return reflect.DeepEqual(p, another)
}

func (p *TimerEvent) IncrementNotificationTime() {
	p.NotificationTime = p.NotificationTime.Add(time.Duration(p.IntervalMin) * time.Minute)
}

func (p TimerEvent) Queued() bool {
	return p.State == _timerEventStateQueued
}

func (p *TimerEvent) SetQueued() {
	p.State = _timerEventStateQueued
}

func (p *TimerEvent) SetWait() {
	p.State = _timerEventStateWait
}
