package prod

import (
	"slacktimer/internal/app/adapter/enqueue"
	"slacktimer/internal/app/adapter/notify"
	"slacktimer/internal/app/adapter/settime"
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

// Container implements DI on production env.
type Container struct {
}

var _ di.DI = (*Container)(nil)

// Get returns interfaces corresponding name.
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
	case "settime.URLVerificationRequestHandler":
		c = settime.NewURLVerificationRequestHandlerFunctor()
	case "settime.SaveEventHandler":
		c = settime.NewSaveEventHandlerFunctor()
	case "settime.ControllerHandler":
		c = settime.NewController()
	case "updatetimerevent.InputPort":
		c = updatetimerevent.NewInteractor()
	case "updatetimerevent.Repository":
		c = repository.NewDynamoDb(nil)
	}
	return c, c != nil
}

func getEnqueueConcrete(name string) (interface{}, bool) {
	var c interface{}
	switch name {
	case "enqueue.ControllerHandler":
		c = enqueue.NewController()
	case "enqueueevent.InputPort":
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
	case "notify.ControllerHandler":
		c = notify.NewController()
	case "notifyevent.InputPort":
		c = notifyevent.NewInteractor()
	case "notifyevent.OutputPort":
		c = notify.NewCloudWatchLogsPresenter()
	case "notifyevent.Repository":
		c = repository.NewDynamoDb(nil)
	case "notifyevent.Notifier":
		c = slackhandler.NewSlackHandler()
	case "slack.API":
		c = slack.NewAPIDriver()
	}
	return c, c != nil
}
