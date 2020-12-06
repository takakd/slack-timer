package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"slacktimer/internal/app/driver/lambdahandler/enqueue"
)

// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func main() {
	h := enqueue.NewEnqueueLambdaHandler()
	lambda.Start(h.Handle)
}
