package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"cybermats/gograph/internal/searcher"
	"cybermats/gograph/internal/version"

	"github.com/spf13/pflag"
)

type config struct {
	BasicsFilename   string
	EpisodesFilename string
	RatingsFilename  string
	DatabaseFilename string
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

	pflag.StringVarP(&cfg.DatabaseFilename, "database", "d", "", "Path to the database file.")
	pflag.StringVarP(&cfg.BasicsFilename, "basics", "",
		"gs://matsf-imdb-data/datasets.imdbws.com/title.basics.tsv.gz",
		"Path to the basics file.")
	pflag.StringVarP(&cfg.EpisodesFilename, "episodes", "",
		"gs://matsf-imdb-data/datasets.imdbws.com/title.episode.tsv.gz",
		"Path to the episodes file.")
	pflag.StringVarP(&cfg.RatingsFilename, "ratings", "",
		"gs://matsf-imdb-data/datasets.imdbws.com/title.ratings.tsv.gz",
		"Path to the ratings file.")
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

func main() {
	cfg := initArgs()
	log.Println("createDb", version.Get())
	log.Println("Config: ", cfg)

	err := searcher.CreateDatabase(
		cfg.BasicsFilename,
		cfg.EpisodesFilename,
		cfg.RatingsFilename,
		cfg.DatabaseFilename)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done.")
}
