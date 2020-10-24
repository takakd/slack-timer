package container

import (
	"proteinreminder/internal/app/driver/repository"
	"proteinreminder/internal/app/usecase/updateproteinevent"
)

type Development struct {
}

// Returns interfaces in development environment.
func (d *Development) Get(name string) interface{} {
	var c interface{}

	if name == "UpdateProteinEvent" {
		c = updateproteinevent.NewUsecase()
	} else if name == "Repository" {
		c = repository.NewPostgresRepository()
	} else if name == "UpdateProteinEventOutputPort" {

	}

	return c
}
