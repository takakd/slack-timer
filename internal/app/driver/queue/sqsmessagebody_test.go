package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSqsMessageBody(t *testing.T) {
	assert.NotPanics(t, func() {
		NewSqsMessageBody()
	})
}
