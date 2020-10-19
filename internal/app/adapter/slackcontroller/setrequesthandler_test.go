package slackcontroller

import (
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"proteinreminder/internal/app/usecase/updateproteinevent"
	"testing"
)

func TestSetRequestHandler_Handler(t *testing.T) {
	params := &SlackCallbackRequestParams{
		UserId: "test",
	}
	ctx := context.TODO()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mu := updateproteinevent.NewMockUsecase(ctrl)
	mu.EXPECT().UpdateTimeToDrink(gomock.Eq(ctx), gomock.Eq(params.UserId), gomock.Any())

	h := GotRequestHandler{
		params:  params,
		usecase: mu,
	}
	h.Handler(ctx, httptest.NewRecorder())
}

func TestSetRequestHandler_validate(t *testing.T) {
	cases := []struct {
		name  string
		text  string
		min   int
		valid bool
	}{
		{"ok", "set 10", 10, true},
		{"ok", "set 1", 1, true},
		{"ng", "set -1", 0, false},
		{"ng", "set", 0, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := SetRequestHandler{
				params: &SlackCallbackRequestParams{
					UserId: "test",
					Text:   c.text,
				},
			}
			bag := r.validate()
			_, exists := bag.GetError("interval")
			assert.Equal(t, c.valid, !exists)
			if c.valid {
				assert.Equal(t, c.min, r.remindIntervalInMin)
			}
		})
	}
}

func TestSetRequestOutputPort_Output(t *testing.T) {
	cases := []struct {
		name string
		err  error
		msg  string
	}{
		{name: "ng:find", err: updateproteinevent.ErrFind, msg: "failed to find event"},
		{name: "ng:create", err: updateproteinevent.ErrCreate, msg: "failed to create event"},
		{name: "ng:save", err: updateproteinevent.ErrSave, msg: "failed to save event"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data := &updateproteinevent.OutputData{
				Result: c.err,
			}

			w := httptest.NewRecorder()
			outputPort := &SetRequestOutputPort{
				w: w,
			}

			outputPort.Output(data)
			body, err := makeErrorCallbackResponseBody("failed to save event", ErrSaveEvent)

			assert.Equal(t, w.Code, http.StatusBadRequest)
			assert.NoError(t, err)
			assert.Equal(t, w.Body.Bytes(), body)
		})
	}

	t.Run("ok", func(t *testing.T) {
		data := &updateproteinevent.OutputData{
			Result: nil,
		}

		w := httptest.NewRecorder()
		outputPort := &SetRequestOutputPort{
			w: w,
		}

		outputPort.Output(data)

		want, _ := json.Marshal(SlackCallbackResponse{Message: "success"})
		assert.Equal(t, w.Code, http.StatusOK)
		assert.Equal(t, w.Body.Bytes(), want)
	})
}