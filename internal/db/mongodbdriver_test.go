package db

import (
	"testing"
)

func TestNewMongoDbDriver(t *testing.T) {
	NewMongoDbDriver()
}

func TestMongoDbDriver_ReadString(t *testing.T) {
	d := NewMongoDbDriver()
	d.ReadString("test")
}

func TestMongoDbDriver_WriteString(t *testing.T) {
	d := NewMongoDbDriver()
	d.WriteString("test", "value")
}
