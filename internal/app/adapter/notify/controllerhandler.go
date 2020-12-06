// Package notify provides features that notify events to user.
package notify

import (
	"context"
)

// ControllerHandler is called by Lambda handler.
type ControllerHandler interface {
	Handle(ctx context.Context, input HandleInput) *Response
}

// Response is returns of Controller.
type Response struct {
	Error error
}

// HandleInput is input parameter of Controller.
type HandleInput struct {
	// Notify users identified this ID.
	UserID string
	// Notified message
	Message string
}
