package searcher

import (
	"compress/gzip"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"cybermats/gograph/internal/helper"

	"cloud.google.com/go/storage"
)

func initDbFromGCS(basics, episodes, ratings string) (*Db, error) {
	log.Println("Loading basics from:", basics)
	basicsReader, err := newReader(basics)
	if err != nil {
		return nil, err
	}
	log.Println("Loading episodes from:", episodes)
	episodesReader, err := newReader(episodes)
	if err != nil {
		return nil, err
	}
	log.Println("Loading ratings from:", ratings)
	ratingsReader, err := newReader(ratings)
	if err != nil {
		return nil, err
	}

	return NewDbFromFiles(basicsReader, episodesReader, ratingsReader)
}

func newReader(filename string) (io.ReadCloser, error) {
	gspath, err := ParseGCS(filename)
	if err == PathError {
		log.Println("File is local.")
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		return f, nil
	}

	log.Println("File is in GCS.")
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.Bucket(gspath.Bucket).Object(gspath.Object).NewReader(ctx)
}

func newWriter(filename string) (io.WriteCloser, error) {
	gspath, err := ParseGCS(filename)
	if err == PathError {
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		return f, nil
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	w := client.Bucket(gspath.Bucket).Object(gspath.Object).NewWriter(ctx)
	return w, nil

}

// CreateDatabase initializes a database based on the supplied files and then stores it.
// Files can reside locally or in Google Cloud Storage.
func CreateDatabase(basics, episodes, ratings, database string) error {
	log.Println("Loading from GCS...")
	b := helper.NewBenchmark("Loading")
	s, err := initDbFromGCS(basics, episodes, ratings)
	if err != nil {
		return err
	}
	b.Println()
	helper.PrintMemUsage()
	log.Printf("Saving database to %s...\n", database)
	b = helper.NewBenchmark("Saving")
	w, err := newWriter(database)
	if err != nil {
		return err
	}
	zw, err := gzip.NewWriterLevel(w, gzip.BestCompression)
	if err != nil {
		return err
	}

	zw.Name = filepath.Base(database)
	zw.Comment = "Database for graph tv"

	err = s.Write(zw)
	if err != nil {
		return err
	}
	b.Println()
	helper.PrintMemUsage()

	err = zw.Close()
	if err != nil {
		return err
	}
	return w.Close()
}

// LoadDatabase loads a database from a file.
// File can reside locally or in Google Cloud Storage.
func LoadDatabase(database string) (*Db, error) {
	log.Printf("Loading database from %s...", database)
	s := NewDb()
	b := helper.NewBenchmark("Loading")

	r, err := newReader(database)
	if err != nil {
		return nil, err
	}

	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	err = s.Read(zr)
	if err != nil {
		return nil, err
	}
	err = zr.Close()
	if err != nil {
		return nil, err
	}
	err = r.Close()
	if err != nil {
		return nil, err
	}

	b.Println()
	helper.PrintMemUsage()
	return s, nil
}
