package testutil

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
)

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

// Returns the handler response of GET.
func HttpGetRequest(h http.HandlerFunc, target string, args url.Values) *httptest.ResponseRecorder {
	target = fmt.Sprintf("%s?%s", target, args.Encode())
	r := httptest.NewRequest(http.MethodGet, target, nil)
	w := httptest.NewRecorder()
	h(w, r)
	return w
}

// Returns the handler respose of application/x-www-form-urlencoded POST.
func HttpFormPostRequest(h http.HandlerFunc, target string, body io.Reader) *httptest.ResponseRecorder {
	r := httptest.NewRequest(http.MethodPost, target, body)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	h(w, r)
	return w
}
