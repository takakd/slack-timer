package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakePanicMessage(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		s := NewPanicMessage("Hi")
		assert.Equal(t, s, "PANIC: Hi")
	})
}
