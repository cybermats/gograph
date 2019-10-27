package main

import (
	"encoding/json"

	"cybermats/gograph/internal/repository"
)

func getTVDBMetaData(id string) (interface{}, string, error) {
	infoStr, err := repository.GetTvdbInfo(id)
	if err != nil {
		return "", "", err
	}
	var info interface{}
	err = json.Unmarshal(infoStr, &info)
	if err != nil {
		return nil, "", err
	}

	image, err := repository.GetTvdbImage(id)
	if err != nil {
		return nil, "", err
	}
	return info, string(image) + "=s256", nil
}
