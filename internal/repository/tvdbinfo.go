package repository

import (
	"context"

	"cloud.google.com/go/datastore"
)

func GetTvdbInfo(imdbID string) ([]byte, error) {
	ctx := context.Background()
	projectID := "matsf-cloud-playpen"

	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	k := datastore.NameKey("TvDbInfo", imdbID, nil)
	e := new(tvdbInfo)
	if err = client.Get(ctx, k, e); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, nil
		}
		return nil, err
	}

	return e.Data, nil
}

func GetTvdbImage(imdbID string) (string, error) {
	ctx := context.Background()
	projectID := "matsf-cloud-playpen"

	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return "", err
	}

	k := datastore.NameKey("TvDbImage", imdbID, nil)
	e := new(tvdbImage)
	if err = client.Get(ctx, k, e); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return "", nil
		}
		return "", err
	}

	return e.Path, nil
}
