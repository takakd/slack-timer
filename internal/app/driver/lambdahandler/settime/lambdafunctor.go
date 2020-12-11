package settime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slacktimer/internal/app/adapter/settime"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"slacktimer/internal/pkg/helper"
	"time"
)

// LambdaFunctor provides the method that is set to AWS Lambda.
type LambdaFunctor struct {
	ctrl settime.ControllerHandler
}

// NewLambdaFunctor creates new struct.
func NewLambdaFunctor() *LambdaFunctor {
	h := &LambdaFunctor{}
	h.ctrl = di.Get("settime.ControllerHandler").(settime.ControllerHandler)
	return h
}

var _ LambdaHandler = (*LambdaFunctor)(nil)

// Handle is called by API Gateway.
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func (s LambdaFunctor) Handle(ctx context.Context, input LambdaInput) (*LambdaOutput, error) {
	ac, err := appcontext.NewLambdaAppContext(ctx, time.Now())
	if err != nil {
		log.Error("context error", err, ctx)
		o := &LambdaOutput{
			IsBase64Encoded: false,
			StatusCode:      http.StatusInternalServerError,
			// Not have to do use json.Marshal.
			Body: fmt.Sprintf(`{"message":"context error", "detail":"%s"}`, err),
		}
		return o, fmt.Errorf("context error: %w", err)
	}

	log.InfoWithContext(ac, "lambda handler input", input)

	// Extract Slack event data.
	data, err := input.HandleInput()
	if err != nil {
		log.InfoWithContext(ac, "invalid request", err)
		o := &LambdaOutput{
			IsBase64Encoded: false,
			StatusCode:      http.StatusInternalServerError,
			// Not have to do use json.Marshal.
			Body: `{"message":"invalid request", "detail":"parameters are wrong"}`,
		}
		return o, nil
	}

	// Call controller method.
	resp := s.ctrl.Handle(ac, *data)

	// Create a response.
	var respBody string
	if helper.IsStruct(resp.Body) {
		body, err := json.Marshal(resp.Body)
		if err != nil {
			log.InfoWithContext(ac, "internal sesrver error", err)
			return nil, errors.New("internal server error")
		}
		respBody = string(body)
	} else {
		respBody = fmt.Sprintf("%v", resp.Body)
	}

	output := &LambdaOutput{
		IsBase64Encoded: false,
		StatusCode:      resp.StatusCode,
		Body:            respBody,
	}

	log.InfoWithContext(ac, "handler output", output)

	return output, nil
}
