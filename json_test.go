package server

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testStruct struct {
	Name string  `json:"name"`
	Description string `json:"description"`
}
func TestParseJsonRequestStruct(t *testing.T) {
	t.Run("valid json", func(t *testing.T) {
		result := &testStruct{}
		r := httptest.NewRequest("POST","/",bytes.NewBufferString(`{"name":"bob","description":"friend"}`))
		r.Header.Add("Content-Type","application/json")
		err := ParseJsonRequestStruct(r,result)
		assert.Nil(t, err)
		assert.Equal(t, &testStruct{
			Name:        "bob",
			Description: "friend",
		},result)
	})
	t.Run("valid json to map", func(t *testing.T) {
		result := make(map[string]interface{})
		r := httptest.NewRequest("POST","/",bytes.NewBufferString(`{"name":"bob","description":"friend"}`))
		r.Header.Add("Content-Type","application/json")
		err := ParseJsonRequestStruct(r,result)
		assert.NotNil(t, err)
		assert.Equal(t,ErrPointerError,err)
	})
	t.Run("invalid json", func(t *testing.T) {
		result := &testStruct{}
		r := httptest.NewRequest("POST","/",bytes.NewBufferString(`"name":"bob","description":"friend"}`))
		r.Header.Add("Content-Type","application/json")
		err := ParseJsonRequestStruct(r,result)
		assert.NotNil(t, err)
		assert.Equal(t, err,ErrInvalidJson)
	})
	t.Run("invalid content-type", func(t *testing.T) {
		result := &testStruct{}
		r := httptest.NewRequest("POST","/",bytes.NewBufferString(`{"name":"bob","description":"friend"}`))
		r.Header.Add("Content-Type","application/xml")
		err := ParseJsonRequestStruct(r,result)
		assert.NotNil(t, err)
		assert.Equal(t, err,ErrContentType)
	})
	t.Run("non-pointer", func(t *testing.T) {
		result := testStruct{}
		r := httptest.NewRequest("POST","/",bytes.NewBufferString(`{"name":"bob","description":"friend"}`))
		r.Header.Add("Content-Type","application/json")
		err := ParseJsonRequestStruct(r,result)
		assert.NotNil(t, err)
		assert.Equal(t, err,ErrPointerError)
		assert.NotEqual(t, testStruct{
			Name:        "bob",
			Description: "friend",
		},result)
	})
}

func TestParseJsonRequestMap(t *testing.T) {
	t.Run("valid json to map", func(t *testing.T) {
		r := httptest.NewRequest("POST","/",bytes.NewBufferString(`{"name":"bob","description":"friend"}`))
		r.Header.Add("Content-Type","application/json")
		m, err := ParseJsonRequestMap(r)
		assert.Nil(t, err)
		assert.Equal(t,map[string]interface{}{
			"name":"bob",
			"description":"friend",
		},m)
	})
	t.Run("invalid json to map", func(t *testing.T) {
		r := httptest.NewRequest("POST","/",bytes.NewBufferString(`"name":"bob","description":"friend"}`))
		r.Header.Add("Content-Type","application/json")
		_, err := ParseJsonRequestMap(r)
		assert.NotNil(t, err)
		assert.Equal(t,ErrInvalidJson,err)
	})
	t.Run("invalid content type", func(t *testing.T) {
		r := httptest.NewRequest("POST","/",bytes.NewBufferString(`"{name":"bob","description":"friend"}`))
		r.Header.Add("Content-Type","application/xml")
		_, err := ParseJsonRequestMap(r)
		assert.NotNil(t, err)
		assert.Equal(t,ErrContentType,err)
	})
}

func TestWriteStructJson(t *testing.T) {
	t.Run("valid struct", func(t *testing.T) {
		test := testStruct{
			Name:        "test",
			Description: "test",
		}
		w := httptest.NewRecorder()
		err := WriteStructJson(200,test,w)
		assert.Nil(t, err)
		assert.Equal(t, "application/json",w.Header().Get("Content-Type"))
		assert.Equal(t, `{"name":"test","description":"test"}`,w.Body.String())
		assert.Equal(t, http.StatusOK,w.Code)
	})
	t.Run("valid struct, change status", func(t *testing.T) {
		test := testStruct{
			Name:        "test",
			Description: "test",
		}
		w := httptest.NewRecorder()
		err := WriteStructJson(http.StatusAccepted,test,w)
		assert.Nil(t, err)
		assert.Equal(t, "application/json",w.Header().Get("Content-Type"))
		assert.Equal(t, `{"name":"test","description":"test"}`,w.Body.String())
		assert.Equal(t, http.StatusAccepted,w.Code)
	})
	t.Run("not a struct ", func(t *testing.T) {
		test := map[string]interface{}{
			"name": "test",
			"description": "test",
		}
		w := httptest.NewRecorder()
		err := WriteStructJson(200,test,w)
		assert.NotNil(t, err)
		assert.Equal(t, ErrNotStruct,err)
	})
}

func TestWriteMapJson(t *testing.T) {
	t.Run("valid map", func(t *testing.T) {
		test := map[string]interface{}{
			"name": "test",
			"description": "test",
		}
		w := httptest.NewRecorder()
		err := WriteMapJson(200,test,w)
		assert.Nil(t, err)
		assert.Equal(t, "application/json",w.Header().Get("Content-Type"))
		assert.Equal(t, `{"description":"test","name":"test"}`,w.Body.String())
		assert.Equal(t, http.StatusOK,w.Code)
	})
	t.Run("valid map, change status", func(t *testing.T) {
		test := map[string]interface{}{
			"name": "test",
			"description": "test",
		}
		w := httptest.NewRecorder()
		err := WriteMapJson(http.StatusAccepted,test,w)
		assert.Nil(t, err)
		assert.Equal(t, "application/json",w.Header().Get("Content-Type"))
		assert.Equal(t, `{"description":"test","name":"test"}`,w.Body.String())
		assert.Equal(t, http.StatusAccepted,w.Code)
	})
	t.Run("valid  multi-level map", func(t *testing.T) {
		test := map[string]interface{}{
			"name": "test",
			"description": map[string]interface{}{
				"test": "test",
			},
		}
		w := httptest.NewRecorder()
		err := WriteMapJson(200,test,w)
		assert.Nil(t, err)
		assert.Equal(t, "application/json",w.Header().Get("Content-Type"))
		assert.Equal(t, `{"description":{"test":"test"},"name":"test"}`,w.Body.String())
		assert.Equal(t, http.StatusOK,w.Code)
	})
}
