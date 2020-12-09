// Package notifyevent provides usecase that notify an event to the user.
package notifyevent

import (
	"slacktimer/internal/app/util/appcontext"
)

// InputPort defines notifying events usecase.
type InputPort interface {
	// Notify an event to the user and update the entity.
	NotifyEvent(ac appcontext.AppContext, input InputData) error
}

// InputData is parameter of InputPort.
type InputData struct {
	UserID  string
	Message string
}
