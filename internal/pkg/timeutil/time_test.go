package timeutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseUnixStr(t *testing.T) {
	cases := []struct {
		name      string
		s         string
		t         time.Time
		errExists bool
	}{
		{"ok", "1606830655", time.Unix(1606830655, 0), false},
		{"ng", "abc", time.Now(), true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tm, err := ParseUnixStr(c.s)

			if c.errExists {
				assert.Error(t, err)
			} else {
				assert.Equal(t, c.t, tm)
				assert.NoError(t, err)
			}
		})
	}
}
