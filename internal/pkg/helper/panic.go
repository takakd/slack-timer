package helper

import "fmt"

// NewPanicMessage create panic message.
func NewPanicMessage(detail interface{}) string {
	if detail == nil {
		panic("PANIC: detail should not be nil.")
	}
	message := fmt.Sprintf("PANIC: %s", detail)
	return message
}
