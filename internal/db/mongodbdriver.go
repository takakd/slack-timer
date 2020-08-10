package db

import (
	"proteinreminder/internal/ioc"
)

type MongoDbDriver struct {
}

func NewMongoDbDriver() *MongoDbDriver {
	d := &MongoDbDriver{}
	return d
}

func (d *MongoDbDriver) ReadString(key string) string {
	logger := ioc.GetLogger()
	logger.Debug("called ReadStrng")

	return "need to implement."
}

func (d *MongoDbDriver) WriteString(key, value string) {
	logger := ioc.GetLogger()
	logger.Debug("called WriteString")
}
