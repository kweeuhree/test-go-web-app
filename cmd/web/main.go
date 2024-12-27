package main

import (
	"log"
	"net/http"
)

// Application-wide dependencies
type applicaton struct {
}

func main() {
	// set up an application config
	app := applicaton{}
	// get application routes
	router := app.routes()
	// print out a message
	log.Println("Starting server on port :8080 ...")
	// start the server
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
