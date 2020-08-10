package model

import (
	"proteinreminder/internal/testutil"
	"testing"
	"time"
)

// Repository

func TestGetProteinEvent(t *testing.T) {
	GetProteinEvent("user_id")
}

func TestSaveProteinEvent(t *testing.T) {
	p := NewProteinEvent("user_id")
	SaveProteinEvent(p)
}

// Entity

func TestNewProteinEvent(t *testing.T) {
	NewProteinEvent("user_id")
}

func TestProteinEvent_RecordedTime(t *testing.T) {
	p := NewProteinEvent("user_id")
	testTime := time.Now()
	p.setRecordedTime(testTime)
	if testTime != p.RecordedTime() {
		t.Errorf(testutil.MakeTestMessageWithGotWant(p.RecordedTime(), testTime))
	}
}
