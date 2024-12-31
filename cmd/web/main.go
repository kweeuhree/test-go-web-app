package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// Error and info logs
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// set up an application config
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// HTTP server config
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	// print out a message
	app.infoLog.Printf("Starting server on port %v...", *addr)
	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
