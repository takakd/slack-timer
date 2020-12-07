package helper

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRequestBody(t *testing.T) {
	// ref: https://stackoverflow.com/questions/45682353/httptest-newrequest-vs-http-newrequest-which-one-to-use-in-tests-and-why

	t.Run("ok", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			body, err := GetRequestBody(r)
			if err != nil || string(body) != "hi" {
				io.WriteString(w, "ng")
			}
			io.WriteString(w, "ok")
		}
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", strings.NewReader("hi"))
		handler(w, req)
		assert.Equal(t, "ok", w.Body.String())
	})

	t.Run("ng:empty", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			body, err := GetRequestBody(r)
			if err != nil || string(body) != "" {
				io.WriteString(w, "ng")
			}
			io.WriteString(w, "ok")
		}
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/empty", nil)
		handler(w, req)
		assert.Equal(t, "ok", w.Body.String())
	})
}

func TestGetResponseBody(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			io.WriteString(w, "ok")
			return
		}
		io.WriteString(w, "")
	}))
	defer func() { testServer.Close() }()

	t.Run("ok", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, testServer.URL, nil)
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		body, err := GetResponseBody(resp)
		assert.NoError(t, err)
		assert.Equal(t, "ok", string(body))
	})

	t.Run("ng:empty", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, testServer.URL, nil)
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		body, err := GetResponseBody(resp)
		assert.NoError(t, err)
		assert.Equal(t, "", string(body))
	})
}

type TestSetFormValueToStructStruct struct {
	Value1 string `json:"test1"`
	Value2 string `json:"test2"`
}

func TestSetFormValueToStruct(t *testing.T) {
	t.Run("ng:empty", func(t *testing.T) {
		formValues := url.Values{}

		got := &TestSetFormValueToStructStruct{}
		err := SetFormValueToStruct(formValues, got)
		assert.NoError(t, err)
		assert.Equal(t, got, &TestSetFormValueToStructStruct{})
	})

	t.Run("ng:nil", func(t *testing.T) {
		formValues := url.Values{}
		formValues.Set("test1", "test1_value")
		err := SetFormValueToStruct(formValues, nil)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		formValues := url.Values{}
		formValues.Set("test1", "test1_value")
		formValues.Set("test2", "test2_value")

		got := &TestSetFormValueToStructStruct{
			Value1: "be overwrite",
		}
		err := SetFormValueToStruct(formValues, got)
		assert.NoError(t, err)
		assert.Equal(t, &TestSetFormValueToStructStruct{
			Value1: "test1_value",
			Value2: "test2_value",
		}, got)
	})
}
