package enterpriserule

import (
	"github.com/pkg/errors"
	"time"
)

// Holds the time to drink a protein and the interval of drinking.
type ProteinEvent struct {
	UserId               string        `db:"user_id" bson:"user_id"`
	UtcTimeToDrink       time.Time     `db:"utc_time_to_drink" bson:"utc_time_to_drink"`
	DrinkTimeIntervalSec time.Duration `db:"drink_time_interval_sec" bson:"drink_time_interval_sec"`
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
	if p == nil || another == nil {
		return false
	}
	if p.UserId != another.UserId {
		return false
	}
	if p.DrinkTimeIntervalSec != another.DrinkTimeIntervalSec {
		return false
	}
	if p.UtcTimeToDrink.Second() != another.UtcTimeToDrink.Second() {
		return false
	}
	return true
}
