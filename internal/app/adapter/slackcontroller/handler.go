// Package slackcontroller provides the slack callback handler.
// 		Routes
//			POST /slack-callback
// Library exists: https://github.com/slack-go/slack
// Ref: https://api.slack.com/interactivity/slash-commands
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

// Errors.
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

// Parameters in Slack webhook post body.
// Ref: https://api.slack.com/interactivity/slash-commands
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
	params := &SlackCallbackRequestParams{}
	r.ParseForm()
	if err := httputil.SetFormValueToStruct(r.Form, params); err != nil {
		return nil, err
	}

	// e.g. set 10, got
	re := regexp.MustCompile(`^([^\s]*)\s*`)
	m := re.FindStringSubmatch(params.Text)
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
		req = &GotRequestHandler{
			params:  params,
			usecase: usecase,
		}
	} else if subType == CmdSet {
		req = &SetRequestHandler{
			params:  params,
			usecase: usecase,
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
