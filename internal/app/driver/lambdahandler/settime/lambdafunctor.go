package settime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slacktimer/internal/app/adapter/settime"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"slacktimer/internal/pkg/helper"
)

// LambdaFunctor provides the method that is set to AWS Lambda.
type LambdaFunctor struct {
	ctrl settime.ControllerHandler
}

// NewLambdaFunctor create new struct.
func NewLambdaFunctor() *LambdaFunctor {
	h := &LambdaFunctor{}
	h.ctrl = di.Get("settime.ControllerHandler").(settime.ControllerHandler)
	return h
}

var _ LambdaHandler = (*LambdaFunctor)(nil)

// Handle is called by API Gateway.
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func (s LambdaFunctor) Handle(ctx context.Context, input LambdaInput) (*LambdaOutput, error) {
	log.Info("lambda handler input", input)

	// Extract Slack event data.
	data, err := input.HandleInput()
	if err != nil {
		log.Info("invalid request", err)
		o := &LambdaOutput{
			IsBase64Encoded: true,
			StatusCode:      http.StatusInternalServerError,
			// Not have to do use json.Marshal.
			Body: `{"message":"invalid request", "detail":"parameters are wrong"}`,
		}
		return o, nil
	}

	// Call controller method.
	// TOOD: di
	resp := s.ctrl.Handle(ctx, *data)

	// Create a response.
	var respBody string
	if helper.IsStruct(resp.Body) {
		body, err := json.Marshal(resp.Body)
		if err != nil {
			log.Info("internal sesrver error", err)
			return nil, errors.New("internal server error")
		}
		respBody = string(body)
	} else {
		respBody = fmt.Sprintf("%v", resp.Body)
	}

	output := &LambdaOutput{
		//IsBase64Encoded: resp.IsBase64Encoded,
		// TODO: check
		IsBase64Encoded: true,
		StatusCode:      resp.StatusCode,
		Body:            respBody,
	}

	log.Info("handler output", output)

	return output, nil
}
