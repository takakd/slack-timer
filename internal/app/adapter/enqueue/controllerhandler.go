package enqueue

import (
	"context"
)

// ControllerHandler is called by Lambda handler.
type ControllerHandler interface {
	Handle(ctx context.Context, input HandleInput)
}

// HandleInput is input parameter of Controller.
type HandleInput struct {
	// Nothing currently
}
