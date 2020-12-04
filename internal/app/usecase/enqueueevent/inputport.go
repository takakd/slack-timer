package enqueueevent

import (
	"context"
	"time"
)

type InputPort interface {
	// Enqueue notification event, which notification time overs eventTime.
	EnqueueEvent(ctx context.Context, data InputData)
}

type InputData struct {
	EventTime time.Time
}
