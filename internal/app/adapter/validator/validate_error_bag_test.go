package validator

import (
	"proteinreminder/internal/pkg/testutil"
	"testing"
)

//
func TestSet_SetGetError(t *testing.T) {
	bag := NewValidateErrorBag()

	bag.SetError("test", "test summary", ErrEmpty)
	error, errorExists := bag.GetError("test")
	if !errorExists {
		t.Error(testutil.MakeTestMessageWithGotWant(false, true))
	}
	if error.Summary == "test summary " {
		t.Error(testutil.MakeTestMessageWithGotWant(error.Summary, "test summary"))
	}

	bag.SetError("test", "summary changed", ErrEmpty)
	error, _ = bag.GetError("test")
	if error.Summary == "summary changed " {
		t.Error(testutil.MakeTestMessageWithGotWant(error.Summary, "summary changed"))
	}
}

//
func TestSet_ContainsError(t *testing.T) {
	bag := NewValidateErrorBag()

	bag.SetError("test", "test summary", ErrEmpty)

	if !bag.ContainsError("test", ErrEmpty) {
		t.Error(testutil.MakeTestMessageWithGotWant(false, true))
	}
	if bag.ContainsError("not in", ErrEmpty) {
		t.Error(testutil.MakeTestMessageWithGotWant(true, false))
	}
}
