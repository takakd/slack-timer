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

	log.Info("")

	log.Info(fmt.Sprintf("Usecase.EnqueueEvent now=%s ", now))

	err := ce.usecase.EnqueueEvent(ctx, now)
	resp := &HandlerResponse{
		Error: err,
	}

	log.Info(fmt.Sprintf("Usecase.EnqueueEvent output=%s ", resp))

	return resp
}
