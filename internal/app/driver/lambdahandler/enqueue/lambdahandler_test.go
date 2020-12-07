package enqueue

import (
	"slacktimer/internal/app/adapter/enqueue"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLambdaInput_HandlerInput(t *testing.T) {
	caseInput := LambdaInput{}
	assert.Equal(t, enqueue.HandleInput{}, caseInput.HandleInput())
}
