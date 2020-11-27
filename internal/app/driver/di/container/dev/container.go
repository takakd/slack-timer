package dev

import (
	"slacktimer/internal/app/driver/queue"
	"slacktimer/internal/app/driver/repository"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/usecase/updatetimerevent"
)

type Container struct {
}

// Returns interfaces in development environment.
func (d *Container) Get(name string) interface{} {
	if c, ok := getSetTimerConcrete(name); ok {
		return c
	}
	if c, ok := getEnqueueConcrete(name); ok {
		return c
	}
	return nil
}

func getSetTimerConcrete(name string) (interface{}, bool) {
	var c interface{}
	switch name {
	case "UpdateTimerEvent":
		c = updatetimerevent.NewUsecase()
	case "Repository":
		c = repository.NewDynamoDbRepository(nil)
	case "Queue":
		c = queue.NewSQSMessageQueue(nil)
		//case "UpdateTimerEventOutputPort":
	}
	return c, c != nil
}

func getEnqueueConcrete(name string) (interface{}, bool) {
	var c interface{}
	switch name {
	case "enqueuecontroller.EnqueueNotification":
		c = enqueueevent.NewUsecase()
	case "enqueueevent.Repository":
		c = repository.NewDynamoDbRepository(nil)
	case "enqueueevent.OutputPort":
		c = enqueueevent.NewCloudWatchLogsOutputPort()
	case "enqueueevent.Queue":
		c = queue.NewSQSMessageQueue(nil)
	}
	return c, c != nil
}
