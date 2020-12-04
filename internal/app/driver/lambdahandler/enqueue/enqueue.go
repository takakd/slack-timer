package enqueue

import (
	"context"
	"slacktimer/internal/app/adapter/enqueuecontroller"
	"slacktimer/internal/app/util/appinit"
	"slacktimer/internal/app/util/log"
)

// Lambda handler input data
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

func (s *LambdaInput) HandlerInput() enqueuecontroller.HandlerInput {
	return enqueuecontroller.HandlerInput{}
}

// Lambda callback
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func LambdaHandleEvent(ctx context.Context, input LambdaInput) error {
	appinit.AppInit()

	log.Info("handler called", input)

	h := enqueuecontroller.NewHandler()
	resp := h.Handler(ctx, input.HandlerInput())

	log.Info("handler output", *resp)

	return resp.Error
}
