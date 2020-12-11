package timeronevent

import (
	"slacktimer/internal/app/enterpriserule"
)

// Repository defines repository methods used in updating timer events usecase.
type Repository interface {
	FindTimerEvent(userID string) (*enterpriserule.TimerEvent, error)
	SaveTimerEvent(event *enterpriserule.TimerEvent) (*enterpriserule.TimerEvent, error)
}
