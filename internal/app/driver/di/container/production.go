package container

import (
	"slacktimer/internal/app/driver/repository"
	"slacktimer/internal/app/usecase/updatetimerevent"
)

type Production struct {
}

// Returns interfaces in production environment.
func (d *Production) Get(name string) interface{} {
	var c interface{}

	if name == "UpdateTimerEvent" {
		c = updatetimerevent.NewUsecase()
	} else if name == "Repository" {
		c = repository.NewPostgresRepository()
	}

	return c
}
