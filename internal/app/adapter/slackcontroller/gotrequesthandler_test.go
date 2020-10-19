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

func TestTestGotRequestHandler_Handler(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
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
	})
}

func TestGotRequestOutputPort_Output(t *testing.T) {
	cases := []struct {
		name    string
		err     error
		msg     string
		wantErr error
	}{
		{name: "ng:find", err: updateproteinevent.ErrFind, msg: "failed to save event", wantErr: ErrSaveEvent},
		{name: "ng:create", err: updateproteinevent.ErrCreate, msg: "failed to save event", wantErr: ErrSaveEvent},
		{name: "ng:save", err: updateproteinevent.ErrSave, msg: "failed to save event", wantErr: ErrSaveEvent},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data := &updateproteinevent.OutputData{
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
		data := &updateproteinevent.OutputData{
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
