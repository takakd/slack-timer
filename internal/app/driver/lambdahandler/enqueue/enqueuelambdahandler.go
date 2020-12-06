package enqueue

import (
	"context"
	"slacktimer/internal/app/adapter/enqueue"
	"slacktimer/internal/app/util/appinitializer"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

type EnqueueLambdaHandler struct {
	ctrl enqueue.Controller
}

var _ LambdaHandler = (*EnqueueLambdaHandler)(nil)

func NewEnqueueLambdaHandler() *EnqueueLambdaHandler {
	h := &EnqueueLambdaHandler{}
	h.ctrl = di.Get("enqueue.Controller").(enqueue.Controller)
	return h
}

// CloudWatchEvent calls this function.
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func (e EnqueueLambdaHandler) Handle(ctx context.Context, input LambdaInput) {
	appinitializer.AppInit()

	log.Info("handler called", input)

	e.ctrl.Handle(ctx, input.HandleInput())

	log.Info("handler done")
}
