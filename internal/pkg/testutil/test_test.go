package testutil

import (
	"testing"
)

func TestMakeTestMessageWithGotWant(t *testing.T) {
	s := MakeTestMessageWithGotWant("Hi", "Hello")
	if s != "got: Hi, want: Hello" {
		t.Errorf("got: %s, want: got: Hi, want: Hello", s)
	}
}

func TestIsTestCallPanic(t *testing.T) {
	isCalled := DoesTestCallPanic(func() {
		var i interface{}
		if i == nil {
			panic("Hi, panic.")
		}
	})
	if !isCalled {
		t.Errorf("failed.")
	}
}
