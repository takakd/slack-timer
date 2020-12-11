// Package timeronevent provides usecase that start event notification.
package timeronevent

import (
	"slacktimer/internal/app/util/appcontext"
)

// InputPort defines updating timer usecase.
type InputPort interface {
	// SetEventOn start to notify event to user which corresponds to userID.
	// Pass OutputPort interface if overwrite presenter implementation.
	//		e.g. HTTPResponse that needs http.ResponseWrite
	SetEventOn(ac appcontext.AppContext, input InputData, presenter OutputPort)
}

// InputData is parameter of UpdateNotificationTimeInputData .
type InputData struct {
	UserID string
}
