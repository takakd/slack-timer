// Package enqueue the entry point of the enqueueing notification lambda function.
package enqueue

import (
	"context"
	"slacktimer/internal/app/adapter/enqueue"
)

type LambdaHandler interface {
	Handle(ctx context.Context, input LambdaInput)
}

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

// To input data for controller.
func (l LambdaInput) HandleInput() enqueue.HandleInput {
	return enqueue.HandleInput{}
}
