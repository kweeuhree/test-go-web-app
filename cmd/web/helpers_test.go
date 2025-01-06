package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestServerError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{
			name:     "Server error",
			err:      fmt.Errorf("test error"),
			expected: http.StatusInternalServerError,
		},
		{
			name:     "Empty error",
			err:      fmt.Errorf(""),
			expected: http.StatusInternalServerError,
		},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			app.ServerError(resp, entry.err)

			if resp.Result().StatusCode != entry.expected {
				t.Errorf("Expected %d, but got %d", entry.expected, resp.Code)
			}
		})
	}
}

func TestClientError(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{"400 Bad Request", 400},
		{"404 Not Found", 404},
		{"405 Method Not Allowed", 405},
		{"418 I'm a teapot", 418},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			app.ClientError(resp, entry.status)

			if resp.Result().StatusCode != entry.status {
				t.Errorf("Expected %d, but got %d", entry.status, resp.Code)
			}
		})
	}
}

func TestNotFound(t *testing.T) {
	resp := httptest.NewRecorder()
	app.NotFound(resp)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Expected status Not Found, but got %d", resp.Code)
	}
}

func TestDecodeJSON(t *testing.T) {
	tests := []struct {
		name    string
		reqBody string
		err     error
	}{
		{"Valid JSON payload", `{"test":"test"}`, nil},
		{"Invalid JSON payload", `{test:"test"}`, fmt.Errorf("The receipt is invalid.")},
		{"Invalid JSON payload", `{"test":}`, fmt.Errorf("The receipt is invalid.")},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/home", bytes.NewReader([]byte(entry.reqBody)))
			var test map[string]interface{}
			err := app.DecodeJSON(resp, req, &test)

			if err != nil && entry.err == nil {
				t.Errorf("Expected nil, but got %v", err)
			}

			if err == nil && entry.err != nil {
				t.Errorf("Expected to receive an error, but did not")
			}
		})
	}
}

func TestEncodeJSON(t *testing.T) {
	tests := []struct {
		name         string
		data         interface{}
		status       int
		expectedBody string
	}{
		{"Valid JSON payload", map[string]string{"test": "test"}, http.StatusOK, `{"test":"test"}`},
		{"No JSON payload", nil, http.StatusBadRequest, "null"},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			resp := httptest.NewRecorder()

			err := app.EncodeJSON(resp, entry.status, entry.data)
			if err != nil {
				t.Fatalf("EncodeJSON failed: %v", err)
			}

			header := resp.Header().Get("Content-Type")
			if header != "application/json" {
				t.Errorf("Expected Content-Type: application/json, but got %v", header)
			}

			if resp.Code != entry.status {
				t.Errorf("Expected status code %d, but got %d", entry.status, resp.Code)
			}

			respBody := strings.TrimSpace(resp.Body.String())
			if respBody != entry.expectedBody {
				t.Errorf("Expected body: %v, but got: %v", entry.expectedBody, respBody)
			}
		})
	}
}

func TestGetIdFromParams(t *testing.T) {
	tests := []struct {
		name       string
		paramsId   string
		url        string
		paramKey   string
		paramValue string
		expected   string
	}{
		{
			name:       "Existing id",
			paramsId:   "id",
			url:        "/user/123",
			paramKey:   "id",
			paramValue: "123",
			expected:   "123",
		},
		{
			name:       "No id",
			paramsId:   "id",
			url:        "/user/",
			paramKey:   "id",
			paramValue: "",
			expected:   "",
		},
		{
			name:       "Incorrect paramsId",
			paramsId:   "uuid",
			url:        "/user/123",
			paramKey:   "id",
			paramValue: "123",
			expected:   "",
		},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			// Create a request with the required URL
			r := httptest.NewRequest(http.MethodGet, entry.url, nil)

			// Set up httprouter parameters and inject them into the request context
			params := httprouter.Params{
				httprouter.Param{Key: entry.paramKey, Value: entry.paramValue},
			}
			ctx := context.WithValue(r.Context(), httprouter.ParamsKey, params)
			r = r.WithContext(ctx)

			// Call GetIdFromParams with the request with context
			id := app.GetIdFromParams(r, entry.paramsId)

			if id != entry.expected {
				t.Errorf("Expected %v, received %v", entry.expected, id)
			}
		})
	}

}
