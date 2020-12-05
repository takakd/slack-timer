package updatetimerevent

import (
	"context"
	"time"
)

type InputPort interface {
	// Set notificationTime to the notification time of the event which corresponds to userId.
	// Pass OutputPort interface if overwrite presenter implementation.
	//		e.g. HTTPResponse that needs http.ResponseWrite
	UpdateNotificationTime(ctx context.Context, userId string, notificationTime time.Time, presenter OutputPort)

	// Set notification interval to the event which corresponds to userId.
	// Use currentTime in calculating notification time if the event is not created.
	// Pass OutputPort interface if overwrite presenter implementation.
	//		e.g. HTTPResponse that needs http.ResponseWrite
	SaveIntervalMin(ctx context.Context, userId string, currentTime time.Time, minutes int, presenter OutputPort)
}
