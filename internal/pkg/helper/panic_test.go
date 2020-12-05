package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakePanicMessage(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		s := MakePanicMessage("Hi")
		assert.Equal(t, s, "PANIC: Hi")
	})

	t.Run("ng", func(t *testing.T) {
		ok := DoesTestCallPanic(func() {
			MakePanicMessage(nil)
		})
		assert.True(t, ok)
	})
}
