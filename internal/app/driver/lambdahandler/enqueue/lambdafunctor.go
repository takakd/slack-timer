package enqueue

import (
	"context"
	"fmt"
	"slacktimer/internal/app/adapter/enqueue"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"time"
)

// LambdaFunctor provides the method that is set to AWS Lambda.
type LambdaFunctor struct {
	ctrl enqueue.ControllerHandler
}

var _ LambdaHandler = (*LambdaFunctor)(nil)

// NewLambdaFunctor create new struct.
func NewLambdaFunctor() *LambdaFunctor {
	h := &LambdaFunctor{}
	h.ctrl = di.Get("enqueue.ControllerHandler").(enqueue.ControllerHandler)
	return h
}

// Handle is called by CloudWatchEvent.
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func (e LambdaFunctor) Handle(ctx context.Context, input LambdaInput) error {
	ac, err := appcontext.NewLambdaAppContext(ctx, time.Now())
	if err != nil {
		log.Error("context error", err, ctx)
		return fmt.Errorf("context error: %w", err)
	}

	log.InfoWithContext(ac, "handler called", input)

	e.ctrl.Handle(ac, input.HandleInput())

	log.InfoWithContext(ac, "handler done")

	return nil
}
