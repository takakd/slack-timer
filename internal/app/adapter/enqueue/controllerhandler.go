package enqueue

import "github.com/aws/aws-lambda-go/lambdacontext"

// ControllerHandler is called by Lambda handler.
type ControllerHandler interface {
	Handle(ctx lambdacontext.LambdaContext, input HandleInput)
}

// HandleInput is input parameter of Controller.
type HandleInput struct {
	// Nothing currently
}
