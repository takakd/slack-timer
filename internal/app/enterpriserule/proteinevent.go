package enterpriserule

import (
	"github.com/pkg/errors"
	"reflect"
	"time"
)

// Holds the time to drink a protein and the interval of drinking.
type ProteinEvent struct {
	UserId               string    `db:"user_id" bson:"user_id"`
	UtcTimeToDrink       time.Time `db:"utc_time_to_drink" bson:"utc_time_to_drink"`
	DrinkTimeIntervalMin int       `db:"drink_time_interval_min" bson:"drink_time_interval_min"`
}

func NewProteinEvent(userId string) (*ProteinEvent, error) {
	if userId == "" {
		return nil, errors.New("must set userId")
	}

	p := &ProteinEvent{
		UserId: userId,
	}
	return p, nil
}

func (p *ProteinEvent) Equal(another *ProteinEvent) bool {
	return reflect.DeepEqual(p, another)
}
