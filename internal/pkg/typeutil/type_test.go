package typeutil

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
		{"int", 1, false},
		{"string", "test", false},
		{"float64", 0.123, false},
	}
	for _, c := range cases {

		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, IsStruct(c.value))
		})
	}
}
