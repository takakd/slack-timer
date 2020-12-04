package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"slacktimer/internal/app/driver/lambdahandler/notify"
)

func main() {
	lambda.Start(notify.NotifyLambdaHandler)
}
