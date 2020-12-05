package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"slacktimer/internal/app/driver/lambdahandler/settime"
)

func main() {
	h := settime.NewSetTimeLambdaHandler()
	lambda.Start(h.Handle)
}
