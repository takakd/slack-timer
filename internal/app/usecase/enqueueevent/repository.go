package enqueueevent

import (
	"slacktimer/internal/app/enterpriserule"
	"time"
)

// Repository defines repository methods used in enqueueing usecase.
type Repository interface {
	FindTimerEventsByTime(eventTime time.Time) ([]*enterpriserule.TimerEvent, error)
	SaveTimerEvent(event *enterpriserule.TimerEvent) (saved *enterpriserule.TimerEvent, err error)
}
