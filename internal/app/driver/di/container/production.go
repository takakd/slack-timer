package container

import (
	"proteinreminder/internal/app/driver/repository"
	"proteinreminder/internal/app/usecase/updateproteinevent"
)

type Production struct {
}

// Returns interfaces in production environment.
func (d *Production) Get(name string) interface{} {
	var c interface{}

	if name == "UpdateProteinEvent" {
		c = updateproteinevent.NewUsecase()
	} else if name == "Repository" {
		c = repository.NewPostgresRepository()
	}

	return c
}
