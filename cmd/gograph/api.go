package main

import (
	"log"
	"net/http"

	"cybermats/gograph/internal/apihelper"
	"cybermats/gograph/internal/repository"
	"cybermats/gograph/internal/searcher"

	"github.com/gorilla/mux"
)

type imagePayload struct {
	URL string `json:"url"`
}

func searchHandler(w http.ResponseWriter, r *http.Request, s *searcher.Db) {
	filter := r.FormValue("filter")
	log.Println("Search: ", filter)
	titles := s.Search(filter)
	apihelper.WriteOK(w, titles)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id != "" {
		info, err := repository.GetTvdbInfo(id)
		if err != nil {
			apihelper.Write500(w, err)
			return
		}
		if info != nil {
			// The result is a serialized json string
			apihelper.WriteJSONFromText(w, info, http.StatusOK)
			return
		}
	}
	apihelper.Write404(w)
	return
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id != "" {
		info, err := repository.GetTvdbImage(id)
		if err != nil {
			apihelper.Write500(w, err)
			return
		}
		if info != "" {
			apihelper.WriteOK(w, imagePayload{info + "=s256"})
			return
		}
	}
	apihelper.Write404(w)
	return
}

func emptyHandler(w http.ResponseWriter, r *http.Request) {
	apihelper.Write404(w)
}

func makeSearcherHandler(
	fn func(http.ResponseWriter, *http.Request, *searcher.Db),
	s *searcher.Db) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, s)
	}
}

func initAPI(router *mux.Router, s *searcher.Db) error {
	sr := router.PathPrefix("/api").Subrouter()
	sr.HandleFunc("/info/", emptyHandler)
	sr.HandleFunc("/info/{id}", infoHandler)
	sr.HandleFunc("/image/", emptyHandler)
	sr.HandleFunc("/image/{id}", imageHandler)
	sr.HandleFunc("/search/", makeSearcherHandler(searchHandler, s))
	sr.NotFoundHandler = http.HandlerFunc(emptyHandler)
	return nil
}
