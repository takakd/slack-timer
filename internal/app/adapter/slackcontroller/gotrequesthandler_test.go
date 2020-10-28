package slackcontroller

import (
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"testing"
)

func TestGotRequestHandler_Handler(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		data := &EventCallbackData{
			MessageEvent: MessageEvent{
				User: "test",
			},
		}
		ctx := context.TODO()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mu := updatetimerevent.NewMockUsecase(ctrl)
		mu.EXPECT().UpdateNotificationTime(gomock.Eq(ctx), gomock.Eq(data.MessageEvent.User), gomock.Any())

		h := GotRequestHandler{
			messageEvent: &data.MessageEvent,
			usecase:      mu,
		}
		h.Handler(ctx, httptest.NewRecorder())
	})
}

func TestGotRequestOutputPort_Output(t *testing.T) {
	cases := []struct {
		name    string
		err     error
		msg     string
		wantErr error
	}{
		{name: "ng:find", err: updatetimerevent.ErrFind, msg: "failed to save event", wantErr: ErrSaveEvent},
		{name: "ng:create", err: updatetimerevent.ErrCreate, msg: "failed to save event", wantErr: ErrSaveEvent},
		{name: "ng:save", err: updatetimerevent.ErrSave, msg: "failed to save event", wantErr: ErrSaveEvent},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data := &updatetimerevent.OutputData{
				Result: c.err,
			}

			w := httptest.NewRecorder()
			outputPort := &GotRequestOutputPort{
				w: w,
			}

			outputPort.Output(data)

			body, err := makeErrorCallbackResponseBody(c.msg, c.wantErr)
			assert.NoError(t, err)
			assert.Equal(t, w.Code, http.StatusBadRequest)
			assert.Equal(t, w.Body.Bytes(), body)
		})
	}

	t.Run("ok", func(t *testing.T) {
		data := &updatetimerevent.OutputData{
			Result: nil,
		}

		w := httptest.NewRecorder()
		outputPort := &GotRequestOutputPort{
			w: w,
		}

		outputPort.Output(data)

		want, _ := json.Marshal(SlackCallbackResponse{Message: "success"})
		assert.Equal(t, w.Code, http.StatusOK)
		assert.Equal(t, w.Body.Bytes(), want)
	})
}
