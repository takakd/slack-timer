package slackcontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"proteinreminder/internal/app/usecase"
	"testing"
	"time"
)

func TestSetRequestHandler_validate(t *testing.T) {
	cases := []struct {
		name  string
		text  string
		min   time.Duration
		valid bool
	}{
		{"OK", "set 10", 10, true},
		{"OK", "set 1", 1, true},
		{"NG", "set -1", 0, false},
		{"NG", "set", 0, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := SetRequestHandler{
				params: &SlackCallbackRequestParams{
					UserId: "test",
					Text:   c.text,
				},
				datetime: time.Now(),
			}
			bag := r.validate()
			_, exists := bag.GetError("interval")
			assert.Equal(t, c.valid, !exists)
			if c.valid {
				assert.Equal(t, time.Duration(c.min)*time.Minute, r.remindIntervalInMin)
			}
		})
	}
}

func TestSetRequestHandler_Handler(t *testing.T) {
	// Validate error
	t.Run("validate error", func(t *testing.T) {
		r := SetRequestHandler{
			params: &SlackCallbackRequestParams{
				UserId: "test",
				Text:   "",
			},
			datetime: time.Now(),
		}

		w := httptest.NewRecorder()
		r.Handler(context.TODO(), w)
		want, _ := json.Marshal(ErrorSlackCallbackResponse{
			Message: "invalid format",
			Error:   ErrInvalidParameters,
		})
		assert.Equal(t, w.Body.Bytes(), want)
	})

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
			interval := time.Duration(120) * time.Minute
			ctrl := gomock.NewController(t)
			saver := usecase.NewMockProteinEventSaver(ctrl)
			saver.EXPECT().SaveIntervalSec(gomock.Eq(ctx), gomock.Eq(userId), gomock.Eq(interval)).Return(c.err)

			h := &SetRequestHandler{
				saver: saver,
				params: &SlackCallbackRequestParams{
					UserId: userId,
					Text:   fmt.Sprintf("set %v", interval.Minutes()),
				},
			}

			w := httptest.NewRecorder()
			h.Handler(context.TODO(), w)
			assert.Equal(t, w.Code, http.StatusBadRequest)
			assert.Equal(t, w.Body.Bytes(), makeErrorCallbackResponseBody(c.msg, c.err))
		})
	}

	// Success
	t.Run("success", func(t *testing.T) {
		ctx := context.TODO()
		userId := "test user"
		interval := time.Duration(120) * time.Minute
		ctrl := gomock.NewController(t)
		saver := usecase.NewMockProteinEventSaver(ctrl)
		saver.EXPECT().SaveIntervalSec(gomock.Eq(ctx), gomock.Eq(userId), gomock.Eq(interval)).Return(nil)

		h := &SetRequestHandler{
			saver: saver,
			params: &SlackCallbackRequestParams{
				UserId: userId,
				Text:   fmt.Sprintf("set %v", interval.Minutes()),
			},
		}

		w := httptest.NewRecorder()
		h.Handler(context.TODO(), w)
		want, _ := json.Marshal(SlackCallbackResponse{Message: "success"})
		assert.Equal(t, w.Code, http.StatusOK)
		assert.Equal(t, w.Body.Bytes(), want)
	})
}
