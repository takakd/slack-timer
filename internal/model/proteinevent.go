package model

import (
	"time"
)

type ProteinEvent struct {
	userId   string
	recorded time.Time
}

// Repository

func GetProteinEvent(userId string) *ProteinEvent {
	// TODO: get db from ioc
	p := NewProteinEvent(userId)
	return p
}

func SaveProteinEvent(p *ProteinEvent) error {
	return nil
}

// Entity

func NewProteinEvent(userId string) *ProteinEvent {
	p := &ProteinEvent{
		userId: userId,
	}
	return p
}

// NOTE: Allow this module to set the value.
func (p *ProteinEvent) setRecordedTime(value time.Time) {
	p.recorded = value
}

func (p *ProteinEvent) RecordedTime() time.Time {
	return p.recorded
}
