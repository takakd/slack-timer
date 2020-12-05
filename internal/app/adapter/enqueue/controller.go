// Package enqueuecontroller provides that events reached the time enqueue queue.
package enqueue

import (
	"context"
)

type Controller interface {
	Handle(ctx context.Context, input HandleInput)
}

type HandleInput struct {
}
