package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDynamoDbWrapperAdapter(t *testing.T) {
	assert.NotPanics(t, func() {
		NewDynamoDbWrapperAdapter()
	})
}
