package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"slacktimer/internal/app/adapter/enqueuecontroller"
)

func main() {
	lambda.Start(enqueuecontroller.LambdaHandleEvent)
}
