package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"slacktimer/internal/app/driver/lambdahandler"
)

func main() {
	lambda.Start(lambdahandler.NotifyLambdaHandler)
}
