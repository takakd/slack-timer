package enqueue

import "slacktimer/internal/app/util/appcontext"

// ControllerHandler is called by Lambda handler.
type ControllerHandler interface {
	Handle(ac appcontext.AppContext, input HandleInput)
}

// HandleInput is input parameter of Controller.
type HandleInput struct {
	// Nothing currently
}
