package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter" // router
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// static assets
	fileServer := http.FileServer(http.Dir("./static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// register routes
	router.Handler(http.MethodGet, "/", http.HandlerFunc(app.Home))

	// chi.Mux satisfies http.Handler type
	// Chain middleware
	standard := alice.New(app.recoverPanic, app.logRequest)

	// Return the 'standard' middleware chain
	return standard.Then(router)
}
