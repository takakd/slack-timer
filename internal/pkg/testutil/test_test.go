package testutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeTestMessageWithGotWant(t *testing.T) {
	s := MakeTestMessageWithGotWant("Hi", "Hello")
	assert.Equal(t, "got: Hi, want: Hello", s)
}

func TestDoesTestCallPanic(t *testing.T) {
	isCalled := DoesTestCallPanic(func() {
		var i interface{}
		if i == nil {
			panic("Hi, panic.")
		}
	})
	assert.True(t, isCalled)
}
