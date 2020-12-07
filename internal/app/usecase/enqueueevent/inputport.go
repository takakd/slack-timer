// Package enqueueevent provides usecase that enqueue event.
package enqueueevent

import (
	"context"
	"time"
)

// InputPort defines enqueueing events usecase.
type InputPort interface {
	// Enqueue notification event, which notification time overs eventTime.
	EnqueueEvent(ctx context.Context, data InputData)
}

// InputData is parameter of InputPort.
type InputData struct {
	EventTime time.Time
}
