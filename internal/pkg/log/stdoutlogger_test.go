package log

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"regexp"
	"slacktimer/internal/pkg/testutil"
	"testing"
)

func TestNewStdoutLogger(t *testing.T) {
	testutil.DoesTestCallPanic(func() {
		NewStdoutLogger()
	})
}

func TestStdoutLogger_Print(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{name: "ok", want: "test log."},
		{name: "ng", want: ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			pattern := ""
			if c.want != "" {
				dateTimePattern := "\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2}"
				pattern = fmt.Sprintf("%s %s", dateTimePattern, regexp.QuoteMeta(c.want))
			}

			// Ref: https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			logger := NewStdoutLogger()
			logger.Print(c.want)

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			got := buf.String()

			re := regexp.MustCompile(pattern)
			assert.True(t, re.MatchString(got))
		})
	}
}
