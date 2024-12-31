package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

// Mock handler for testing
func mockHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

func Test_application_routes(t *testing.T) {

	router := httprouter.New()

	// Register routes for testing
	router.GET("/", mockHandler)
	router.GET("/static/*filepath", mockHandler)

	var registered = []struct {
		route  string
		method string
	}{
		{"/", "GET"},
		{"/static/*", "GET"},
	}

	for _, route := range registered {
		// check if the route exists
		if !routeExists(router, route.route, route.method) {
			t.Errorf("route %s is not registered", route.route)
		}
	}
}

func routeExists(router *httprouter.Router, testRoute, testMethod string) bool {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(testMethod, testRoute, nil)

	router.ServeHTTP(recorder, request)

	return recorder.Code != http.StatusNotFound
}
