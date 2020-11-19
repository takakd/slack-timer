package container

import (
	"slacktimer/internal/app/driver/repository"
	"slacktimer/internal/app/usecase/updatetimerevent"
)

type Development struct {
}

// Returns interfaces in development environment.
func (d *Development) Get(name string) interface{} {
	var c interface{}

	if name == "UpdateTimerEvent" {
		c = updatetimerevent.NewUsecase()
	} else if name == "Repository" {
		c = repository.NewDynamoDbRepository(nil)
	} else if name == "UpdateTimerEventOutputPort" {

	}

	return c
}
