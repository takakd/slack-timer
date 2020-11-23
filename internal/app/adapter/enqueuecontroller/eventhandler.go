package enqueuecontroller

import (
	"context"
	"slacktimer/internal/app/usecase/enqueueevent"
	"time"
)

// EventHandler handles CloudWatchEvent.
type CloudWatchEventHandler struct {
	usecase enqueueevent.Usecase
}

func (ce *CloudWatchEventHandler) Handler(ctx context.Context) *HandlerResponse {
	now := time.Now().UTC()
	err := ce.usecase.EnqueueEvent(ctx, now)
	resp := &HandlerResponse{
		Error: err,
	}
	return resp
}
