package notify

import (
	"testing"

	"encoding/json"
	"slacktimer/internal/app/driver/queue"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSqsMessage_HandleInput(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		c := queue.SqsMessageBody{
			UserID: "test user",
			Text:   "test text",
		}
		cBody, err := json.Marshal(c)
		require.NoError(t, err)

		m := &SqsMessage{
			Body: string(cBody),
		}
		h, err := m.HandleInput()

		assert.NoError(t, err)
		assert.Equal(t, c.UserID, h.UserID)
		assert.Equal(t, c.Text, h.Message)
	})

	t.Run("ng:marshal", func(t *testing.T) {
		m := &SqsMessage{
			Body: string("{syntax error"),
		}
		h, err := m.HandleInput()

		assert.Empty(t, h)
		assert.Error(t, err)
	})
}
