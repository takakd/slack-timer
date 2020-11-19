package enterpriserule

import (
	"github.com/pkg/errors"
	"reflect"
	"time"
)

// Holds the time to drink a timer event and the interval of drinking.
type TimerEvent struct {
	UserId           string    `dynamodbav:"UserId" db:"user_id" bson:"user_id"`
	NotificationTime time.Time `dynamodbav:"NotificationTime" db:"notification_time_utc" bson:"notification_time_utc"`
	IntervalMin      int       `dynamodbav:"IntervalMin" db:"interval_min" bson:"interval_min"`
}

func NewTimerEvent(userId string) (*TimerEvent, error) {
	if userId == "" {
		return nil, errors.New("must set userId")
	}

	p := &TimerEvent{
		UserId: userId,
	}
	return p, nil
}

func (p *TimerEvent) Equal(another *TimerEvent) bool {
	return reflect.DeepEqual(p, another)
}
