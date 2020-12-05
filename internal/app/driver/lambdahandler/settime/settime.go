package settime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slacktimer/internal/app/adapter/slackcontroller"
	"slacktimer/internal/app/util/appinit"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"slacktimer/internal/pkg/typeutil"
)

type LambdaHandler interface {
	LambdaHandler(ctx context.Context, input LambdaInput) (*LambdaOutput, error)
}

type SetTimerLambdaHandler struct {
	ctrl slackcontroller.Handler
}

func NewSetTimerLambdaHandler() LambdaHandler {
	h := &SetTimerLambdaHandler{}
	h.ctrl = di.Get("settime.Handler").(slackcontroller.Handler)
	return h
}

// Lambda handler input data
// API Gateway passes this.
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

// To a controller input data.
func (s *LambdaInput) HandlerInput() (data *slackcontroller.HandlerInput, err error) {
	// Extract Slack event data.
	var body slackcontroller.EventCallbackData
	err = json.Unmarshal([]byte(s.Body), &body)
	if err != nil {
		return
	}

	// Call controller method.
	data = &slackcontroller.HandlerInput{
		EventData: body,
	}
	return
}

// Lambda handler output data
// this lambda function returns this.
// Ref: Output format of a Lambda function for proxy integration
// 	https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integrations.html
type LambdaOutput struct {
	IsBase64Encoded bool              `json:"isBase64Encoded"`
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
}

// TODO: deprecated
//// Set this to LambdaOutput.Body as JSON if errors happened.
//type HandlerResponseErrorBody struct {
//	Message string      `json:"message"`
//	Detail  interface{} `json:"detail"`
//}

// API Gateway calls this function.
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func (s *SetTimerLambdaHandler) LambdaHandler(ctx context.Context, input LambdaInput) (*LambdaOutput, error) {
	appinit.AppInit()

	log.Info("lambda handler input", input)

	// Extract Slack event data.
	data, err := input.HandlerInput()
	if err != nil {
		log.Info(fmt.Errorf("invalid request: %w", err))
		o := &LambdaOutput{
			IsBase64Encoded: true,
			StatusCode:      http.StatusInternalServerError,
			// Not have to do use json.Marshal.
			Body: `{"message":"invalid request", "detail":"parameters are wrong"}`,
		}
		return o, nil
	}

	// Call controller method.
	// TOOD: di
	resp := s.ctrl.Handler(ctx, *data)

	// Create a response.
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

	output := &LambdaOutput{
		//IsBase64Encoded: resp.IsBase64Encoded,
		// TODO: check
		IsBase64Encoded: true,
		StatusCode:      resp.StatusCode,
		Body:            respBody,
	}

	log.Info("handler output", output)

	return output, nil
}
