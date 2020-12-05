package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsStruct(t *testing.T) {
	cases := []struct {
		name  string
		value interface{}
		want  bool
	}{
		{"struct", struct{ test string }{"test"}, true},
		{"struct pointer", &struct{ test string }{"test"}, true},
		{"int", 1, false},
		{"string", "test", false},
		{"float64", 0.123, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, IsStruct(c.value))
		})
	}

	t.Run("struct pointer", func(t *testing.T) {
		assert.Equal(t, true, IsStruct(&struct{ test string }{"test"}))
	})

	t.Run("other pointer", func(t *testing.T) {
		var ps *string
		var pi *int

		s := "test"
		ps = &s
		i := 100
		pi = &i

		assert.Equal(t, false, IsStruct(ps))
		assert.Equal(t, false, IsStruct(pi))
	})
}
