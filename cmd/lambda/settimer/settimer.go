package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"slacktimer/internal/app/adapter/slackcontroller"
)

func main() {
	lambda.Start(slackcontroller.LambdaHandleRequest)
}
