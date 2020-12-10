// Package updatetimerevent provides usecase that finding or saving event entity.
package updatetimerevent

import (
	"slacktimer/internal/app/util/appcontext"
	"time"
)

// InputPort defines updating timer usecase.
type InputPort interface {
	// UpdateNotificationTime sets notificationTime to the notification time of the event which corresponds to userID.
	// Pass OutputPort interface if overwrite presenter implementation.
	//		e.g. HTTPResponse that needs http.ResponseWrite
	UpdateNotificationTime(ac appcontext.AppContext, userID string, notificationTime time.Time, presenter OutputPort)

	// SaveIntervalMin sets notification interval to the event which corresponds to userID.
	// Use currentTime in calculating notification time if the event is not created.
	// Pass OutputPort interface if overwrite presenter implementation.
	//		e.g. HTTPResponse that needs http.ResponseWrite
	SaveIntervalMin(ac appcontext.AppContext, userID string, currentTime time.Time, minutes int, presenter OutputPort)
}
