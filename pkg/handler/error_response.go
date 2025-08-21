package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrResponse struct {
	Title   string `json:"title"`
	Details string `json:"details"`
}

func ErrorResponse(w http.ResponseWriter, status int, ErrRes ErrResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(ErrRes); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}
