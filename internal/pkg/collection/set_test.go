package collection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		s := NewSet()

		s.Set("test")
		assert.True(t, s.Contains("test"))

		s.Remove("test")
		assert.False(t, s.Contains("test"))

		s.Set(1)
		assert.True(t, s.Contains(1))

		s.Remove(1)
		assert.False(t, s.Contains(1))
	})
}
