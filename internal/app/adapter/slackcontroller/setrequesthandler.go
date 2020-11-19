package slackcontroller

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"slacktimer/internal/app/adapter/validator"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/pkg/log"
	"strconv"
	"time"
)

// SetRequestHandler represents the API command "Set".
type SetRequestHandler struct {
	messageEvent *MessageEvent
	// Time to notify user next
	remindIntervalInMin int
	usecase             updatetimerevent.Usecase
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

func (sr *SetRequestHandler) Handler(ctx context.Context) *HandlerResponse {
	if validateErrors := sr.validate(); len(validateErrors.GetErrors()) > 0 {
		var firstError *validator.ValidateError
		for _, v := range validateErrors.GetErrors() {
			firstError = v
			break
		}
		return makeErrorHandlerResponse(firstError.Summary, ErrInvalidParameters)
	}

	outputPort := &SetRequestOutputPort{}
	now := time.Now().UTC()
	sr.usecase.SaveIntervalMin(ctx, sr.messageEvent.User, now, sr.remindIntervalInMin, outputPort)
	log.Debug(outputPort)
	return outputPort.Resp
}

type SetRequestOutputPort struct {
	Resp *HandlerResponse
}

func (s *SetRequestOutputPort) Output(data *updatetimerevent.OutputData) {
	err := data.Result
	errRaised := false
	if errors.Is(err, updatetimerevent.ErrFind) {
		errRaised = true
	} else if errors.Is(err, updatetimerevent.ErrCreate) {
		errRaised = true
	} else if errors.Is(err, updatetimerevent.ErrSave) {
		errRaised = true
	}

	if errRaised {
		log.Error(err)
		s.Resp = makeErrorHandlerResponse("failed to save event", ErrSaveEvent)
		return
	}

	s.Resp = &HandlerResponse{
		StatusCode: http.StatusOK,
		Body:       "success",
	}
}
