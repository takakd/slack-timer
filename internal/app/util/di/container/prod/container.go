package prod

import (
	"slacktimer/internal/app/adapter/enqueuecontroller"
	"slacktimer/internal/app/driver/queue"
	"slacktimer/internal/app/driver/repository"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/usecase/updatetimerevent"
)

type Container struct {
}

// Returns interfaces in production environment.
func (d *Container) Get(name string) interface{} {
	if name == "UpdateTimerEvent" {
		return updatetimerevent.NewUsecase()
	} else if name == "Repository" {
		return repository.NewPostgresRepository()
	}

	if c, ok := getEnqueueConcrete(name); ok {
		return c
	}

	return nil
}

func getEnqueueConcrete(name string) (interface{}, bool) {
	var c interface{}
	switch name {
	case "enqueuecontroller.InputPort":
		c = enqueueevent.NewInteractor()
	case "enqueueevent.OutputPort":
		c = enqueuecontroller.NewCloudWatchLogsOutputPort()
	case "enqueueevent.Repository":
		c = repository.NewDynamoDbRepository(nil)
	case "enqueueevent.Queue":
		c = queue.NewSQSMessageQueue(nil)
	}
	return c, c != nil
}
