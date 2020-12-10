package settime

import (
	"errors"
	"regexp"
	"slacktimer/internal/app/adapter/validator"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/app/util/appcontext"
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
	Handle(ac appcontext.AppContext, data EventCallbackData) *Response
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

// validate and parse parameters.
func (se *SaveEventHandlerFunctor) parse(data EventCallbackData) *validator.ValidateErrorBag {
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
func (se SaveEventHandlerFunctor) Handle(ac appcontext.AppContext, data EventCallbackData) *Response {
	if validateErrors := se.parse(data); len(validateErrors.GetErrors()) > 0 {
		var firstError *validator.ValidateError
		for _, v := range validateErrors.GetErrors() {
			firstError = v
			break
		}
		return newErrorHandlerResponse(ac, "invalid parameter", firstError.Summary)
	}

	log.InfoWithContext(ac, "call inputport", map[string]interface{}{
		"user":              data.MessageEvent.User,
		"interval":          se.remindIntervalInMin,
		"notification time": se.notificationTime,
	})

	presenter := NewSaveEventOutputReceivePresenter()
	se.inputPort.SaveIntervalMin(ac, data.MessageEvent.User, se.notificationTime, se.remindIntervalInMin, presenter)

	log.InfoWithContext(ac, "return from inputport", presenter.Resp)

	return &presenter.Resp
}
