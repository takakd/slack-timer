package notify

import (
	"context"
	"fmt"
	"slacktimer/internal/app/adapter/notify"
	"slacktimer/internal/app/util/appinitializer"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

// LambdaFunctor provides the method that is set to AWS Lambda.
type LambdaFunctor struct {
	ctrl notify.ControllerHandler
}

// NewLambdaFunctor create new struct.
func NewLambdaFunctor() *LambdaFunctor {
	h := &LambdaFunctor{}
	h.ctrl = di.Get("notify.ControllerHandler").(notify.ControllerHandler)
	return h
}

var _ LambdaHandler = (*LambdaFunctor)(nil)

// Handle is called by SQS.
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func (n LambdaFunctor) Handle(ctx context.Context, input LambdaInput) error {
	appinitializer.AppInit()

	log.Info(fmt.Sprintf("lambda handler input count=%d, recourds=%v", len(input.Records), input.Records))

	count := 0
	for _, m := range input.Records {
		resp := n.ctrl.Handle(ctx, m.HandleInput())
		if resp.Error != nil {
			count++
		}
	}

	var err error

	if count > 0 {
		err = fmt.Errorf("count=%d", count)
	}

	log.Info("lambda handler output", err)

	return err
}
