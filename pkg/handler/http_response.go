package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

// SendResponse is a utility function to send the response to the client
func SendResponse(w http.ResponseWriter, response any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
