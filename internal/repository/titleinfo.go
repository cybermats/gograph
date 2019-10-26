package repository

import (
	"context"
	"log"
	"time"

	"cybermats/gograph/internal/helper"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

func GetTop(days int, count int) ([]TitleInfo, error) {
	ctx := context.Background()
	projectID := "matsf-cloud-playpen"

	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal("Failed to create client: ", err)
		return nil, err
	}

	cutoff := time.Now().AddDate(0, 0, -days)
	log.Println("Cutoff: ", cutoff)

	query := datastore.NewQuery("TitleView").Filter("datetime >=", cutoff)
	it := client.Run(ctx, query)

	titleMap := make(map[string]int)
	for {
		var item titleView
		_, err := it.Next(&item)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal("Error executing query: ", err)
			return nil, err
		}
		titleMap[item.TID]++
	}
	sortedTitles := helper.SortMapByValue(titleMap)

	keys := make([]*datastore.Key, 0, count)
	for i, p := range sortedTitles {
		if i == count {
			break
		}
		keys = append(keys, datastore.NameKey("ImdbInfo", p.Key, nil))
	}

	imdbInfo := make([]imdbInfo, len(keys))

	err = client.GetMulti(ctx, keys, imdbInfo)
	if err != nil {
		log.Fatal("Failed to get title info: ", err)
		return nil, err
	}

	output := make([]TitleInfo, len(imdbInfo))

	for i, key := range keys {
		id := key.Name
		output[i].ID = id
		output[i].Title = imdbInfo[i].Title
		output[i].Year = imdbInfo[i].Year
		output[i].Count = titleMap[id]
	}
	log.Println(output)

	return output, nil
}
