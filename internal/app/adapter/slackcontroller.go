package adapter

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"proteinreminder/internal/app/apprule"
	"proteinreminder/internal/app/usecase"
	"proteinreminder/internal/pkg/config"
	"proteinreminder/internal/pkg/errorutil"
	"proteinreminder/internal/pkg/httputil"
	"proteinreminder/internal/pkg/log"
	"regexp"
	"strconv"
	"time"
)

//
// POST slack-callback
//
// Library exists: https://github.com/slack-go/slack
// Ref: https://api.slack.com/interactivity/slash-commands

const (
	SlackErrorCodeNoError             = 0
	SlackErrorCodeParse               = 1
	SlackErrorCodeVaidate             = 2
	SlackErrorCodeSavingProteinEvent1 = 3
	SlackErrorCodeSavingProteinEvent2 = 4
	SlackErrorCodeCreateResponse      = 5
	SlackErrorCodeInvalidSubtype      = 6

	SubTypeGot CommandSubType = "got"
	SubTypeSet CommandSubType = "set"
)

type CommandSubType string

type SlackCallbackRequest struct {
	request *http.Request
	params  SlackCallbackRequestParams

	// The subtype of command is set after command.
	// e.g. /protein <sub type>
	// got: Mark the time when the protein was drunk.
	// set: Set the interval of minutes to drink.
	subType CommandSubType

	// The time of entering a message on Slack.
	datetime time.Time
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

type Validator interface {
	validate() (bool, *ValidateErrorBag)
}

func parseRequest(r *http.Request) (Validator, error) {
	request := &SlackCallbackRequest{
		request: r,
		params:  SlackCallbackRequestParams{},
	}

	r.ParseForm()
	if err := httputil.SetFormValueToStruct(r.Form, &request.params); err != nil {
		return nil, err
	}

	re := regexp.MustCompile("(.*)\\s*")
	m := re.FindStringSubmatch(request.params.Text)
	if m == nil {
		return nil, errors.New("invalid Text format")
	}

	request.subType = CommandSubType(m[1])

	var validator Validator
	if request.subType == SubTypeGot {
		validator = MakeSlackCallbackGotRequest(request)
	} else if request.subType == SubTypeGot {
		var err error
		validator, err = MakeSlackCallbackSetRequest(request)
		if err != nil {
			return nil, err
		}
	}

	//hour, err := strconv.Atoi(m[2])
	//if err != nil {
	//	return err
	//}
	//minute, err := strconv.Atoi(m[3])
	//if err != nil {
	//	return err
	//}
	//t := time.Now()
	//r.datetime = time.Date(t.Year(), t.Month(), t.Day(), hour, minute, 0, 0, time.UTC)

	return validator, nil
}

func (r *SlackCallbackRequest) validate() (bool, *ValidateErrorBag) {
	valid := true
	bag := NewValidateErrorBag()
	if r.params.UserId == "" {
		valid = false
		bag.SetError("user_id", "need user_id.", Empty)
	}
	return valid, bag
}

type SlackCallbackGotRequest struct {
	SlackCallbackRequest
}

func MakeSlackCallbackGotRequest(r *SlackCallbackRequest) *SlackCallbackGotRequest {
	return &SlackCallbackGotRequest{
		*r,
	}
}

func (r *SlackCallbackGotRequest) validate() (bool, *ValidateErrorBag) {
	valid, bag := r.SlackCallbackRequest.validate()
	if !valid {
		return valid, bag
	}
	// TOOD: check datetime
	return true, nil
}

type SlackCallbackSetRequest struct {
	SlackCallbackRequest

	remindIntervalInMin time.Duration
}

func MakeSlackCallbackSetRequest(r *SlackCallbackRequest) (*SlackCallbackSetRequest, error) {
	req := &SlackCallbackSetRequest{
		SlackCallbackRequest: *r,
	}

	re := regexp.MustCompile("(.*)\\s+([0-9]+)")
	m := re.FindStringSubmatch(r.params.Text)
	if m == nil {
		return nil, errors.New("invalid Text format")
	}

	if minutes, err := strconv.Atoi(m[2]); err != nil {
		// the process doesn't come here.
		return nil, err
	} else {
		req.remindIntervalInMin = time.Duration(minutes)
	}

	return req, nil
}

func (r *SlackCallbackSetRequest) validate() (bool, *ValidateErrorBag) {
	valid, bag := r.SlackCallbackRequest.validate()
	if !valid {
		return valid, bag
	}
	// TODO: check duration
	return true, nil
}

type SlackCallbackResponse struct {
	Message string `json:"message"`
}

// Ref: https://developer.github.com/v3/
type ErrorSlackCallbackResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func MakeErrorCallbackResponseBody(message string, code int) []byte {
	resp := ErrorSlackCallbackResponse{
		Message: message,
		Code:    code,
	}
	body, err := json.Marshal(resp)
	if err != nil {
		panic(errorutil.MakePanicMessage(err))
	}
	return body
}

// POST handler.
func SlackCallbackHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	validator, err := parseRequest(r)
	if err != nil {
		log.Error("%v", err.Error())
		httputil.WriteJsonResponse(w, http.StatusBadRequest, MakeErrorCallbackResponseBody("parameter error", SlackErrorCodeParse))
		return
	}

	if ok, validateErrors := validator.validate(); !ok {
		var firstError *ValidateError
		for _, v := range validateErrors.errors {
			firstError = v
			break
		}
		httputil.WriteJsonResponse(w, http.StatusBadRequest, MakeErrorCallbackResponseBody(firstError.Summary, SlackErrorCodeVaidate))
		return
	}

	saveProteinEvent, err := usecase.NewSaveProteinEvent(apprule.NewMongoDbRepository(config.GetConfig("")))
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}
	var errCode usecase.SaveProteinEventError

	switch req := validator.(type) {
	case *SlackCallbackGotRequest:
		errCode = saveProteinEvent.SaveTimeToDrink(ctx, req.params.UserId, req.datetime)
	case *SlackCallbackSetRequest:
		errCode = saveProteinEvent.SaveIntervalSec(ctx, req.params.UserId, req.remindIntervalInMin)
	}

	if errCode != usecase.SaveProteinEventNoError {
		if errCode == usecase.SaveProteinEventErrorFind {
			httputil.WriteJsonResponse(w, http.StatusBadRequest, MakeErrorCallbackResponseBody("failed to find event", SlackErrorCodeSavingProteinEvent1))
		} else if errCode == usecase.SaveProteinEventErrorCreate {
			httputil.WriteJsonResponse(w, http.StatusBadRequest, MakeErrorCallbackResponseBody("failed to create event", SlackErrorCodeSavingProteinEvent1))
		} else if errCode == usecase.SaveProteinEventErrorSave {
			httputil.WriteJsonResponse(w, http.StatusBadRequest, MakeErrorCallbackResponseBody("failed to save event", SlackErrorCodeSavingProteinEvent1))
		}
		return
	}

	// Make response.
	resp := &SlackCallbackResponse{
		Message: "success",
	}
	respBody, err := json.Marshal(resp)
	if err != nil {
		log.Error("%v", err.Error())
		httputil.WriteJsonResponse(w, http.StatusBadRequest, MakeErrorCallbackResponseBody("failed to create response", SlackErrorCodeCreateResponse))
	}
	httputil.WriteJsonResponse(w, http.StatusOK, respBody)
}
