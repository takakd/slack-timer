package httputil

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

// Get body as string from http.request.
// Ref: https://developer.twitter.com/en/docs/basics/response-codes
func GetRequestBody(req *http.Request) ([]byte, error) {
	if req == nil {
		return nil, errors.New("arguments should not be nil.")
	}

	body := make([]byte, req.ContentLength)
	if req.ContentLength > 0 {
		req.Body.Read(body)
	}
	return body, nil
}

// Get body as string from http.Response.
// More about details are https://golang.org/pkg/net/http/
func GetResponseBody(resp *http.Response) ([]byte, error) {
	if resp == nil {
		return nil, errors.New("arguments should not be nil.")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

type ErrorJsonResponse struct {
	Summary   string `json:"summary"`
	ErrorCode string `json:"error_code"`
	Detail    string `json:"detail"`
}

// Create error response in json format.
func NewErrorJsonResponse(summary string, errorCode string, detail string) ([]byte, error) {
	response := &ErrorJsonResponse{
		Summary:   summary,
		ErrorCode: errorCode,
		Detail:    detail,
	}
	return json.Marshal(response)
}

// Write error response data to http.ResponseWriter.
func WriteErrorJsonResponse(w http.ResponseWriter, headers map[string]string, httpStatusCode int, summary, errorCode, detail string) error {
	body, err := NewErrorJsonResponse(summary, errorCode, detail)
	if err != nil {
		return err
	}
	return WriteJsonResponse(w, headers, httpStatusCode, body)
}

// Write response data to http.ResponseWriter.
func WriteJsonResponse(w http.ResponseWriter, headers map[string]string, httpStatusCode int, body []byte) error {
	w.Header().Set("Content-Type", "application/json")
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(httpStatusCode)
	w.Write(body)
	return nil
}

// Set FormValues to struct
func SetFormValueToStruct(values url.Values, structPtr interface{}) error {
	if len(values) == 0 {
		// No values
		return nil
	}

	if structPtr == nil {
		return errors.New("http/SetFormValueToStruct: structPtr is nil")
	}

	// Get struct pointer
	ptr := reflect.ValueOf(structPtr)
	if ptr.Type().Kind() != reflect.Ptr {
		return errors.New("http/SetFormValueToStruct: structPtr is not pointer")
	}

	// Get struct value
	value := ptr.Elem()
	if value.Type().Kind() != reflect.Struct {
		return errors.New("http/SetFormValueToStruct: structPtr is not struct pointer")
	}

	// Set value to struct field
	valueType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		jsonTag := valueType.Field(i).Tag.Get("json")
		if value.Field(i).CanSet() == false {
			return errors.New("cannot set value to field")
		}
		value.Field(i).Set(reflect.ValueOf(values.Get(jsonTag)))
	}

	return nil
}
