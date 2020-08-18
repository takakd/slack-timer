package controller

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"proteinreminder/internal/httputil"
	"proteinreminder/internal/ioc"
	"proteinreminder/internal/model"
	"proteinreminder/internal/panicutil"
	"regexp"
	"strconv"
	"time"
	"context"
)

//
// POST slack-callback
//
// Library exists: https://github.com/slack-go/slack

const (
	SlackErrorCodeNoError= 0
	SlackErrorCodeParse = 1
	SlackErrorCodeVaidate = 2
	SlackErrorCodeSavingProteinEvent1 = 3
	SlackErrorCodeSavingProteinEvent2 = 4
	SlackErrorCodeCreateResponse = 5
)

type SlackCallbackRequest struct {
	request *http.Request
	params  SlackCallbackRequestParams
	// The word entered on Slack.
	keyword string
	// The time of entering a message on Slack.
	datetime time.Time
}

type SlackCallbackResponse struct {
	Message string `json:"message"`
}

// Ref: https://developer.github.com/v3/
type ErrorSlackCallbackResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

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

func NewSlackCallbackRequest(r *http.Request) *SlackCallbackRequest {
	return &SlackCallbackRequest{
		request: r,
		params:  SlackCallbackRequestParams{},
	}
}

func (r *SlackCallbackRequest) parse() error {
	r.request.ParseForm()
	if err := httputil.SetFormValueToStruct(r.request.Form, &r.params); err != nil {
		return err
	}

	re := regexp.MustCompile("(.*)\\s+([0-9]+):([0-9]+)")
	m := re.FindStringSubmatch(r.params.Text)
	if m == nil {
		return errors.New("invalid Text format.")
	}

	r.keyword = m[1]

	hour, err := strconv.Atoi(m[2])
	if err != nil {
		return err
	}
	minute, err := strconv.Atoi(m[3])
	if err != nil {
		return err
	}
	t := time.Now()
	r.datetime = time.Date(t.Year(), t.Month(), t.Day(), hour, minute, 0, 0, time.UTC)

	return err
}

func (r *SlackCallbackRequest) validate() (bool, *ValidateErrorBag) {
	valid := true
	bag := NewValidateErrorBag()
	if r.keyword == "" {
		valid = false
		bag.SetError("keyword", "need keyword.", Empty)
	}
	if r.params.UserId == "" {
		valid = false
		bag.SetError("user_id", "need user_id.", Empty)
	}
	return valid, bag
}

func MakeErrorCallbackResponseBody(message string, code int) []byte {
	resp := ErrorSlackCallbackResponse{
		Message: message,
		Code:    code,
	}
	body, err := json.Marshal(resp)
	if err != nil {
		panic(panicutil.MakePanicMessage(err))
	}
	return body
}

// POST handler.
func SlackCallbackHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	logger := ioc.GetLogger()

	req := NewSlackCallbackRequest(r)
	if err := req.parse(); err != nil {
		logger.Error("%v", err.Error())
		httputil.WriteJsonResponse(w, http.StatusBadRequest, MakeErrorCallbackResponseBody("parameter error", SlackErrorCodeParse))
		return
	}

	if ok, validateErrors := req.validate(); !ok {
		var firstError *ValidateError
		for _, v := range validateErrors.errors {
			firstError = v
			break
		}
		httputil.WriteJsonResponse(w, http.StatusBadRequest, MakeErrorCallbackResponseBody(firstError.Summary, SlackErrorCodeVaidate))
		return
	}

	// Save event.
	event, err := model.NewProteinEvent(req.params.UserId)
	if err != nil {
		logger.Error("%v", err.Error())
		httputil.WriteJsonResponse(w, http.StatusBadRequest, MakeErrorCallbackResponseBody("failed to create event", SlackErrorCodeSavingProteinEvent1))
		return
	}
	if err := model.SaveProteinEvent(event); err != nil {
		httputil.WriteJsonResponse(w, http.StatusBadRequest, MakeErrorCallbackResponseBody("failed to save event", SlackErrorCodeSavingProteinEvent2))
		return
	}

	resp := &SlackCallbackResponse{
		Message: "success",
	}
	respBody, err := json.Marshal(resp)
	if err != nil {
		logger.Error("%v", err.Error())
		httputil.WriteJsonResponse(w, http.StatusBadRequest, MakeErrorCallbackResponseBody("failed to create response", SlackErrorCodeCreateResponse))
		return
	}

	httputil.WriteJsonResponse(w, http.StatusOK, respBody)
}
