// Package slackcontroller provides the slack Event API callback handler.
// Ref: https://api.slack.com/events-api#the-events-api__receiving-events
package slackcontroller

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"regexp"
	"slacktimer/internal/app/driver/di"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/pkg/log"
)

// Errors
var (
	ErrInvalidRequest    = errors.New("invalid request")
	ErrInvalidParameters = errors.New("invalid parameters")
	ErrSaveEvent         = errors.New("failed to save timer event")
	ErrCreateResponse    = errors.New("failed to create response")
)

// Command types entered by users.
const (
	CmdSet = "set"
)

// JSON to be sent
// Defined what app needs
type EventCallbackData struct {
	Token  string `json:"token"`
	TeamId string `json:"team_id"`
	// Ref. https://api.slack.com/events
	MessageEvent MessageEvent `json:"event"`
	Type         string       `json:"type"`
	EventTime    int          `json:"event_time"`

	// This field is only included in URL Verification Event.
	// Ref: https://api.slack.com/events/url_verification
	Challenge string `json:"challenge"`
}

func (e *EventCallbackData) isVerificationEvent() bool {
	return e.Type == "url_verification"
}

type MessageEvent struct {
	Type    string `json:"type"`
	EventTs string `json:"event_ts"`
	User    string `json:"user"`
	Ts      string `json:"ts"`
	Text    string `json:"text"`
}

type EventCallbackResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status"`
	Detail     string `json:"detail"`
}

func NewRequestHandler(data *EventCallbackData) (RequestHandler, error) {
	// URL Verification callback
	if data.isVerificationEvent() {
		log.Info("url verification event")
		return &UrlVerificationRequestHandler{
			Data: data,
		}, nil
	}

	// Normal Event callback
	supportEvent := data.MessageEvent.Type == "message"
	if !supportEvent {
		return nil, fmt.Errorf("invalid event type, type=%s", data.MessageEvent.Type)
	}

	// e.g. set 10
	re := regexp.MustCompile(`^([^\s]*)\s*`)
	m := re.FindStringSubmatch(data.MessageEvent.Text)
	if m == nil {
		return nil, fmt.Errorf("invalid Text format, text=%s", data.MessageEvent.Text)
	}

	subType := m[1]
	if subType != CmdSet {
		return nil, fmt.Errorf("invalid sub type, subtype=%s", subType)
	}

	usecase := di.Get("UpdateTimerEvent").(updatetimerevent.Usecase)

	log.Info(fmt.Sprintf("set event text=%s", data.MessageEvent.Text))
	req := &SetRequestHandler{
		messageEvent: &data.MessageEvent,
		usecase:      usecase,
	}

	return req, nil
}

// Provides handlers to each request.
type RequestHandler interface {
	Handler(ctx context.Context) EventCallbackResponse
}

func makeErrorCallbackResponse(message string, err error) *EventCallbackResponse {
	resp := &EventCallbackResponse{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
	if err != nil {
		resp.Detail = err.Error()
	}
	return resp
}

// Lambda callback
func LambdaHandleRequest(ctx context.Context, event EventCallbackData) (EventCallbackResponse, error) {
	h, err := NewRequestHandler(&event)
	if err != nil {
		log.Error(err.Error())
		return *makeErrorCallbackResponse("parameter error", ErrInvalidRequest), nil
	}

	return h.Handler(ctx), nil
}
