package slackcontroller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"proteinreminder/internal/app/usecase"
	"proteinreminder/internal/pkg/httputil"
	"time"
)

// GotRequestHandler represents the API command "Got".
type GotRequestHandler struct {
	params *SlackCallbackRequestParams
	// User entered time on Slack
	datetime time.Time
	// Usecase to save entity
	saver usecase.ProteinEventSaver
}

func (gr *GotRequestHandler) Handler(ctx context.Context, w http.ResponseWriter) {
	// Save protein event.
	err := gr.saver.SaveTimeToDrink(ctx, gr.params.UserId, gr.datetime)
	if errors.Is(err, usecase.ErrFind) {
		httputil.WriteJsonResponse(w, http.StatusBadRequest, makeErrorCallbackResponseBody("failed to find event", ErrSaveEvent))
		return
	} else if errors.Is(err, usecase.ErrCreate) {
		httputil.WriteJsonResponse(w, http.StatusBadRequest, makeErrorCallbackResponseBody("failed to create event", ErrSaveEvent))
		return
	} else if errors.Is(err, usecase.ErrSave) {
		httputil.WriteJsonResponse(w, http.StatusBadRequest, makeErrorCallbackResponseBody("failed to save event", ErrSaveEvent))
		return
	}

	resp := &SlackCallbackResponse{
		Message: "success",
	}
	respBody, _ := json.Marshal(resp)
	httputil.WriteJsonResponse(w, http.StatusOK, respBody)
	return
}
