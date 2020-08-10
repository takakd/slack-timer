package errorutil

import "fmt"

func MakePanicMessage(detail interface{}) string {
	if detail == nil {
		panic("PANIC: detail should not be nil.")
	}
	message := fmt.Sprintf("PANIC: %s", detail)
	return message
}
