package notify

import (
	"context"
	"fmt"
	"slacktimer/internal/app/adapter/notify"
	"slacktimer/internal/app/util/appinitializer"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

type NotifyLambdaHandler struct {
	ctrl notify.Controller
}

func NewNotifyLambdaHandler() LambdaHandler {
	h := &NotifyLambdaHandler{}
	h.ctrl = di.Get("notify.Controller").(notify.Controller)
	return h
}

// SQS calls this function.
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func (n *NotifyLambdaHandler) Handle(ctx context.Context, input LambdaInput) error {
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
