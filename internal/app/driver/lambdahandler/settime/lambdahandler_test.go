package settime

import (
	"encoding/json"
	"slacktimer/internal/app/adapter/settime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLambdaInput_HandleInput(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		data := settime.EventCallbackData{
			Token:  "test",
			TeamID: "test id",
			MessageEvent: settime.MessageEvent{
				Type:    "message",
				EventTs: "1234.0000001",
				User:    "YIG35ADg",
				Ts:      "1234.0000001",
				Text:    "message",
			},
			Challenge: "challenge",
		}
		dataJSON, _ := json.Marshal(&data)

		caseInput := LambdaInput{
			Body: string(dataJSON),
		}

		want := &settime.HandleInput{
			EventData: data,
		}

		got, err := caseInput.HandleInput()
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("ng", func(t *testing.T) {
		caseInput := LambdaInput{
			Body: "{syntax error",
		}

		got, err := caseInput.HandleInput()
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}
