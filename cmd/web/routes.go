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

	// static assets
	fileServer := http.FileServer(http.Dir("./static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// register routes
	router.Get("/", app.Home)

	// chi.Mux satisfies http.Handler type
	return router
}
