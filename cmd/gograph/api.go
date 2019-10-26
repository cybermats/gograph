package main

import (
	"encoding/json"
	"log"
	"net/http"

	"cybermats/gograph/internal/repository"

	"github.com/gorilla/mux"
)

type errorPayload struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
}

type imagePayload struct {
	URL string `json:"url"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	//	ctx := appengine.NewContext(r)
}

func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("Error: ", err)
	}
}

func writeJSONFromText(w http.ResponseWriter, data []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	_, err := w.Write(data)
	if err != nil {
		log.Println("Error: ", err)
	}
}

func writeOK(w http.ResponseWriter, data interface{}) {
	writeJSON(w, data, http.StatusOK)
}

func writeError(w http.ResponseWriter, message string, statusCode int) {
	data := errorPayload{message, statusCode}
	writeJSON(w, data, statusCode)
}

func write404(w http.ResponseWriter) {
	writeError(w, "item not found", http.StatusNotFound)
}

func write500(w http.ResponseWriter, err error) {
	writeError(w, err.Error(), http.StatusInternalServerError)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id != "" {
		info, err := repository.GetTvdbInfo(id)
		if err != nil {
			write500(w, err)
			return
		}
		if info != nil {
			// The result is a serialized json string
			writeJSONFromText(w, info, http.StatusOK)
			return
		}
	}
	write404(w)
	return
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id != "" {
		info, err := repository.GetTvdbImage(id)
		if err != nil {
			write500(w, err)
			return
		}
		if info != "" {
			writeOK(w, imagePayload{info + "=s256"})
			return
		}
	}
	write404(w)
	return
}

func emptyHandler(w http.ResponseWriter, r *http.Request) {
	write404(w)
}

func initAPI(router *mux.Router) error {
	router.HandleFunc("/api/info/", emptyHandler)
	router.HandleFunc("/api/info/{id}", infoHandler)
	router.HandleFunc("/api/image/", emptyHandler)
	router.HandleFunc("/api/image/{id}", imageHandler)
	return nil
}
