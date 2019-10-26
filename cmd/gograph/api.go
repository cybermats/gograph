package main

import (
	"encoding/json"
	"net/http"
	"regexp"

	"google.golang.org/appengine"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
}

func infoHandler(w http.ResponseWriter, r *http.Request, id string) {
	result := []byte("{}")

	if id != "" {
		info, err := getTvdbInfo(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if info != nil {
			result = info
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func imageHandler(w http.ResponseWriter, r *http.Request, id string) {
	result := make(map[string]string)

	if id != "" {
		info, err := getTvdbImage(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if info != "" {
			result["url"] = info + "=s256"
		}
	}

	js, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func makeAPIHandler(fn func(http.ResponseWriter, *http.Request, string),
	idRegexp *regexp.Regexp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := idRegexp.FindStringSubmatch(r.URL.Path)
		fn(w, r, m[1])
	}
}

func initAPI() error {
	idRegexp := regexp.MustCompile("([^/]*)$")
	http.HandleFunc("/api/info/", makeAPIHandler(infoHandler, idRegexp))
	http.HandleFunc("/api/image/", makeAPIHandler(imageHandler, idRegexp))
	return nil
}
