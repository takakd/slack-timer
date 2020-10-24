package errorutil

import (
	"github.com/stretchr/testify/assert"
	"slacktimer/internal/pkg/testutil"
	"testing"
)

func TestMakePanicMessage(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		s := MakePanicMessage("Hi")
		assert.Equal(t, s, "PANIC: Hi")
	})

	t.Run("ng", func(t *testing.T) {
		ok := testutil.DoesTestCallPanic(func() {
			MakePanicMessage(nil)
		})
		assert.True(t, ok)
	})
}
