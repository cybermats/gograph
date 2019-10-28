package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"cybermats/gograph/internal/apihelper"
	"cybermats/gograph/internal/helper"
	"cybermats/gograph/internal/searcher"
	"cybermats/gograph/internal/version"

	"cloud.google.com/go/storage"
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

func getUpdated(filename string) (time.Time, error) {
	gsInfo, err := searcher.ParseGCS(filename)
	if err != nil {
		return time.Time{}, err
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return time.Time{}, err
	}

	attr, err := client.Bucket(gsInfo.Bucket).Object(gsInfo.Object).Attrs(ctx)
	if err != nil {
		return time.Time{}, err
	}

	return attr.Updated, nil
}

func isDatabaseNewest(cfg config) (bool, error) {
	du, err := getUpdated(cfg.DatabaseFilename)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return false, nil
		}
		return false, err
	}
	bu, err := getUpdated(cfg.BasicsFilename)
	if err != nil {
		return false, err
	}

	if du.Before(bu) {
		return false, nil
	}

	eu, err := getUpdated(cfg.EpisodesFilename)
	if err != nil {
		return false, err
	}
	if du.Before(eu) {
		return false, nil
	}

	ru, err := getUpdated(cfg.RatingsFilename)
	if err != nil {
		return false, err
	}
	if du.Before(ru) {
		return false, nil
	}

	return true, nil
}

func updateHandlerFunc(cfg config) func(http.ResponseWriter,
	*http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			apihelper.Write404(w)
			return
		}
		upToDate, err := isDatabaseNewest(cfg)
		if err != nil {
			apihelper.Write500(w, err)
			log.Print(err)
			return
		}
		if upToDate {
			apihelper.WriteOKMessage(w, "No need to update")
			return
		}
		b := helper.NewBenchmark("Database")
		err = searcher.CreateDatabase(
			cfg.BasicsFilename,
			cfg.EpisodesFilename,
			cfg.RatingsFilename,
			cfg.DatabaseFilename)
		d := b.Duration()
		if err != nil {
			apihelper.Write500(w, err)
			return
		}
		apihelper.WriteOKMessage(w, fmt.Sprintf("Update took %v", d))
	}

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	apihelper.WriteOKMessage(w, "up and running")
}

func main() {
	cfg := initArgs()
	log.Println("createDb", version.Get())
	log.Println("Config: ", cfg)

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/create", updateHandlerFunc(cfg))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
