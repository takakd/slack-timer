package lambdahandler

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slacktimer/internal/app/adapter/notifycontroller"
	"slacktimer/internal/app/driver/di"
	"slacktimer/internal/app/driver/di/container"
	"slacktimer/internal/pkg/config"
	"slacktimer/internal/pkg/config/driver"
	"slacktimer/internal/pkg/errorutil"
	"slacktimer/internal/pkg/fileutil"
	"slacktimer/internal/pkg/log"
)

// Lambda handler input data
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/with-sqs.html
type LambdaInput struct {
	Records []*SqsMessage `json:"records"`
}

type SqsMessage struct {
	MessageId     string            `json:"messageId"`
	ReceiptHandle string            `json:"receiptHandle"`
	Body          string            `json:"body":`
	Attributes    map[string]string `json:"attributes"`
	// TODO: check schema
	MessageAttributes map[string]interface{} `json:"messageAttributes"`
	MD5OfBody         string                 `json:"md5OfBody"`
	EventSource       string                 `json:"eventSource"`
	EventSourceArn    string                 `json:"eventSourceARN"`
	AwsRegion         string                 `json:"awsRegion"`
}

func (s *SqsMessage) HandlerInput() *notifycontroller.HandlerInput {
	return &notifycontroller.HandlerInput{
		UserId: s.Body,
		// TODO: if it needs.
		Message: "",
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
func NotifyLambdaHandler(ctx context.Context, input LambdaInput) error {
	log.Debug(fmt.Sprintf("handler, input=%v", input))

	setConfig()
	setDi()

	count := 0
	for _, m := range input.Records {
		h := notifycontroller.NewHandler()
		resp := h.Handler(ctx, m.HandlerInput())
		if resp.Error != nil {
			count++
		}
	}

	if count > 0 {
		return fmt.Errorf("error happend count=%d", count)
	}

	return nil
}
