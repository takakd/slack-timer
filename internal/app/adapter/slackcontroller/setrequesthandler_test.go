package slackcontroller

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSetRequestHandler_validate(t *testing.T) {
	cases := []struct {
		name  string
		text  string
		min   time.Duration
		valid bool
	}{
		{"OK", "set 10", 10, true},
		{"OK", "set 1", 1, true},
		{"NG", "set -1", 0, false},
		{"NG", "set", 0, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := SetRequestHandler{
				params: &SlackCallbackRequestParams{
					UserId: "test",
					Text: c.text,
				},
				datetime: time.Now(),
			}
			bag := r.validate()
			_, exists := bag.GetError("interval")
			assert.Equal(t, c.valid, !exists)
			if c.valid {
				assert.Equal(t, c.min, r.remindIntervalInMin)
			}
		})
	}
}

func TestSetRequestHandler_Handler(t *testing.T) {
	// TODO:

}

//
//import (
//	"bytes"
//	"github.com/stretchr/testify/require"
//	"net/http"
//	"net/http/httptest"
//	"proteinreminder/internal/pkg/testutil"
//	"strings"
//	"testing"
//	"time"
//	"context"
//	"net/url"
//	"github.com/stretchr/testify/assert"
//	"proteinreminder/internal/app/usecase"
//	"github.com/golang/mock/gomock"
//)
//
////
//// POST slack-callback
////
//
//func TestParseRequest(t *testing.T) {
//	cases := []struct {
//		name    string
//		text    string
//		subType CommandSubType
//	}{
//		{"set", "set 1", SubTypeSet},
//		{"got", "got", SubTypeGot},
//		{"nil", "invalid", ""},
//	}
//	for _, c := range cases {
//		t.Run(c.name, func(t *testing.T) {
//			// https://golang.org/src/net/http/request_test.go
//			body := strings.NewReader(`text=` + c.text)
//			httpReq := httptest.NewRequest(http.MethodPost, "/", body)
//			httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
//
//			req, err := parseRequest(httpReq)
//			if c.subType != "" {
//				if err != nil {
//					t.Error(testutil.MakeTestMessageWithGotWant(err, nil))
//				}
//				require.NotNil(t, req)
//			}
//
//			switch r := req.(type) {
//			case *SlackCallbackGotRequest:
//				if r.subType != SubTypeGot {
//					t.Error(testutil.MakeTestMessageWithGotWant(r.subType, SubTypeGot))
//				}
//			case *SlackCallbackSetRequest:
//				if r.subType != SubTypeSet {
//					t.Error(testutil.MakeTestMessageWithGotWant(r.subType, SubTypeSet))
//				}
//			default:
//				if c.subType != "" {
//					t.Errorf("wrong type. type=%v,%v", r, c.subType)
//				}
//			}
//		})
//	}
//}
//
//func TestValidate(t *testing.T) {
//	cases := []struct {
//		name       string
//		body       string
//		containErr bool
//	}{
//		{"OK", "text=got&user_id=abc", false},
//		{"NG", "text=got&user_id=", true},
//	}
//	for _, c := range cases {
//		t.Run(c.name, func(t *testing.T) {
//			body := strings.NewReader(c.body)
//			httpReq := httptest.NewRequest(http.MethodPost, "/", body)
//			httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
//			req, _ := parseRequest(httpReq)
//			r, _ := req.(*SlackCallbackGotRequest)
//			_, bag := r.SlackCallbackRequest.validate()
//			if c.containErr != bag.ContainsError("user_id", Empty) {
//				t.Error(testutil.MakeTestMessageWithGotWant(bag.ContainsError("user_id", Empty), c.containErr))
//			}
//		})
//	}
//}
//
//func TestSlackCallbackGotRequest_validate(t *testing.T) {
//	cases := []struct {
//		name    string
//		userId  string
//		valid   bool
//		errType ErrorType
//	}{
//		{"OK", "id1", true, ""},
//		{"NG", "", false, Empty},
//		{"NG:not match error type", "", false, ""},
//	}
//	for _, c := range cases {
//		t.Run(c.name, func(t *testing.T) {
//			r := SlackCallbackGotRequest{
//				SlackCallbackRequest: SlackCallbackRequest{
//					request: nil,
//					params: SlackCallbackRequestParams{
//						UserId: c.userId,
//					},
//				},
//			}
//			valid, bag := r.validate()
//			if valid != c.valid {
//				t.Error(testutil.MakeTestMessageWithGotWant(valid, c.valid))
//			}
//			if !c.valid {
//				if !bag.ContainsError("user_id", c.errType) {
//					t.Errorf("should has %s error.", c.errType)
//				}
//			}
//
//		})
//	}
//}
//
//func TestMakeSlackCallbackSetRequest(t *testing.T) {
//	cases := []struct {
//		name    string
//		text    string
//		min     time.Duration
//		valid   bool
//		wantErr string
//	}{
//		{"OK", "set 10", 10, true, ""},
//		{"OK", "set 1", 1, true, ""},
//		{"NG", "set -1", 0, false, "invalid Text format"},
//		{"NG", "set", 0, false, "invalid Text format"},
//	}
//	validReq := &SlackCallbackRequest{
//		request: nil,
//		params: SlackCallbackRequestParams{
//			UserId: "id1",
//		},
//	}
//	for _, c := range cases {
//		t.Run(c.name, func(t *testing.T) {
//			validReq.params.Text = c.text
//			r, err := MakeSlackCallbackSetRequest(validReq)
//			if err == nil && c.wantErr != "" {
//				t.Error(testutil.MakeTestMessageWithGotWant(err, c.wantErr))
//				return
//
//			} else if err != nil {
//				if err.Error() != c.wantErr {
//					t.Error(testutil.MakeTestMessageWithGotWant(err.Error(), c.wantErr))
//				}
//				return
//			}
//
//			if r.remindIntervalInMin != c.min {
//				t.Error(testutil.MakeTestMessageWithGotWant(c.min, r.remindIntervalInMin))
//			}
//		})
//	}
//}
//
//func TestSlackCallbackSetRequest_validate(t *testing.T) {
//	t.Log("TestSlackCallbackGotRequest_validate covers this test.")
//}
//
//func TestMakeErrorCallbackResponseBody(t *testing.T) {
//	t.Run("OK", func(t *testing.T) {
//		got := makeErrorCallbackResponseBody("hi", 123)
//		want := []byte(`{"message":"hi","code":123}`)
//		if !bytes.Equal(got, want) {
//			t.Error(testutil.MakeTestMessageWithGotWant(string(want), string(got)))
//		}
//	})
//}
//
////func TestSlackCallbackPostRequest_validate(t *testing.T) {
////	// https://golang.org/src/net/http/request_test.go
////	body := strings.NewReader(`text=protein 11:00&user_id=user1`)
////	httpReq := httptest.NewRequest(http.MethodPost, "/", body)
////	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
////	req := NewSlackCallbackRequest(httpReq)
////	req.parse()
////	ok, _ := req.validate()
////	if !ok {
////		t.Error(testutil.MakeTestMessageWithGotWant(ok, true))
////	}
////
////	body = strings.NewReader(`text=&user_id=`)
////	httpReq = httptest.NewRequest(http.MethodPost, "/", body)
////	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
////
////	req = NewSlackCallbackRequest(httpReq)
////	req.parse()
////	ok, errors := req.validate()
////	if ok {
////		t.Error(testutil.MakeTestMessageWithGotWant(ok, false))
////	}
////	if !errors.ContainsError("keyword", Empty) {
////		t.Error(testutil.MakeTestMessageWithGotWant(false, true))
////	}
////	error, exists := errors.GetError("keyword")
////	if !exists {
////		t.Error(testutil.MakeTestMessageWithGotWant(false, true))
////	}
////	if error.Summary != "need keyword." {
////		t.Error(testutil.MakeTestMessageWithGotWant(error.Summary, "need keyword."))
////	}
////	if !errors.ContainsError("user_id", Empty) {
////		t.Error(testutil.MakeTestMessageWithGotWant(false, true))
////	}
////	error, exists = errors.GetError("user_id")
////	if !exists {
////		t.Error(testutil.MakeTestMessageWithGotWant(false, true))
////	}
////	if error.Summary != "need user_id." {
////		t.Error(testutil.MakeTestMessageWithGotWant(error.Summary, "need user_id."))
////	}
////}
//
////
//func TestSlackCallbackHandler(t *testing.T) {
//	// Ref: https://golang.org/src/net/http/httptest/example_test.go#L38
//	//ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	//		fmt.Fprintln(w, "Hello, client")
//	//	}))
//	//	defer ts.Close()
//	//
//	//	res, err := http.Get(ts.URL)
//	//	if err != nil {
//	//		log.Fatal(err)
//	//	}
//	//	greeting, err := ioutil.ReadAll(res.Body)
//	//	res.Body.Close()
//	//	if err != nil {
//	//		log.Fatal(err)
//	//	}
//	//
//	//	fmt.Printf("%s", greeting)
//	//	// Output: Hello, client
//	//}
//
//	t.Run("wrong method", func(t *testing.T) {
//		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			ctx := context.TODO()
//			handler(ctx, nil, w, r)
//		})
//		w := testutil.HttpGetRequest(h, "/", url.Values{})
//		assert.Equal(t, 404, w.Result().StatusCode)
//	})
//
//	t.Run("invalid parameter", func(t *testing.T) {
//		//h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		//	ctx := context.TODO()
//		//	ctrl := gomock.NewController(t)
//		//	m := usecase.NewMockProteinEventSaver(ctrl)
//		//	handler(ctx, m, w, r)
//		//})
//		//w := testutil.HttpGetRequest(handler, "/", url.Values{})
//		//assert.Equal(t, 404, w.Result().StatusCode)
//
//		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			ctx := context.TODO()
//			handler(ctx, nil, w, r)
//		})
//
//		params := url.Values{
//			"text": {""},
//		}
//		w := testutil.HttpFormPostRequest(handler, "/", strings.NewReader(params.Encode()))
//		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
//
//		params = url.Values{
//			"text": {"abc"},
//		}
//		w = testutil.HttpFormPostRequest(handler, "/", strings.NewReader(params.Encode()))
//		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
//
//		params = url.Values{
//			"text": {"set"},
//		}
//		w = testutil.HttpFormPostRequest(handler, "/", strings.NewReader(params.Encode()))
//		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
//
//		params = url.Values{
//			"text": {"set a"},
//		}
//		w = testutil.HttpFormPostRequest(handler, "/", strings.NewReader(params.Encode()))
//		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
//
//		params = url.Values{
//			"text": {"got 1"},
//		}
//		w = testutil.HttpFormPostRequest(handler, "/", strings.NewReader(params.Encode()))
//		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
//	})
//
//	//	res, err := http.Get(ts.URL)
//	//	if err != nil {
//	//		log.Fatal(err)
//	//	}
//	//	greeting, err := ioutil.ReadAll(res.Body)
//	//	res.Body.Close()
//	//	if err != nil {
//	//		log.Fatal(err)
//	//	}
//	//
//	//	fmt.Printf("%s", greeting)
//
//	//mux := http.NewServeMux()
//	//mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//	//	SlackCallbackHandler(context.TODO(), w, r)
//	//})
//
//	// valiadtor error
//	// usecase error
//	// success
//
//	// NOTE: below this is old code
//
//	//// OK case
//	//// https://golang.org/src/net/http/request_test.go
//	//w := httptest.NewRecorder()
//	//body := strings.NewReader(`text=protein 11:00&user_id=user1`)
//	//req, _ := http.NewRequest(http.MethodPost, "/", body)
//	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
//	//mux.ServeHTTP(w, req)
//	//if w.Code != 200 {
//	//	t.Error(testutil.MakeTestMessageWithGotWant(w.Code, 200))
//	//}
//	//respStr, _ := httputil.GetResponseBody(w.Result())
//	//
//	//testResp := &SlackCallbackResponse{
//	//	Message: "success",
//	//}
//	//testRespStr, _ := json.Marshal(testResp)
//	//
//	//if !bytes.Equal(respStr, testRespStr) {
//	//	t.Error(testutil.MakeTestMessageWithGotWant(string(respStr), string(testRespStr)))
//	//}
//	//
//	//// NG case
//	//w = httptest.NewRecorder()
//	//body = strings.NewReader(`text=&user_id=`)
//	//req, _ = http.NewRequest(http.MethodPost, "/", body)
//	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
//	//mux.ServeHTTP(w, req)
//	//if w.Code != 400 {
//	//	t.Error(testutil.MakeTestMessageWithGotWant(w.Code, 400))
//	//}
//	//testErrRespStr := makeErrorCallbackResponseBody("parameter error", ErrorCode1)
//	//respStr, _ = httputil.GetResponseBody(w.Result())
//	//if bytes.Compare(testErrRespStr, respStr) != 0 {
//	//	t.Error(testutil.MakeTestMessageWithGotWant(string(respStr), string(testErrRespStr)))
//	//}
//}
