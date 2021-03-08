package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

var (
	ErrContentType = fmt.Errorf("content-type header is not set to application/json")
	ErrPointerError = fmt.Errorf("encountered non pointer value,interface passed to ParseJsonRequest should be a pointer")
	ErrInvalidJson = fmt.Errorf("payload is not valid json")
	ErrNotStruct = fmt.Errorf("data is not a struct")
	ErrNotMap    = fmt.Errorf("data is not a map")
)

// ParseJsonRequest parses the body of an http.Request and unmarshals it's json payload into a struct to be used later
// interface needs to be a pointer to avoid errors. It also expects the content-type header to be set to application/json
func ParseJsonRequest(r *http.Request, i interface{}) error {
	if strings.ToLower(r.Header.Get("Content-Type")) != "application/json" {
		return ErrContentType
	}

	if reflect.ValueOf(i).Kind() != reflect.Ptr {
		return ErrPointerError
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	err = json.Unmarshal(body, i)
	if err != nil {
		return ErrInvalidJson
	}
	return nil
}

// WriteStructJson takes a struct and marshals it into a []byte that gets sent to the user as a json payload with appropriate headers
// Also returns provided status code
func WriteStructJson(statusCode int,data interface{}, w http.ResponseWriter) error {
	if reflect.ValueOf(data).Kind() != reflect.Struct {
		return ErrNotStruct
	}
	err := writeJson(statusCode,data,w)
	if err != nil {
		return err
	}
	return nil
}

// WriteMapJson takes a map[string]interface and marshals it into a []byte that gets sent to the user as a json payload with appropriate headers
// Also returns provided status code
func WriteMapJson(statusCode int,data map[string]interface{}, w http.ResponseWriter) error {
	err := writeJson(statusCode,data,w)
	if err != nil {
		return err
	}
	return nil
}

func writeJson(statusCode int,data interface{}, w http.ResponseWriter)error {
	jData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(jData)
	if err != nil {
		return err
	}
	return nil
}