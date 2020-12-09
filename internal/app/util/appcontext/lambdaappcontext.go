package appcontext

import (
	"context"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/pkg/errors"
)

// LambdaAppContext implements AppContext with lambdacontext.LambdaContext.
type LambdaAppContext struct {
	requestID string
}

// FromContext creates new struct.
func FromContext(ctx context.Context) (*LambdaAppContext, error) {
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		return nil, errors.New("context error")
	}

	return &LambdaAppContext{
		requestID: lc.AwsRequestID,
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
