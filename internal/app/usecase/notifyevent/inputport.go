package notifyevent

import (
	"context"
)

type InputPort interface {
	// Notify an event to user and update entity.
	NotifyEvent(ctx context.Context, input *InputData) error
}

type InputData struct {
	UserId  string
	Message string
}
