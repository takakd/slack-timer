// Package notifyevent provides usecase that notify an event to the user.
package notifyevent

import (
	"context"
)

// InputPort defines notifying events usecase.
type InputPort interface {
	// Notify an event to the user and update the entity.
	NotifyEvent(ctx context.Context, input InputData) error
}

// InputData is parameter of InputPort.
type InputData struct {
	UserID  string
	Message string
}
