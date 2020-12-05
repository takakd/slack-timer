package enqueue

import (
	"context"
	"slacktimer/internal/app/adapter/enqueuecontroller"
	"slacktimer/internal/app/util/appinit"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

type LambdaHandler interface {
	LambdaHandler(ctx context.Context, input LambdaInput)
}

type EnqueueLambdaHandler struct {
	ctrl enqueuecontroller.Handler
}

func NewEnqueueLambdaHandler() LambdaHandler {
	h := &EnqueueLambdaHandler{}
	h.ctrl = di.Get("enqueue.Handler").(enqueuecontroller.Handler)
	return h
}

// Lambda handler input data
// CloudWatchEvent passes this.
type LambdaInput struct {
	Version    string   `json:"version"`
	Id         string   `json:"id"`
	DetailType string   `json:"detail-type"`
	Source     string   `json:"source"`
	Account    string   `json:"account"`
	Time       string   `json:"time"`
	Region     string   `json:"region"`
	Resources  []string `json:"resources"`
	Detail     struct {
		EventCategories  []string `json:"EventCategories"`
		SourceType       string   `json:"SourceType"`
		SourceArn        string   `json:"SourceArn"`
		Date             string   `json:"Date"`
		Message          string   `json:"Message"`
		SourceIdentifier string   `json:"SourceIdentifier"`
	} `json:"detail"`
}

// To a controller input data.
func (s *LambdaInput) HandlerInput() enqueuecontroller.HandlerInput {
	return enqueuecontroller.HandlerInput{}
}

// CloudWatchEvent calls this function.
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func (e *EnqueueLambdaHandler) LambdaHandler(ctx context.Context, input LambdaInput) {
	appinit.AppInit()

	log.Info("handler called", input)

	e.ctrl.Handler(ctx, input.HandlerInput())

	log.Info("handler done")
}
