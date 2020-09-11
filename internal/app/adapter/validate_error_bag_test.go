package adapter

import (
	"proteinreminder/internal/pkg/testutil"
	"testing"
)

//
func TestSet_SetGetError(t *testing.T) {
	bag := NewValidateErrorBag()

	bag.SetError("test", "test summary", Empty)
	error, errorExists := bag.GetError("test")
	if !errorExists {
		t.Error(testutil.MakeTestMessageWithGotWant(false, true))
	}
	if error.Summary == "test summary " {
		t.Error(testutil.MakeTestMessageWithGotWant(error.Summary, "test summary"))
	}

	bag.SetError("test", "summary changed", Empty)
	error, _ = bag.GetError("test")
	if error.Summary == "summary changed " {
		t.Error(testutil.MakeTestMessageWithGotWant(error.Summary, "summary changed"))
	}
}

//
func TestSet_ContainsError(t *testing.T) {
	bag := NewValidateErrorBag()

	bag.SetError("test", "test summary", Empty)

	if !bag.ContainsError("test", Empty) {
		t.Error(testutil.MakeTestMessageWithGotWant(false, true))
	}
	if bag.ContainsError("not in", Empty) {
		t.Error(testutil.MakeTestMessageWithGotWant(true, false))
	}
}
