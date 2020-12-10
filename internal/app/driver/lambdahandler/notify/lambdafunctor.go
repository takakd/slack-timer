package notify

import (
	"context"
	"fmt"
	"slacktimer/internal/app/adapter/notify"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"time"
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
	ac, err := appcontext.NewLambdaAppContext(ctx, time.Now())
	if err != nil {
		log.Error("context error", err, ctx)
		return fmt.Errorf("context error: %w", err)
	}

	log.InfoWithContext(ac, "lambda handler", map[string]interface{}{
		"count":   len(input.Records),
		"records": input.Records,
	})

	count := 0
	for _, m := range input.Records {
		hi, err := m.HandleInput()
		if err != nil {
			count++
			continue
		}

		resp := n.ctrl.Handle(ac, hi)
		if resp.Error != nil {
			count++
		}
	}

	if count > 0 {
		err = fmt.Errorf("count=%d", count)
	}

	log.InfoWithContext(ac, "lambda handler output", err)

	return err
}
