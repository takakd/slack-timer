// Package updatetimerevent provides usecase that finding or saving event entity.
package updatetimerevent

import (
	"slacktimer/internal/app/util/appcontext"
	"time"
)

// InputPort defines updating timer usecase.
type InputPort interface {
	// SaveIntervalMin sets notification interval to the event which corresponds to userID.
	// Use currentTime in calculating notification time if the event is not created.
	// Pass OutputPort interface if overwrite presenter implementation.
	//		e.g. HTTPResponse that needs http.ResponseWrite
	SaveIntervalMin(ac appcontext.AppContext, input SaveEventInputData, presenter OutputPort)
}

// UpdateNotificationTimeInputData is parameter of UpdateNotificationTimeInputData .
type UpdateNotificationTimeInputData struct {
	UserID           string
	NotificationTime time.Time
}

// SaveEventInputData is parameter of SaveEventInputData.
type SaveEventInputData struct {
	UserID      string
	CurrentTime time.Time
	Minutes     int
	Text        string
}
