package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"slacktimer/internal/app/driver/lambdahandler/settime"
)

// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func main() {
	h := settime.NewSetTimeLambdaHandler()
	lambda.Start(h.Handle)
}
