package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"slacktimer/internal/app/driver/lambdahandler/enqueue"
)

func main() {
	lambda.Start(enqueue.LambdaHandleEvent)
}
