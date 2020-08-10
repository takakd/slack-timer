package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"proteinreminder/internal/httputil"
	"proteinreminder/internal/testutil"
	"strings"
	"testing"
)

//
// POST slack-callback
//

func TestSlackCallbackPostRequest_parse(t *testing.T) {
	// https://golang.org/src/net/http/request_test.go
	body := strings.NewReader(`text=protein 11:00&user_id=user1`)
	httpReq := httptest.NewRequest(http.MethodPost, "/", body)
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	req := NewSlackCallbackRequest(httpReq)
	err := req.parse()
	if err != nil {
		t.Error(testutil.MakeTestMessageWithGotWant(err, nil))
	}
	if req.keyword != "protein" {
		t.Error(testutil.MakeTestMessageWithGotWant(req.keyword, "protein"))
	}
	if req.datetime.Hour() != 11 {
		t.Error(testutil.MakeTestMessageWithGotWant(req.datetime.Hour(), 11))
	}
	if req.datetime.Minute() != 0 {
		t.Error(testutil.MakeTestMessageWithGotWant(req.datetime.Minute(), 0))
	}
}

func TestSlackCallbackPostRequest_validate(t *testing.T) {
	// https://golang.org/src/net/http/request_test.go
	body := strings.NewReader(`text=protein 11:00&user_id=user1`)
	httpReq := httptest.NewRequest(http.MethodPost, "/", body)
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req := NewSlackCallbackRequest(httpReq)
	req.parse()
	ok, _ := req.validate()
	if !ok {
		t.Error(testutil.MakeTestMessageWithGotWant(ok, true))
	}

	body = strings.NewReader(`text=&user_id=`)
	httpReq = httptest.NewRequest(http.MethodPost, "/", body)
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	req = NewSlackCallbackRequest(httpReq)
	req.parse()
	ok, errors := req.validate()
	if ok {
		t.Error(testutil.MakeTestMessageWithGotWant(ok, false))
	}
	if !errors.ContainsError("keyword", Empty) {
		t.Error(testutil.MakeTestMessageWithGotWant(false, true))
	}
	error, exists := errors.GetError("keyword")
	if !exists {
		t.Error(testutil.MakeTestMessageWithGotWant(false, true))
	}
	if error.Summary != "need keyword." {
		t.Error(testutil.MakeTestMessageWithGotWant(error.Summary, "need keyword."))
	}
	if !errors.ContainsError("user_id", Empty) {
		t.Error(testutil.MakeTestMessageWithGotWant(false, true))
	}
	error, exists = errors.GetError("user_id")
	if !exists {
		t.Error(testutil.MakeTestMessageWithGotWant(false, true))
	}
	if error.Summary != "need user_id." {
		t.Error(testutil.MakeTestMessageWithGotWant(error.Summary, "need user_id."))
	}
}

//
func TestSlackCallbackHandler(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", SlackCallbackHandler)

	// OK case
	// https://golang.org/src/net/http/request_test.go
	w := httptest.NewRecorder()
	body := strings.NewReader(`text=protein 11:00&user_id=user1`)
	req, _ := http.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Error(testutil.MakeTestMessageWithGotWant(w.Code, 200))
	}
	respStr, _ := httputil.GetResponseBody(w.Result())

	testResp := &SlackCallbackResponse{
		Message: "success",
	}
	testRespStr, _ := json.Marshal(testResp)

	if !bytes.Equal(respStr, testRespStr) {
		t.Error(testutil.MakeTestMessageWithGotWant(string(respStr), string(testRespStr)))
	}

	// NG case
	w = httptest.NewRecorder()
	body = strings.NewReader(`text=&user_id=`)
	req, _ = http.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	mux.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Error(testutil.MakeTestMessageWithGotWant(w.Code, 400))
	}
	testErrRespStr := MakeErrorCallbackResponseBody("parameter error", ErrorCode1)
	respStr, _ = httputil.GetResponseBody(w.Result())
	if bytes.Compare(testErrRespStr, respStr) != 0 {
		t.Error(testutil.MakeTestMessageWithGotWant(string(respStr), string(testErrRespStr)))
	}
}
