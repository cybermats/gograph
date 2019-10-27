package main

import (
	"context"
	"log"
	"os"

	"cybermats/gograph/internal/helper"
	"cybermats/gograph/internal/searcher"

	"cloud.google.com/go/storage"
)

func initDbFromGCS() (*searcher.Db, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	bucketName := "matsf-imdb-data"
	bucket := client.Bucket(bucketName)

	basics, err :=
		bucket.Object("datasets.imdbws.com/title.basics.tsv.gz").NewReader(ctx)
	if err != nil {
		return nil, err
	}
	episodes, err :=
		bucket.Object("datasets.imdbws.com/title.episode.tsv.gz").NewReader(ctx)
	if err != nil {
		return nil, err
	}
	ratings, err :=
		bucket.Object("datasets.imdbws.com/title.ratings.tsv.gz").NewReader(ctx)
	if err != nil {
		return nil, err
	}

	s, err := searcher.NewSearcherFromFiles(basics, episodes, ratings)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func createDatabase(filename string) error {
	log.Println("Loading from GCS...")
	b := helper.NewBenchmark("Loading")
	s, err := initDbFromGCS()
	if err != nil {
		return err
	}
	b.Println()
	helper.PrintMemUsage()
	log.Printf("Saving database to %s...\n", filename)
	b = helper.NewBenchmark("Saving")
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	err = s.Write(f)
	b.Println()
	helper.PrintMemUsage()
	return err
}

func loadDatabase(filename string) (*searcher.Db, error) {
	log.Printf("Loading database from %s...", filename)
	s := searcher.NewSearcher()
	b := helper.NewBenchmark("Loading")
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	err = s.Read(f)
	if err != nil {
		return nil, err
	}
	b.Println()
	helper.PrintMemUsage()
	return s, nil
}
