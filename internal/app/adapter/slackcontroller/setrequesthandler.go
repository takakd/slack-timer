package slackcontroller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"proteinreminder/internal/app/adapter"
	"proteinreminder/internal/app/usecase"
	"proteinreminder/internal/pkg/httputil"
	"proteinreminder/internal/pkg/log"
	"regexp"
	"strconv"
	"time"
)

// SetRequestHandler represents the API command "Set".
type SetRequestHandler struct {
	params *SlackCallbackRequestParams
	// The time entered by the user.
	datetime time.Time
	// The time the user is notified next time.
	remindIntervalInMin time.Duration
	saver               usecase.ProteinEventSaver
}

//
func (sr *SetRequestHandler) validate() (*adapter.ValidateErrorBag, error) {
	bag := adapter.NewValidateErrorBag()

	re := regexp.MustCompile(`(.*)\s+([0-9]+)`)
	m := re.FindStringSubmatch(sr.params.Text)
	if m == nil {
		return nil, fmt.Errorf("invalid Text format")
	}

	if minutes, err := strconv.Atoi(m[2]); err != nil {
		// the process doesn't come here.
		return bag, err
	} else {
		sr.remindIntervalInMin = time.Duration(minutes)
	}

	return bag, nil
}

//
func (sr *SetRequestHandler) Handler(ctx context.Context, w http.ResponseWriter) {
	if validateErrors, err := sr.validate(); err != nil {
		var firstError *adapter.ValidateError
		for _, v := range validateErrors.GetErrors() {
			firstError = v
			break
		}
		httputil.WriteJsonResponse(w, http.StatusBadRequest, makeErrorCallbackResponseBody(firstError.Summary, ErrInvalidParameters))
		return
	}

	err := sr.saver.SaveIntervalSec(ctx, sr.params.UserId, sr.remindIntervalInMin)
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
	respBody, err := json.Marshal(resp)
	if err != nil {
		log.Error("%v", err.Error())
		httputil.WriteJsonResponse(w, http.StatusBadRequest, makeErrorCallbackResponseBody("failed to create response", ErrCreateResponse))
	}
	httputil.WriteJsonResponse(w, http.StatusOK, respBody)
}
