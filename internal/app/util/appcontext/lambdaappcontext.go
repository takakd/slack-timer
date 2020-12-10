package appcontext

import (
	"context"

	"time"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/pkg/errors"
)

// LambdaAppContext implements AppContext with lambdacontext.LambdaContext.
type LambdaAppContext struct {
	requestID string
	called    time.Time
}

// NewLambdaAppContext creates new struct.
func NewLambdaAppContext(ctx context.Context, called time.Time) (*LambdaAppContext, error) {
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		return nil, errors.New("context error")
	}

	return &LambdaAppContext{
		requestID: lc.AwsRequestID,
		called:    called.UTC(),
	}, nil
}

// TODO returns empty context, which is only used in unit test.
func TODO() *LambdaAppContext {
	return &LambdaAppContext{}
}

// RequestID returns the current request ID.
func (r *LambdaAppContext) RequestID() string {
	return r.requestID
}

// HandlerCalledTime returns time of calling handler.
func (r *LambdaAppContext) HandlerCalledTime() time.Time {
	return r.called
}
