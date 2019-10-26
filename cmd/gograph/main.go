package main

import (
	"cybermats/gograph/internal/version"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("GoGraph", version.Get())
	env := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	log.Println("Env: ", env)

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
