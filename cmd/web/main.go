package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	errorLog := log.New(os.Stdout, "ERROR\t", log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.LstdFlags)
	application := application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// start the web server
	port := ":1313"
	server := &http.Server{
		Addr:     port,
		ErrorLog: errorLog,
		Handler:  application.routes(),
	}
	application.infoLog.Println("Server started at port", port)
	err := server.ListenAndServe()
	if err != nil {
		application.errorLog.Fatal(err)
	}
}
