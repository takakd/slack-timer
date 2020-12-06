// Package enqueue provides features that enqueue events reached the time.
package enqueue

import (
	"context"
)

// Called by Lambda handler.
type Controller interface {
	Handle(ctx context.Context, input HandleInput)
}

type HandleInput struct {
	// Nothing currently
}
