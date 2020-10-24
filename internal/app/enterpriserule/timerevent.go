package enterpriserule

import (
	"github.com/pkg/errors"
	"reflect"
	"time"
)

// Holds the time to drink a timer event and the interval of drinking.
type TimerEvent struct {
	UserId               string    `db:"user_id" bson:"user_id"`
	UtcTimeToDrink       time.Time `db:"utc_time_to_drink" bson:"utc_time_to_drink"`
	DrinkTimeIntervalMin int       `db:"drink_time_interval_min" bson:"drink_time_interval_min"`
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
