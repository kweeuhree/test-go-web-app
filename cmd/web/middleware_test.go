package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Declare and initialize application instance for all tests
var app *application
var logBuffer bytes.Buffer

func TestMain(m *testing.M) {
	app = &application{
		infoLog:  log.New(&logBuffer, "", log.LstdFlags),
		errorLog: log.New(&logBuffer, "", log.LstdFlags),
	}
	m.Run()
}

func testHandler(testPanic bool) http.HandlerFunc {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if testPanic {
			panic("Test panic triggered")
		} else {
			w.WriteHeader(http.StatusOK)
		}

	}
	return http.HandlerFunc(fn)
}

func Test_logRequest(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		url      string
		expected []string
	}{
		{
			name:     "GET",
			method:   http.MethodGet,
			url:      "/",
			expected: []string{"HTTP/1", "GET", "/"},
		},
		{
			name:     "POST",
			method:   http.MethodPost,
			url:      "/home",
			expected: []string{"HTTP/1", "POST", "/home"},
		},
		{
			name:     "Invalid URL",
			method:   http.MethodGet,
			url:      "/hello-world",
			expected: []string{"HTTP/1", "GET", "/hello-world"},
		},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			// Create test HTTP server with the handler that does not trigger panic
			middleware := app.logRequest(testHandler(false))
			ts := httptest.NewServer(middleware)
			defer ts.Close()

			// Trigger a request
			resp, err := http.NewRequest(entry.method, fmt.Sprintf("%s%s", ts.URL, entry.url), nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			// Send the request and return the response
			http.DefaultClient.Do(resp)

			// Validate log output
			logOutput := logBuffer.String()
			if logOutput == "" {
				t.Error("No information logged to the logger")
			}

			// Check the structure of the log entry
			for _, part := range entry.expected {
				if !strings.Contains(logOutput, part) {
					t.Errorf("Expected log output to contain '%s', but it didn't. Log: %s", part, logOutput)
				}
			}
		})
	}
}

func Test_recoverPanic(t *testing.T) {
	tests := []struct {
		name             string
		panicOccurred    bool
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:             "No panic",
			panicOccurred:    false,
			expectedStatus:   http.StatusOK,
			expectedResponse: "OK",
		},
		{
			name:             "Panic occurred",
			panicOccurred:    true,
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "Internal Server Error",
		},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			middleware := app.recoverPanic(testHandler(entry.panicOccurred))

			// Create a new HTTP request
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Create a ResponseRecorder to capture the response
			resp := httptest.NewRecorder()

			// Serve the HTTP request
			middleware.ServeHTTP(resp, req)

			if resp.Code != entry.expectedStatus {
				t.Errorf("Expected %d, got '%d'", entry.expectedStatus, resp.Code)
			}
			// Check if the "Connection" header is set to "close"
			if entry.panicOccurred {
				if connectionHeader := resp.Header().Get("Connection"); connectionHeader != "close" {
					t.Errorf("Expected 'Connection' header to be 'close', got '%s'", connectionHeader)
				}
			}
		})
	}
}
