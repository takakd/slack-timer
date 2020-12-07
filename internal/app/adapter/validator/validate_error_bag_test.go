package validator

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSet_SetGetError(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		bag := NewValidateErrorBag()

		errMsg := errors.New("error")
		bag.SetError("test", "test summary", errMsg)
		error, errorExists := bag.GetError("test")
		assert.True(t, errorExists)
		assert.Equal(t, "test summary", error.Summary)
	})
}

func TestSet_ContainsError(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		bag := NewValidateErrorBag()
		errMsg := errors.New("error")
		bag.SetError("test", "test summary", errMsg)
		assert.True(t, bag.ContainsError("test", errMsg))
		assert.False(t, bag.ContainsError("not in", errMsg))
	})
}
