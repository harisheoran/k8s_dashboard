package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	// create a router
	router := http.NewServeMux()

	fileserver := http.FileServer(http.Dir("./ui"))
	router.Handle("/static/", http.StripPrefix("/static/", fileserver))

	router.HandleFunc("/", roothandler)

	fmt.Println("Starting the web server")
	// start the web server
	err := http.ListenAndServe(":1313", router)

	if err != nil {
		log.Fatal(err)
	}
}

func roothandler(w http.ResponseWriter, request *http.Request) {
	templates := []string{
		"./ui/footer.partial.tmpl",
		"./ui/header.partial.tmpl",
		"./ui/home.page.tmpl",
	}

	template, err := template.ParseFiles(templates...)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = template.Execute(w, request)
	if err != nil {
		log.Fatal(err.Error())
	}
}
