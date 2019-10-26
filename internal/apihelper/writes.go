package apihelper

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorPayload struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
}

func WriteJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("Error: ", err)
	}
}

func WriteJSONFromText(w http.ResponseWriter, data []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	_, err := w.Write(data)
	if err != nil {
		log.Println("Error: ", err)
	}
}

func WriteOK(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, data, http.StatusOK)
}

func WriteError(w http.ResponseWriter, message string, statusCode int) {
	data := errorPayload{message, statusCode}
	WriteJSON(w, data, statusCode)
}

func Write404(w http.ResponseWriter) {
	WriteError(w, "item not found", http.StatusNotFound)
}

func Write500(w http.ResponseWriter, err error) {
	WriteError(w, err.Error(), http.StatusInternalServerError)
}
