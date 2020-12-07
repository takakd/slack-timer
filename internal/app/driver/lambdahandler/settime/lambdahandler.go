// Package settime the entry point of setting event time lambda function.
package settime

import (
	"context"
	"encoding/json"
	"slacktimer/internal/app/adapter/settime"
)

// LambdaHandler defines the interface called by AWS Lambda.
type LambdaHandler interface {
	Handle(ctx context.Context, input LambdaInput) (*LambdaOutput, error)
}

// LambdaInput is passed from API Gateway.
type LambdaInput struct {
	Resource                        string                `json:"resource,omitempty"`
	Path                            string                `json:"path,omitempty"`
	HTTPMethod                      string                `json:"httpMethod,omitempty"`
	Headers                         map[string]string     `json:"headers,omitempty"`
	MultiValueHeaders               map[string][]string   `json:"multiValueHeaders,omitempty"`
	QueryStringParameters           map[string]string     `json:"queryStringParameters,omitempty"`
	MultiValueQueryStringParameters []map[string][]string `json:"multiValueQueryStringParameters,omitempty"`
	PathParameters                  map[string]string     `json:"pathParameters,omitempty"`
	StageVaribales                  map[string]string     `json:"stageVariables,omitempty"`
	RequestContext                  struct {
		AccountID  string `json:"accountId,omitempty"`
		ResourceID string `json:"resourceId,omitempty"`
		Stage      string `json:"stage,omitempty"`
		RequestID  string `json:"requestId,omitempty"`
		Identity   struct {
			CognitoIdentityPoolID         string `json:"cognitoIdentityPoolId,omitempty"`
			AccountID                     string `json:"accountId,omitempty"`
			CognitoIdentityID             string `json:"cognitoIdentityId,omitempty"`
			Caller                        string `json:"caller,omitempty"`
			APIKey                        string `json:"apiKey,omitempty"`
			SourceIP                      string `json:"sourceIp,omitempty"`
			CognitoAuthenticationType     string `json:"cognitoAuthenticationType,omitempty"`
			CognitoAuthenticationProvider string `json:"cognitoAuthenticationProvider,omitempty"`
			UserArn                       string `json:"userArn,omitempty"`
			UserAgent                     string `json:"userAgent,omitempty"`
			User                          string `json:"user,omitempty"`
		} `json:"identity,omitempty"`
		ResourcePath string `json:"resourcePath,omitempty"`
		HTTPMethod   string `json:"httpMethod,omitempty"`
		APIID        string `json:"apiId,omitempty"`
	} `json:"requestContext,omitempty"`
	Body            string `json:"body,omitempty"`
	IsBase64Encoded bool   `json:"isBase64Encoded,omitempty"`
}

// HandleInput convert to the data for controller.
func (l LambdaInput) HandleInput() (data *settime.HandleInput, err error) {
	// Extract Slack event data.
	var body settime.EventCallbackData
	err = json.Unmarshal([]byte(l.Body), &body)
	if err != nil {
		return
	}

	// Call controller method.
	data = &settime.HandleInput{
		EventData: body,
	}
	return
}

// LambdaOutput is the type of Lambda function returns.
// Ref: Output format of a Lambda function for proxy integration
// 	https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integrations.html
type LambdaOutput struct {
	IsBase64Encoded bool              `json:"isBase64Encoded"`
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
}
