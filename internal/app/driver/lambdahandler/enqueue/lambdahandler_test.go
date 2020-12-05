package enqueue

import (
	"github.com/stretchr/testify/assert"
	"slacktimer/internal/app/adapter/enqueue"
	"testing"
)

func TestLambdaInput_HandlerInput(t *testing.T) {
	caseInput := LambdaInput{}
	assert.Equal(t, enqueue.HandleInput{}, caseInput.HandleInput())
}
