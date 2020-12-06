// Package notifyevent provides usecase that notify an event to the user.
package notifyevent

import (
	"context"
)

type InputPort interface {
	// Notify an event to the user and update the entity.
	NotifyEvent(ctx context.Context, input InputData) error
}

type InputData struct {
	UserId  string
	Message string
}
