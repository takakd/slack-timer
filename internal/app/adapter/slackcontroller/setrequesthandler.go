package slackcontroller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slacktimer/internal/app/adapter/validator"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/app/util/log"
	"slacktimer/internal/pkg/timeutil"
	"strconv"
	"time"
)

// SetRequestHandler represents the API command "Set".
type SetRequestHandler struct {
	messageEvent *MessageEvent
	// Time to notify user next
	notificationTime    time.Time
	remindIntervalInMin int
	usecase             updatetimerevent.Usecase
}

// Validate parameters.
func (sr *SetRequestHandler) validate() *validator.ValidateErrorBag {
	bag := validator.NewValidateErrorBag()

	eventTime, err := timeutil.ParseUnixStr(sr.messageEvent.EventTs)
	if err != nil {
		bag.SetError("timestamp", "invalid format", errors.New("invalid format"))
	}
	sr.notificationTime = eventTime.UTC()

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
		return makeErrorHandlerResponse("invalid parameter", firstError.Summary)
	}

	outputPort := &SetRequestOutputPort{}

	log.Info(fmt.Sprintf("Usecase.SaveIntervalMin user=%s notificationtime=%s interval=%d", sr.messageEvent.User, sr.notificationTime, sr.remindIntervalInMin))

	sr.usecase.SaveIntervalMin(ctx, sr.messageEvent.User, sr.notificationTime, sr.remindIntervalInMin, outputPort)

	log.Info(fmt.Sprintf("Usecase.SaveIntervalMin output=%v", *outputPort))

	return outputPort.Resp
}

type SetRequestOutputPort struct {
	Resp *HandlerResponse
}

func (s *SetRequestOutputPort) Output(data *updatetimerevent.OutputData) {
	err := data.Result
	if err != nil {
		log.Info(fmt.Sprintf("SetRequestOutputPort.Output error=%v", err))
		s.Resp = makeErrorHandlerResponse("failed to set timer", "internal server error")
		return
	}

	s.Resp = &HandlerResponse{
		StatusCode: http.StatusOK,
		Body:       "success",
	}

	log.Info(fmt.Sprintf("SetRequestOutputPort.Output resp=%v", s.Resp))
}
