package main

import (
	"cybermats/gograph/internal/version"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
)

type config struct {
	createDatabase   *bool
	databaseFilename *string
	webDirectory     *string
}

func initArgs() config {
	var cfg config
	cfg.createDatabase = pflag.BoolP("create", "c", false,
		"If True a new database file will be created and no service "+
			"will be run. If false a database file will be loaded and "+
			"the service started.")
	cfg.databaseFilename = pflag.StringP("filename", "f", "",
		"Path to the database file.")
	cfg.webDirectory = pflag.StringP("web-dir", "w", "web",
		"Directory where the static and template files are located.")
	help := pflag.BoolP("help", "h", false,
		"Show help for all commands.")

	pflag.Parse()

	if *cfg.databaseFilename == "" || *help {
		fmt.Println("Usage:")
		pflag.PrintDefaults()
		os.Exit(1)
	}
	return cfg
}

func main() {
	cfg := initArgs()
	log.Println("GoGraph", version.Get())
	log.Println("Config: ", cfg)

	router := mux.NewRouter().StrictSlash(true)

	log.Println("Initializing site...")
	if err := initSite(router, *cfg.webDirectory); err != nil {
		log.Fatal("Failed initializing site: ", err)
	}
	log.Println("Initializing API...")
	if err := initAPI(router); err != nil {
		log.Fatal("Failed initializing API: ", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("No PORT set in env vars. Using default.")
		port = "8080"
	}
	log.Printf("Starting service on port %s...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
