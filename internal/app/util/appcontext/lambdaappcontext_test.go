package appcontext

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stretchr/testify/assert"
)

func TestNewLambdaContext(t *testing.T) {
	assert.NotPanics(t, func() {
		lc := &lambdacontext.LambdaContext{
			AwsRequestID: "test ID",
		}
		ctx := lambdacontext.NewContext(context.TODO(), lc)
		FromContext(ctx)
	})
	assert.Panics(t, func() {
		FromContext(nil)
	})
}

func TestLambdaAppContext_RequestID(t *testing.T) {
	ac := &LambdaAppContext{
		requestID: "test ID",
	}
	assert.Equal(t, ac.requestID, ac.RequestID())
}
