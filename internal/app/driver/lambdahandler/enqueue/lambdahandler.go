package enqueue

import (
	"context"
	"slacktimer/internal/app/adapter/enqueue"
)

type LambdaHandler interface {
	Handle(ctx context.Context, input LambdaInput)
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
func (s *LambdaInput) HandleInput() enqueue.HandleInput {
	return enqueue.HandleInput{}
}
