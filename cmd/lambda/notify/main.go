package main

import (
	"slacktimer/internal/app/driver/lambdahandler/notify"

	"github.com/aws/aws-lambda-go/lambda"
)

// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func main() {
	h := notify.NewNotifyLambdaHandler()
	lambda.Start(h.Handle)
}
