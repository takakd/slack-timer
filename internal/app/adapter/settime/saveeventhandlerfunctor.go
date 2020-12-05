package settime

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slacktimer/internal/app/adapter/validator"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"slacktimer/internal/pkg/helper"
	"strconv"
	"strings"
	"time"
)

type SaveEventHandler interface {
	Handle(ctx context.Context, data EventCallbackData) *Response
}

// SetRequestHandler represents the API command "Set".
type SaveEventHandlerFunctor struct {
	// Time to notify user next
	notificationTime    time.Time
	remindIntervalInMin int
	inputPort           updatetimerevent.InputPort
}

func NewSaveEventHandlerFunctor() SaveEventHandler {
	return &SaveEventHandlerFunctor{
		inputPort: di.Get("slackcontroller.InputPort").(updatetimerevent.InputPort),
	}
}

// Validate parameters.
// TODO: naming parse? because set remindinterval value.
func (se *SaveEventHandlerFunctor) validate(data EventCallbackData) *validator.ValidateErrorBag {
	bag := validator.NewValidateErrorBag()

	// Extract second part. e.g.1607054661.000200 -> 160705466.
	s := strings.Split(data.MessageEvent.EventTs, ".")
	if len(s) < 1 {
		bag.SetError("timestamp", "invalid format", errors.New("invalid format"))
		return bag
	}

	eventTime, err := helper.ParseUnixStr(s[0])
	if err != nil {
		bag.SetError("timestamp", "invalid format", errors.New("invalid format"))
	}
	se.notificationTime = eventTime.UTC()

	// e.g. set 10
	re := regexp.MustCompile(`^(.*)\s+([0-9]+)$`)
	m := re.FindStringSubmatch(data.MessageEvent.Text)
	if m == nil {
		bag.SetError("interval", "invalid format", errors.New("invalid format"))
		return bag
	}
	minutes, _ := strconv.Atoi(m[2])
	se.remindIntervalInMin = minutes

	return bag
}

func (se *SaveEventHandlerFunctor) Handle(ctx context.Context, data EventCallbackData) *Response {
	if validateErrors := se.validate(data); len(validateErrors.GetErrors()) > 0 {
		var firstError *validator.ValidateError
		for _, v := range validateErrors.GetErrors() {
			firstError = v
			break
		}
		return makeErrorHandlerResponse("invalid parameter", firstError.Summary)
	}

	log.Info(fmt.Sprintf("Usecase.SaveIntervalMin user=%s notificationtime=%s interval=%d", data.MessageEvent.User, se.notificationTime, se.remindIntervalInMin))

	presenter := NewSaveEventOutputReceivePresenter()
	se.inputPort.SaveIntervalMin(ctx, data.MessageEvent.User, se.notificationTime, se.remindIntervalInMin, presenter)

	log.Info(fmt.Sprintf("Usecase.SaveIntervalMin output.resp=%v", *presenter.Resp))

	return presenter.Resp
}
