package notify

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSqsMessage_HandleInput(t *testing.T) {
	m := &SqsMessage{
		Body: "test user",
	}
	h := m.HandleInput()
	assert.Equal(t, m.Body, h.UserId)
	assert.Equal(t, "test user", h.UserId)
}
