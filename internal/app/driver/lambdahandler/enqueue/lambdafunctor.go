package enqueue

import (
	"context"
	"slacktimer/internal/app/adapter/enqueue"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"github.com/aws/aws-lambda-go/lambdacontext"
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
func (e LambdaFunctor) Handle(ctx context.Context, input LambdaInput) {
	lc, _ = lambdacontext.FromContext(ctx)
	log.Info("handler called", input)

	e.ctrl.Handle(lc, input.HandleInput())

	log.Info("handler done")
}
