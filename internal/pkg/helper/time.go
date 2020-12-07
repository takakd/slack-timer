package helper

import (
	"strconv"
	"time"
)

// ParseUnixStr parses unix timestamp string to time.Time.
func ParseUnixStr(s string) (time.Time, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Now(), err
	}
	return time.Unix(i, 0), nil
}
