// Package slackcontroller provides the slack Event API callback handler.
// Ref: https://api.slack.com/events-api#the-events-api__receiving-events
package slackcontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slacktimer/internal/app/driver/di"
	"slacktimer/internal/app/driver/di/container"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/pkg/config"
	"slacktimer/internal/pkg/config/driver"
	"slacktimer/internal/pkg/errorutil"
	"slacktimer/internal/pkg/fileutil"
	"slacktimer/internal/pkg/log"
	"slacktimer/internal/pkg/typeutil"
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

// Lambda handler input data
type LambdaInput struct {
	Resource                        string                `json:"resource,omitempty"`
	Path                            string                `json:"path,omitempty"`
	HttpMethod                      string                `json:"httpMethod,omitempty"`
	Headers                         map[string]string     `json:"headers,omitempty"`
	MultiValueHeaders               map[string][]string   `json:"multiValueHeaders,omitempty"`
	QueryStringParameters           map[string]string     `json:"queryStringParameters,omitempty"`
	MultiValueQueryStringParameters []map[string][]string `json:"multiValueQueryStringParameters,omitempty"`
	PathParameters                  map[string]string     `json:"pathParameters,omitempty"`
	StageVaribales                  map[string]string     `json:"stageVariables,omitempty"`
	RequestContext                  struct {
		AccountId  string `json:"accountId,omitempty"`
		ResourceId string `json:"resourceId,omitempty"`
		Stage      string `json:"stage,omitempty"`
		RequestId  string `json:"requestId,omitempty"`
		Identity   struct {
			CognitoIdentityPoolId         string `json:"cognitoIdentityPoolId,omitempty"`
			AccountId                     string `json:"accountId,omitempty"`
			CognitoIdentityId             string `json:"cognitoIdentityId,omitempty"`
			Caller                        string `json:"caller,omitempty"`
			ApiKey                        string `json:"apiKey,omitempty"`
			SourceIp                      string `json:"sourceIp,omitempty"`
			CognitoAuthenticationType     string `json:"cognitoAuthenticationType,omitempty"`
			CognitoAuthenticationProvider string `json:"cognitoAuthenticationProvider,omitempty"`
			UserArn                       string `json:"userArn,omitempty"`
			UserAgent                     string `json:"userAgent,omitempty"`
			User                          string `json:"user,omitempty"`
		} `json:"identity,omitempty"`
		ResourcePath string `json:"resourcePath,omitempty"`
		HttpMethod   string `json:"httpMethod,omitempty"`
		ApiId        string `json:"apiId,omitempty"`
	} `json:"requestContext,omitempty"`
	Body            string `json:"body,omitempty"`
	IsBase64Encoded bool   `json:"isBase64Encoded,omitempty"`
}

// Lambda handler output data
// Ref: Output format of a Lambda function for proxy integration
// 	https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integrations.html
type LambdaOutput struct {
	IsBase64Encoded bool              `json:"isBase64Encoded"`
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
}

// Slack EventAPI Notification data
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

type HandlerResponse struct {
	StatusCode      int
	Body            interface{}
	IsBase64Encoded bool
}

// Set this to HandlerResponse.Body if errors happened.
type HandlerResponseErrorBody struct {
	Message string      `json:"message"`
	Detail  interface{} `json:"detail"`
}

// Create request struct corresponding to input.
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

	log.Debug(req)

	return req, nil
}

// Provides handlers to each request.
type RequestHandler interface {
	Handler(ctx context.Context) *HandlerResponse
}

func makeErrorHandlerResponse(message string, err error) *HandlerResponse {
	body := &HandlerResponseErrorBody{
		Message: message,
	}
	if err != nil {
		body.Detail = err.Error()
	}
	return &HandlerResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       body,
	}
}

// Setup config.
func setConfig() {
	configType := os.Getenv("APP_CONFIG_TYPE")
	if configType == "" {
		configType = "env"
	}

	log.Info(fmt.Sprintf("set config type=%s", configType))

	if configType == "env" {
		// Get .env path
		appDir, err := fileutil.GetAppDir()
		if err != nil {
			panic(errorutil.MakePanicMessage("need app directory path."))
		}
		names := make([]string, 0)
		path := filepath.Join(appDir, ".env")
		if fileutil.FileExists(path) {
			names = append(names, path)
		}
		config.SetConfig(driver.NewEnvConfig(names...))
	}
}

// Setup DI container by env.
func setDi() {
	env := config.Get("APP_ENV", "dev")

	log.Info(fmt.Sprintf("set di env=%s", env))

	if env == "prod" {
		di.SetDi(&container.Production{})
	} else if env == "dev" {
		di.SetDi(&container.Development{})
	} else if env == "test" {
		di.SetDi(&container.Test{})
	}
}

// Lambda callback
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func LambdaHandleRequest(ctx context.Context, input LambdaInput) (interface{}, error) {
	log.Debug(fmt.Sprintf("handler, input=%v", input))

	setConfig()
	setDi()

	var body EventCallbackData
	err := json.Unmarshal([]byte(input.Body), &body)
	if err != nil {
		log.Error(err.Error())
		return makeErrorHandlerResponse("invalid request", ErrInvalidRequest), nil
	}

	h, err := NewRequestHandler(&body)
	if err != nil {
		log.Error(err.Error())
		return makeErrorHandlerResponse("parameter error", ErrInvalidParameters), nil
	}

	resp := h.Handler(ctx)
	log.Debug(resp)
	if resp == nil {
		return nil, errors.New("no response")
	}

	var respBody string
	if typeutil.IsStruct(resp.Body) {
		log.Debug("is struct")
		body, err := json.Marshal(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create response: %w", err)
		}
		respBody = string(body)
	} else {
		respBody = fmt.Sprintf("%v", resp.Body)
	}
	log.Debug("respBody", respBody)

	output := LambdaOutput{
		IsBase64Encoded: resp.IsBase64Encoded,
		StatusCode:      resp.StatusCode,
		Body:            respBody,
	}

	log.Debug(fmt.Sprintf("handler, output=%v", output))

	return output, nil
}
