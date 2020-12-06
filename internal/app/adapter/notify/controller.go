// Package notify provides features that notify events to user.
package notify

import (
	"context"
)

// Called by Lambda handler.
type Controller interface {
	Handle(ctx context.Context, input HandleInput) *Response
}

type Response struct {
	Error error
}

type HandleInput struct {
	// Notify users identified this ID.
	UserId  string
	// Notified message
	Message string
}
