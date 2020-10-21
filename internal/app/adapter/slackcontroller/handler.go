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
	"proteinreminder/internal/pkg/httputil"
	"proteinreminder/internal/pkg/log"
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

// URL verification callback data
type SlackUrlVerificationCallbackData struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
}

// JSON to be sent
// Defined what app needs
type SlackCallbackData struct {
	Token     string             `json:"token"`
	TeamId    string             `json:"team_id"`
	Event     SlackCallbackEvent `json:"event"`
	Type      string             `json:"type"`
	EventTime int                `json:"event_time"`
}

type SlackCallbackEvent struct {
	Type    string `json:"type"`
	EventTs string `json:"event_ts"`
	User    string `json:"user"`
	Ts      string `json:"ts"`
	Item    string `json:"item"`
}

type SlackCallbackRequestParams struct {
	Token          string `json:"token"`
	TeamId         string `json:"team_id"`
	TeamDomain     string `json:"team_domain"`
	EnterpriseId   string `json:"enterprise_id"`
	EnterpriseName string `json:"enterprise_name"`
	ChannelId      string `json:"channel_id"`
	ChannelName    string `json:"channel_name"`
	UserId         string `json:"user_id"`
	UserName       string `json:"user_name"`
	Command        string `json:"command"`
	Text           string `json:"text"`
	ResponseUrl    string `json:"response_url"`
	TriggerId      string `json:"trigger_id"`
}

//
func NewRequestHandler(r *http.Request) (RequestHandler, error) {
	body, err := httputil.GetRequestBody(r)
	if err != nil {
		return nil, err
	}

	data := SlackCallbackData{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	var supportEvent bool
	supportEvent = supportEvent || data.Event.Type == "message"
	if !supportEvent {
		return nil, fmt.Errorf("invalid event type, type=%s", data.Event.Type)
	}

	log.Debug("data")
	log.Debug(data)

	return nil, nil

	//// e.g. set 10, got
	//re := regexp.MustCompile(`^([^\s]*)\s*`)
	//m := re.FindStringSubmatch(params.Text)
	//if m == nil {
	//	return nil, fmt.Errorf("invalid Text format")
	//}
	//
	//subType := m[1]
	//if subType != CmdGot && subType != CmdSet {
	//	return nil, fmt.Errorf("invalid sub type")
	//}
	//
	//usecase := di.Get("UpdateProteinEvent").(updateproteinevent.Usecase)
	//
	//var req RequestHandler
	//if subType == CmdGot {
	//	req = &GotRequestHandler{
	//		params:  params,
	//		usecase: usecase,
	//	}
	//} else if subType == CmdSet {
	//	req = &SetRequestHandler{
	//		params:  params,
	//		usecase: usecase,
	//	}
	//}
	//
	//return req, nil
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

	// TODO: Test
	body, err := httputil.GetRequestBody(r)
	if err != nil {
		makeErrorCallbackResponseBody("get body error", err)
	}

	verify := SlackUrlVerificationCallbackData{}
	err = json.Unmarshal(body, &verify)
	if err != nil {
		makeErrorCallbackResponseBody("unmarshal error", err)
	}
	if verify.Challenge != "" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(verify.Challenge))
		return
	}

	data := SlackCallbackData{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		makeErrorCallbackResponseBody("unmarshal error", err)
	}

	log.Debug("data")
	log.Debug(data)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(body)

	//h, err := NewRequestHandler(r)
	//if err != nil {
	//	log.Error(err.Error())
	//	body, err := makeErrorCallbackResponseBody("parameter error", ErrInvalidRequest)
	//	if err != nil {
	//		body = []byte("internal error")
	//	}
	//	httputil.WriteJsonResponse(w, http.StatusBadRequest, body)
	//	return
	//}
	//
	//h.Handler(ctx, w)
}
