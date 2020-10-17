package slackcontroller

import (
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"proteinreminder/internal/app/usecase"
	"testing"
)

func TestGotRequestHandler_Handler(t *testing.T) {
	// Protein save errors.
	cases := []struct {
		name string
		err  error
		msg  string
	}{
		{name: "save error: find", err: usecase.ErrFind, msg: "failed to find event"},
		{name: "save error: create", err: usecase.ErrCreate, msg: "failed to create event"},
		{name: "save error: save", err: usecase.ErrSave, msg: "failed to save event"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.TODO()
			userId := "test user"
			ctrl := gomock.NewController(t)
			saver := usecase.NewMockProteinEventSaver(ctrl)
			saver.EXPECT().UpdateTimeToDrink(gomock.Eq(ctx), gomock.Eq(userId)).Return(c.err)

			h := &GotRequestHandler{
				saver: saver,
				params: &SlackCallbackRequestParams{
					UserId: userId,
				},
			}

			w := httptest.NewRecorder()
			h.Handler(context.TODO(), w)
			assert.Equal(t, w.Code, http.StatusBadRequest)
			assert.Equal(t, w.Body.Bytes(), makeErrorCallbackResponseBody(c.msg, c.err))
		})
	}

	// success
	t.Run("success", func(t *testing.T) {
		ctx := context.TODO()
		userId := "test user"
		ctrl := gomock.NewController(t)
		saver := usecase.NewMockProteinEventSaver(ctrl)
		saver.EXPECT().UpdateTimeToDrink(gomock.Eq(ctx), gomock.Eq(userId)).Return(nil)

		h := &GotRequestHandler{
			saver: saver,
			params: &SlackCallbackRequestParams{
				UserId: userId,
			},
		}

		w := httptest.NewRecorder()
		h.Handler(context.TODO(), w)
		want, _ := json.Marshal(SlackCallbackResponse{Message: "success"})
		assert.Equal(t, w.Code, http.StatusOK)
		assert.Equal(t, w.Body.Bytes(), want)
	})
}
