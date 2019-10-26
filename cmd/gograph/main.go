package main

import (
	"log"
	"net/http"
)

func main() {
	if err := initAPI(); err != nil {
		log.Fatal("Failed initializing API: ", err)
	}
	if err := initSite(); err != nil {
		log.Fatal("Failed initializing site: ", err)
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/favicon.ico", fs)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
