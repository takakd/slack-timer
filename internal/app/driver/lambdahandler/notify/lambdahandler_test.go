package notify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqsMessage_HandleInput(t *testing.T) {
	m := &SqsMessage{
		Body: "test user",
	}
	h := m.HandleInput()
	assert.Equal(t, m.Body, h.UserID)
	assert.Equal(t, "test user", h.UserID)
}
