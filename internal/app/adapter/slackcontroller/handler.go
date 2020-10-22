// Package slackcontroller provides the slack Event API callback handler.
// 		Routes
//			POST /api/{ver}/slack-callback
// Library exists: https://github.com/slack-go/slack
// Ref.: https://api.slack.com/events-api#the-events-api__receiving-events
package slackcontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"proteinreminder/internal/app/driver/di"
	"proteinreminder/internal/app/usecase/updateproteinevent"
	"proteinreminder/internal/pkg/httputil"
	"proteinreminder/internal/pkg/log"
	"regexp"
)

// Errors
var (
	ErrInvalidRequest    = errors.New("invalid request")
	ErrInvalidParameters = errors.New("invalid parameters")
	ErrSaveEvent         = errors.New("failed to save protein event")
	ErrCreateResponse    = errors.New("failed to create response")
)

// Command types entered by users.
const (
	CmdGot = "got"
	CmdSet = "set"
)

// URL Verification Event callback data
// Ref. https://api.slack.com/events/url_verification
type UrlVerificationEventCallbackData struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
}

func (d *UrlVerificationEventCallbackData) doesMatchType() bool {
	return d.Type == "url_verification"
}

// JSON to be sent
// Defined what app needs
type EventCallbackData struct {
	Token  string `json:"token"`
	TeamId string `json:"team_id"`
	// Ref. https://api.slack.com/events
	MessageEvent MessageEvent `json:"event"`
	Type         string       `json:"type"`
	EventTime    int          `json:"event_time"`
}

type MessageEvent struct {
	Type    string `json:"type"`
	EventTs string `json:"event_ts"`
	User    string `json:"user"`
	Ts      string `json:"ts"`
	Text    string `json:"text"`
}

//
func NewRequestHandler(r *http.Request) (RequestHandler, error) {
	body, err := httputil.GetRequestBody(r)
	if err != nil {
		return nil, err
	}

	log.Info(body)

	// URL Verification callback
	urlVerification := UrlVerificationEventCallbackData{}
	err = json.Unmarshal(body, &urlVerification)
	if err != nil {
		return nil, err
	}
	if urlVerification.doesMatchType() {
		log.Info("url verification event")
		return &UrlVerificationRequestHandler{
			Data: &urlVerification,
		}, nil
	}

	// Normal Event callback
	data := EventCallbackData{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	var supportEvent bool
	supportEvent = supportEvent || data.MessageEvent.Type == "message"
	if !supportEvent {
		return nil, fmt.Errorf("invalid event type, type=%s", data.MessageEvent.Type)
	}

	// e.g. set 10, got
	re := regexp.MustCompile(`^([^\s]*)\s*`)
	m := re.FindStringSubmatch(data.MessageEvent.Text)
	if m == nil {
		return nil, fmt.Errorf("invalid Text format")
	}

	subType := m[1]
	if subType != CmdGot && subType != CmdSet {
		return nil, fmt.Errorf("invalid sub type")
	}

	usecase := di.Get("UpdateProteinEvent").(updateproteinevent.Usecase)

	var req RequestHandler
	if subType == CmdGot {
		log.Info(fmt.Sprintf("got event text=%s", data.MessageEvent.Text))
		req = &GotRequestHandler{
			messageEvent: &data.MessageEvent,
			usecase:      usecase,
		}
	} else if subType == CmdSet {
		log.Info(fmt.Sprintf("set event text=%s", data.MessageEvent.Text))
		req = &SetRequestHandler{
			messageEvent: &data.MessageEvent,
			usecase:      usecase,
		}
	}

	return req, nil
}

// Provides handlers to each request.
type RequestHandler interface {
	Handler(ctx context.Context, w http.ResponseWriter)
}

// Represents this API response.
type SlackCallbackResponse struct {
	Message string `json:"message"`
}

// Represents this API error response.
// Ref: https://developer.github.com/v3/
type ErrorSlackCallbackResponse struct {
	// Error brief.
	Message string `json:"message"`
	Error   string `json:"error"`
}

func makeErrorCallbackResponseBody(message string, err error) ([]byte, error) {
	resp := ErrorSlackCallbackResponse{
		Message: message,
		Error:   err.Error(),
	}
	body, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Web server registers this to themselves and call.
func Handler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	h, err := NewRequestHandler(r)
	if err != nil {
		log.Error(err.Error())
		body, err := makeErrorCallbackResponseBody("parameter error", ErrInvalidRequest)
		if err != nil {
			body = []byte("internal error")
		}
		httputil.WriteJsonResponse(w, http.StatusBadRequest, body)
		return
	}

	h.Handler(ctx, w)
}
