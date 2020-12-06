package main

import (
	"slacktimer/internal/app/driver/lambdahandler/settime"

	"github.com/aws/aws-lambda-go/lambda"
)

// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func main() {
	h := settime.NewSetTimeLambdaHandler()
	lambda.Start(h.Handle)
}
