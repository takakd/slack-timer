package settime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slacktimer/internal/app/adapter/settime"
	"slacktimer/internal/app/util/appinitializer"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"slacktimer/internal/pkg/helper"
)

type SetTimeLambdaHandler struct {
	ctrl settime.Controller
}

func NewSetTimeLambdaHandler() LambdaHandler {
	h := &SetTimeLambdaHandler{}
	h.ctrl = di.Get("settime.Controller").(settime.Controller)
	return h
}

// API Gateway calls this function.
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func (s SetTimeLambdaHandler) Handle(ctx context.Context, input LambdaInput) (*LambdaOutput, error) {
	appinitializer.AppInit()

	log.Info("lambda handler input", input)

	// Extract Slack event data.
	data, err := input.HandleInput()
	if err != nil {
		log.Info(fmt.Errorf("invalid request: %w", err))
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
			log.Info(fmt.Errorf("internal sesrver error: %w", err))
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
