package slackcontroller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"proteinreminder/internal/app/adapter/validator"
	"proteinreminder/internal/app/usecase/updateproteinevent"
	"proteinreminder/internal/pkg/httputil"
	"proteinreminder/internal/pkg/log"
	"regexp"
	"strconv"
)

// SetRequestHandler represents the API command "Set".
type SetRequestHandler struct {
	messageEvent *MessageEvent
	// Time to notify user next
	remindIntervalInMin int
	usecase             updateproteinevent.Usecase
}

// Validate parameters.
func (sr *SetRequestHandler) validate() *validator.ValidateErrorBag {
	bag := validator.NewValidateErrorBag()

	// e.g. set 10
	re := regexp.MustCompile(`^(.*)\s+([0-9]+)$`)
	m := re.FindStringSubmatch(sr.messageEvent.Text)
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
		body, err := makeErrorCallbackResponseBody(firstError.Summary, ErrInvalidParameters)
		if err != nil {
			log.Error(err)
			httputil.WriteJsonResponse(w, http.StatusBadRequest, []byte("internal error"))
		}

		httputil.WriteJsonResponse(w, http.StatusBadRequest, body)
		return
	}

	outputPort := &SetRequestOutputPort{w: w}
	sr.usecase.SaveIntervalMin(ctx, sr.messageEvent.User, sr.remindIntervalInMin, outputPort)
	return
}

type SetRequestOutputPort struct {
	w http.ResponseWriter
}

func (s *SetRequestOutputPort) Output(data *updateproteinevent.OutputData) {
	err := data.Result
	errRaised := false
	if errors.Is(err, updateproteinevent.ErrFind) {
		errRaised = true
	} else if errors.Is(err, updateproteinevent.ErrCreate) {
		errRaised = true
	} else if errors.Is(err, updateproteinevent.ErrSave) {
		errRaised = true
	}

	resp := &SlackCallbackResponse{
		Message: "success",
	}
	respBody, err := json.Marshal(resp)
	if err != nil {
		errRaised = true
	}

	if errRaised {
		log.Error(err)
		body, err := makeErrorCallbackResponseBody("failed to save event", ErrSaveEvent)
		if err != nil {
			log.Error(err)
			body = []byte("internal error")
		}
		httputil.WriteJsonResponse(s.w, http.StatusBadRequest, body)
		return
	}

	httputil.WriteJsonResponse(s.w, http.StatusOK, respBody)
	return
}
