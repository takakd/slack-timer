package main

import (
	"slacktimer/internal/app/driver/lambdahandler/settime"

	"slacktimer/internal/app/util/appinitializer"

	"github.com/aws/aws-lambda-go/lambda"
)

// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func main() {
	appinitializer.AppInit()

	h := settime.NewLambdaFunctor()
	lambda.Start(h.Handle)
}
