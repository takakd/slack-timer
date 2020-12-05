package notify

import (
	"context"
)

type Controller interface {
	Handle(ctx context.Context, input HandleInput) *Response
}

type Response struct {
	Error error
}

type HandleInput struct {
	UserId  string
	Message string
}
