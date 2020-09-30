package slackcontroller

import (
	"context"
	"encoding/json"
	"fmt"
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
	"proteinreminder/internal/app/adapter"
)

//
// POST slack-callback
//
// Library exists: https://github.com/slack-go/slack
// Ref: https://api.slack.com/interactivity/slash-commands

// TODO: Change error definitions' type to Error.
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

type Input interface {
	Parse() (Request, error)
	Request() *http.Request
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

type SlackCallbackRequest struct {
	request *http.Request
	params  SlackCallbackRequestParams

	// The subtype of command is set after command.
	// e.g. /protein <sub type>
	// got: Mark the time when the protein was drunk.
	// set: Set the interval of minutes to drink.
	subType CommandSubType
}

func (c *SlackCallbackRequest) Request() *http.Request {
	return c.request
}

func (c *SlackCallbackRequest) Parse() (Request, error) {
	c.request.ParseForm()
	if err := httputil.SetFormValueToStruct(c.request.Form, &c.params); err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`([^\s]*)\s*`)
	m := re.FindStringSubmatch(c.params.Text)
	if m == nil {
		return nil, fmt.Errorf("invalid Text format")
	}

	c.subType = CommandSubType(m[1])
	if c.subType != SubTypeGot && c.subType != SubTypeSet {
		return nil, fmt.Errorf("invalid sub type")
	}

	var req Request
	if c.subType == SubTypeGot {
		req = &GotRequest{
			params: c.params,
		}
	} else if c.subType == SubTypeSet {
		req = &SetRequest{
			params: c.params,
		}
	}

	return req, nil
}

type Request interface {
	// TODO: リクエストごとの処理を実装する
	// 共通処理はprivateメソッドで切り出す
	Handler(ctx context.Context, w http.ResponseWriter)
}

type SetRequest struct {
	params SlackCallbackRequestParams
	// The time of entering a message on Slack.
	datetime            time.Time
	remindIntervalInMin time.Duration
	saver usecase.ProteinEventSaver
}

type GotRequest struct {
	params SlackCallbackRequestParams
	// The time of entering a message on Slack.
	datetime time.Time
	saver usecase.ProteinEventSaver
}

func (gr *GotRequest) validate() (error, *adapter.ValidateErrorBag) {
	// TODO
	bag := adapter.NewValidateErrorBag()
	return nil, bag
}

func (gr *GotRequest) Handler(ctx context.Context, w http.ResponseWriter) {
	// TODO
	return
}

func (sr *SetRequest) validate() (error, *adapter.ValidateErrorBag) {
	// TODO
	bag := adapter.NewValidateErrorBag()

	re := regexp.MustCompile(`(.*)\s+([0-9]+)`)
	m := re.FindStringSubmatch(sr.params.Text)
	if m == nil {
		return fmt.Errorf("invalid Text format"), bag
	}

	if minutes, err := strconv.Atoi(m[2]); err != nil {
		// the process doesn't come here.
		return err, bag
	} else {
		sr.remindIntervalInMin = time.Duration(minutes)
	}

	return nil, bag
}

func (sr *SetRequest) Handler(ctx context.Context, w http.ResponseWriter) {
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

//
func handler(ctx context.Context, saver usecase.ProteinEventSaver, w http.ResponseWriter, r Input) {
	if r.Request().Method != "POST" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	req, err := r.Parse()
	if err != nil {
		log.Error("%v", err.Error())
		httputil.WriteJsonResponse(w, http.StatusBadRequest, MakeErrorCallbackResponseBody("parameter error", SlackErrorCodeParse))
		return
	}

	req.Handler(ctx, w)
	return

	//if err, validateErrors := validator.Validate(); err != nil {
	//	var firstError *adapter.ValidateError
	//	for _, v := range validateErrors.GetErrors() {
	//		firstError = v
	//		break
	//	}
	//	httputil.WriteJsonResponse(w, http.StatusBadRequest, MakeErrorCallbackResponseBody(firstError.Summary, SlackErrorCodeVaidate))
	//	return
	//}
	//
	//return validator.Handle(ctx, w)

	//var errCode usecase.SaveProteinEventError
	//
	//switch req := validator.(type) {
	//case *GotRequest:
	//	errCode = saver.SaveTimeToDrink(ctx, req.params.UserId, req.datetime)
	//case *SetRequest:
	//	errCode = saver.SaveIntervalSec(ctx, req.params.UserId, req.remindIntervalInMin)
	//}
	//
	//if errCode != usecase.SaveProteinEventNoError {
	//	if errCode == usecase.SaveProteinEventErrorFind {
	//		httputil.WriteJsonResponse(w, http.StatusBadRequest, MakeErrorCallbackResponseBody("failed to find event", SlackErrorCodeSavingProteinEvent1))
	//	} else if errCode == usecase.SaveProteinEventErrorCreate {
	//		httputil.WriteJsonResponse(w, http.StatusBadRequest, MakeErrorCallbackResponseBody("failed to create event", SlackErrorCodeSavingProteinEvent1))
	//	} else if errCode == usecase.SaveProteinEventErrorSave {
	//		httputil.WriteJsonResponse(w, http.StatusBadRequest, MakeErrorCallbackResponseBody("failed to save event", SlackErrorCodeSavingProteinEvent1))
	//	}
	//	return
	//}
	//
	//// Make response.
	//resp := &SlackCallbackResponse{
	//	Message: "success",
	//}
	//respBody, err := json.Marshal(resp)
	//if err != nil {
	//	log.Error("%v", err.Error())
	//	httputil.WriteJsonResponse(w, http.StatusBadRequest, MakeErrorCallbackResponseBody("failed to create response", SlackErrorCodeCreateResponse))
	//}
	//httputil.WriteJsonResponse(w, http.StatusOK, respBody)
}

// POST handler.
func Handler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// TODO: 各リクエスト内の処理で実装
	//saver, err := usecase.NewSaveProteinEvent(apprule.NewPostgresRepository(config.GetConfig()))
	//if err != nil {
	//	httputil.WriteJsonResponse(w, http.StatusInternalServerError, MakeErrorCallbackResponseBody("internal", SlackErrorCodeSavingProteinEvent1))
	//	return
	//}

	input := &SlackCallbackRequest{}

	handler(ctx, saver, w, input)
}
