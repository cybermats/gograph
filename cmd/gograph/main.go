package main

import (
	"cybermats/gograph/internal/searcher"
	"cybermats/gograph/internal/version"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
)

type config struct {
	CreateDatabase   bool
	DatabaseFilename string
	WebDirectory     string
}

func (c config) String() string {
	v := reflect.ValueOf(c)
	t := v.Type()
	fields := make([]string, 0, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		name := t.Field(i).Name
		value := v.Field(i).Interface()
		fields = append(fields, fmt.Sprintf("%s: %v", name, value))
	}
	return "{ " + strings.Join(fields, ", ") + " }"
}

func initArgs() config {
	var cfg config

	pflag.BoolVarP(&cfg.CreateDatabase, "create", "c", false,
		"If True a new database file will be created and no service "+
			"will be run. If false a database file will be loaded and "+
			"the service started.")
	pflag.StringVarP(&cfg.DatabaseFilename, "filename", "f", "",
		"Path to the database file.")
	pflag.StringVarP(&cfg.WebDirectory, "web-dir", "w", "web",
		"Directory where the static and template files are located.")
	help := pflag.BoolP("help", "h", false,
		"Show help for all commands.")

	pflag.Parse()

	if cfg.DatabaseFilename == "" || *help {
		fmt.Println("Usage:")
		pflag.PrintDefaults()
		os.Exit(1)
	}
	return cfg
}

func runServer(webDirectory string, s *searcher.Searcher) {
	router := mux.NewRouter().StrictSlash(true)
	log.Println("Initializing site...")
	if err := initSite(router, webDirectory); err != nil {
		log.Fatal("Failed initializing site: ", err)
	}
	log.Println("Initializing API...")
	if err := initAPI(router, s); err != nil {
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

func main() {
	cfg := initArgs()
	log.Println("GoGraph", version.Get())
	log.Println("Config: ", cfg)
	var err error

	if cfg.CreateDatabase {
		err = createDatabase(cfg.DatabaseFilename)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Done.")
		return
	}

	s, err := loadDatabase(cfg.DatabaseFilename)
	if err != nil {
		log.Fatal(err)
	}
	runServer(cfg.WebDirectory, s)
}
