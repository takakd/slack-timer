package collection

import (
	"proteinreminder/internal/pkg/testutil"
	"testing"
)

func TestSet(t *testing.T) {
	s := NewSet()

	s.Set("test")
	if !s.Contains("test") {
		t.Error(testutil.MakeTestMessageWithGotWant(false, true))
	}

	s.Remove("test")
	if s.Contains("test") {
		t.Error(testutil.MakeTestMessageWithGotWant(true, false))
	}

	s.Set(1)
	if !s.Contains(1) {
		t.Error(testutil.MakeTestMessageWithGotWant(false, true))
	}

	s.Remove(1)
	if s.Contains(1) {
		t.Error(testutil.MakeTestMessageWithGotWant(true, false))
	}
}
