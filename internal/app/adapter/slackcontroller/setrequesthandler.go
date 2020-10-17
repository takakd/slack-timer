package slackcontroller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"proteinreminder/internal/app/adapter/validator"
	"proteinreminder/internal/app/usecase"
	"proteinreminder/internal/pkg/httputil"
	"regexp"
	"strconv"
)

// SetRequestHandler represents the API command "Set".
type SetRequestHandler struct {
	params *SlackCallbackRequestParams
	// Time to notify user next
	remindIntervalInMin int
	// Usecase to save entity
	saver usecase.ProteinEventSaver
}

// Validate parameters.
func (sr *SetRequestHandler) validate() *validator.ValidateErrorBag {
	bag := validator.NewValidateErrorBag()

	// e.g. set 10
	re := regexp.MustCompile(`(.*)\s+([0-9]+)`)
	m := re.FindStringSubmatch(sr.params.Text)
	if m == nil {
		bag.SetError("interval", "invalid format", errors.New("invalid format"))
		return bag
	}

	minutes, _ := strconv.Atoi(m[2])
	sr.remindIntervalInMin = minutes

	return bag
}

func (sr *SetRequestHandler) Handler(ctx context.Context, w http.ResponseWriter) {
	if validateErrors := sr.validate(); len(validateErrors.GetErrors()) > 0 {
		var firstError *validator.ValidateError
		for _, v := range validateErrors.GetErrors() {
			firstError = v
			break
		}
		httputil.WriteJsonResponse(w, http.StatusBadRequest, makeErrorCallbackResponseBody(firstError.Summary, ErrInvalidParameters))
		return
	}

	// Save protein event.
	err := sr.saver.SaveIntervalMin(ctx, sr.params.UserId, sr.remindIntervalInMin)
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
