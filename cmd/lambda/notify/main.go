package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"slacktimer/internal/app/driver/lambdahandler/notify"
)

// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func main() {
	h := notify.NewNotifyLambdaHandler()
	lambda.Start(h.Handle)
}
