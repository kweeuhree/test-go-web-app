package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *applicaton) routes() http.Handler {
	router := chi.NewRouter()

	// register middleware
	router.Use(middleware.Recoverer)

	// register routes
	router.Get("/", app.Home)
	// static assets

	// chi.Mux satisfies http.Handler type
	return router
}
