package testutil

import "fmt"

// Make test message with got and want information.
func MakeTestMessageWithGotWant(got, want interface{}) string {
	message := fmt.Sprintf("got: %v, want: %v", got, want)
	return message
}

// Testing if testFunc calls panic.
// true: called, false: not called.
// e.g.
// DoesTestCallPanic(func(){
//   <place test target here.>
// })
func DoesTestCallPanic(testFunc func()) (called bool) {
	defer func() {
		if err := recover(); err == nil {
			called = false
		}
	}()
	called = true
	testFunc()
	return
}
