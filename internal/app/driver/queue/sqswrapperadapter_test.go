package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSqsWrapperAdapter(t *testing.T) {
	assert.NotPanics(t, func() {
		NewSqsWrapperAdapter()
	})
}
