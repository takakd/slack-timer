package main

import (
	"slacktimer/internal/app/driver/lambdahandler/enqueue"

	"github.com/aws/aws-lambda-go/lambda"
)

// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func main() {
	h := enqueue.NewEnqueueLambdaHandler()
	lambda.Start(h.Handle)
}
