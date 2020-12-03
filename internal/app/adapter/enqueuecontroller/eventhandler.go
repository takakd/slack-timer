package enqueuecontroller

import (
	"context"
	"fmt"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/log"
	"time"
)

// EventHandler handles CloudWatchEvent.
type CloudWatchEventHandler struct {
	usecase enqueueevent.Usecase
}

func (ce *CloudWatchEventHandler) Handler(ctx context.Context) *HandlerResponse {
	now := time.Now().UTC()

	log.Info(fmt.Sprintf("Usecase.EnqueueEvent time=%s", now))

	ce.usecase.EnqueueEvent(ctx, now)
	resp := &HandlerResponse{
		// TODO: modify structures along notification feature
		Error: nil,
	}

	log.Info(fmt.Sprintf("Usecase.EnqueueEvent output=%s", *resp))

	return resp
}
