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

type messagePayload struct {
	Message string `json:"message"`
}

// WriteJSON helps writing headers and marshalling objects for a JSON response
func WriteJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("Error: ", err)
	}
}

// WriteJSONFromText creates a JSON response based on a serialized JSON data.
func WriteJSONFromText(w http.ResponseWriter, data []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	_, err := w.Write(data)
	if err != nil {
		log.Println("Error: ", err)
	}
}

// WriteOK creates a JSON response with StatusOK.
func WriteOK(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, data, http.StatusOK)
}

// WriteOKMessage creates a JSON response with StatusOK.
func WriteOKMessage(w http.ResponseWriter, message string) {
	data := messagePayload{message}
	WriteJSON(w, data, http.StatusOK)
}

// WriteError wraps an error message in a general error JSON message.
func WriteError(w http.ResponseWriter, message string, statusCode int) {
	data := errorPayload{message, statusCode}
	WriteJSON(w, data, statusCode)
}

// Write404 creates a simple 404 JSON message.
func Write404(w http.ResponseWriter) {
	WriteError(w, "item not found", http.StatusNotFound)
}

// Write500 creates a simple 500 message.
func Write500(w http.ResponseWriter, err error) {
	WriteError(w, err.Error(), http.StatusInternalServerError)
}
