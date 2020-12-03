// Package slackcontroller provides the slack Event API callback handler.
// Ref: https://api.slack.com/events-api#the-events-api__receiving-events
package slackcontroller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/app/util/appinit"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
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

	req := &SetRequestHandler{
		messageEvent: &data.MessageEvent,
		usecase:      usecase,
	}

	return req, nil
}

// Provides handlers to each request.
type RequestHandler interface {
	Handler(ctx context.Context) *HandlerResponse
}

func makeErrorHandlerResponse(message string, detail string) *HandlerResponse {
	body := &HandlerResponseErrorBody{
		Message: message,
	}
	if detail != "" {
		body.Detail = detail
	}
	return &HandlerResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       body,
	}
}

// Lambda callback
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func LambdaHandleRequest(ctx context.Context, input LambdaInput) (interface{}, error) {
	appinit.AppInit()

	log.Info("handler input", input)

	var body EventCallbackData
	err := json.Unmarshal([]byte(input.Body), &body)
	if err != nil {
		log.Info(fmt.Errorf("invalid request: %w", err))
		return makeErrorHandlerResponse("invalid request", "parameters are wrong"), nil
	}

	h, err := NewRequestHandler(&body)
	if err != nil {
		log.Info(fmt.Errorf("invalid parameter: %w", err))
		return makeErrorHandlerResponse("invalid parameter", ""), nil
	}

	resp := h.Handler(ctx)

	var respBody string
	if typeutil.IsStruct(resp.Body) {
		body, err := json.Marshal(resp.Body)
		if err != nil {
			log.Info(fmt.Errorf("internal sesrver error: %w", err))
			return nil, errors.New("internal server error")
		}
		respBody = string(body)
	} else {
		respBody = fmt.Sprintf("%v", resp.Body)
	}

	output := LambdaOutput{
		IsBase64Encoded: resp.IsBase64Encoded,
		StatusCode:      resp.StatusCode,
		Body:            respBody,
	}

	log.Info("handler output", output)

	return output, nil
}
