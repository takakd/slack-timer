package settime

import (
	"context"
	"errors"
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

var (
	// ErrInvalidFormat returns if parameter is invalid format.
	ErrInvalidFormat = errors.New("invalid format")
)

// SaveEventHandler handles "set" command.
type SaveEventHandler interface {
	Handle(ctx context.Context, data EventCallbackData) *Response
}

// SaveEventHandlerFunctor handle "Set" command.
type SaveEventHandlerFunctor struct {
	notificationTime    time.Time
	remindIntervalInMin int
	inputPort           updatetimerevent.InputPort
}

var _ SaveEventHandler = (*SaveEventHandlerFunctor)(nil)

// NewSaveEventHandlerFunctor create new struct.
func NewSaveEventHandlerFunctor() *SaveEventHandlerFunctor {
	return &SaveEventHandlerFunctor{
		inputPort: di.Get("updatetimerevent.InputPort").(updatetimerevent.InputPort),
	}
}

// Validate parameters.
// TODO: naming parse? because set remindinterval value.
func (se *SaveEventHandlerFunctor) validate(data EventCallbackData) *validator.ValidateErrorBag {
	bag := validator.NewValidateErrorBag()

	// Extract second part. e.g.1607054661.000200 -> 160705466.
	s := strings.Split(data.MessageEvent.EventTs, ".")
	if len(s) < 1 {
		bag.SetError("timestamp", "invalid format", ErrInvalidFormat)
		return bag
	}

	eventTime, err := helper.ParseUnixStr(s[0])
	if err != nil {
		bag.SetError("timestamp", "invalid format", ErrInvalidFormat)
	}
	se.notificationTime = eventTime.UTC()

	// e.g. set 10
	re := regexp.MustCompile(`^(.*)\s+([0-9]+)$`)
	m := re.FindStringSubmatch(data.MessageEvent.Text)
	if m == nil {
		bag.SetError("interval", "invalid format", ErrInvalidFormat)
		return bag
	}
	minutes, _ := strconv.Atoi(m[2])
	se.remindIntervalInMin = minutes

	return bag
}

// Handle saves event sent by user.
func (se SaveEventHandlerFunctor) Handle(ctx context.Context, data EventCallbackData) *Response {
	if validateErrors := se.validate(data); len(validateErrors.GetErrors()) > 0 {
		var firstError *validator.ValidateError
		for _, v := range validateErrors.GetErrors() {
			firstError = v
			break
		}
		return newErrorHandlerResponse("invalid parameter", firstError.Summary)
	}

	log.Info("call inputport", map[string]interface{}{
		"user":              data.MessageEvent.User,
		"interval":          se.remindIntervalInMin,
		"notification time": se.notificationTime,
	})

	presenter := NewSaveEventOutputReceivePresenter()
	se.inputPort.SaveIntervalMin(ctx, data.MessageEvent.User, se.notificationTime, se.remindIntervalInMin, presenter)

	log.Info("return from inputport", presenter.Resp)

	return &presenter.Resp
}
