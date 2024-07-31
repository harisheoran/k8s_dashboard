package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	// create a router
	router := http.NewServeMux()

	// router map the path with the handler
	router.HandleFunc("/", app.Roothandler)

	return router
}
