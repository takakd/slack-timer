package validator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSet_SetGetError(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		bag := NewValidateErrorBag()

		bag.SetError("test", "test summary", ErrEmpty)
		error, errorExists := bag.GetError("test")
		assert.True(t, errorExists)
		assert.Equal(t, "test summary", error.Summary)
	})
}

func TestSet_ContainsError(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		bag := NewValidateErrorBag()
		bag.SetError("test", "test summary", ErrEmpty)
		assert.True(t, bag.ContainsError("test", ErrEmpty))
		assert.False(t, bag.ContainsError("not in", ErrEmpty))
	})
}
