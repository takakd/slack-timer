package dev

import (
	"slacktimer/internal/app/adapter/enqueue"
	"slacktimer/internal/app/adapter/notify"
	"slacktimer/internal/app/adapter/slackhandler"
	"slacktimer/internal/app/driver/queue"
	"slacktimer/internal/app/driver/repository"
	"slacktimer/internal/app/driver/slack"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/usecase/notifyevent"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log/driver"
)

type Container struct {
}

var _ di.DI = (*Container)(nil)

// Returns interfaces in development environment.
func (d *Container) Get(name string) interface{} {
	if c, ok := getUtilConcrete(name); ok {
		return c
	}
	if c, ok := getSetTimerConcrete(name); ok {
		return c
	}
	if c, ok := getEnqueueConcrete(name); ok {
		return c
	}
	if c, ok := getNotifyConcrete(name); ok {
		return c
	}
	return nil
}

func getUtilConcrete(name string) (interface{}, bool) {
	var c interface{}
	switch name {
	case "Logger":
		c = driver.NewCloudWatchLogger()
	}
	return c, c != nil
}

func getSetTimerConcrete(name string) (interface{}, bool) {
	var c interface{}
	switch name {
	case "UpdateTimerEvent":
		c = updatetimerevent.NewInteractor()
	case "Repository":
		c = repository.NewDynamoDb(nil)
	case "Queue":
		c = queue.NewSqs(nil)
	}
	return c, c != nil
}

func getEnqueueConcrete(name string) (interface{}, bool) {
	var c interface{}
	switch name {
	case "enqueuecontroller.InputPort":
		c = enqueueevent.NewInteractor()
	case "enqueueevent.OutputPort":
		c = enqueue.NewCloudWatchLogsPresenter()
	case "enqueueevent.Repository":
		c = repository.NewDynamoDb(nil)
	case "enqueueevent.Queue":
		c = queue.NewSqs(nil)
	}
	return c, c != nil
}

func getNotifyConcrete(name string) (interface{}, bool) {
	var c interface{}
	switch name {
	case "notifycontroller.InputPort":
		c = notifyevent.NewInteractor()
	case "notifyevent.OutputPort":
		c = notify.NewCloudWatchLogsPresenter()
	case "notifyevent.Repository":
		c = repository.NewDynamoDb(nil)
	case "notifyevent.Notifier":
		c = slackhandler.NewSlackHandler()
	case "slackhandler.SlackApi":
		c = slack.NewSlackApiDriver()
	}
	return c, c != nil
}
