package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/julienschmidt/httprouter"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) ServerError(w http.ResponseWriter, err error) {
	// Use the debug.Stack() function to get a stack trace for the current goroutine and append it to the log message
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// Report the file name and line number one step back in the stack trace
	// to have a clearer idea of where the error actually originated from
	// set frame depth to 2
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user
func (app *application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Not found helper
func (app *application) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}

// Decode the JSON body of a request into the destination struct
func (app *application) DecodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return err
	}
	return nil
}

// Encodes provided data into a JSON response
func (app *application) EncodeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// Get parameter id from the request
func (app *application) GetIdFromParams(r *http.Request, paramsId string) string {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName(paramsId)
	return id
}
