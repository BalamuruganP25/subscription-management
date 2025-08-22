package user

import (
	"fmt"
	"net/http"
	"subscription-management/pkg/handler"

	"github.com/go-chi/chi/v5"
)

func GetUserById(s *handler.ProcessConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from URL parameters
		userID := chi.URLParam(r, "id")
		if userID == "" {
			handler.ErrorResponse(w, http.StatusBadRequest, handler.ErrResponse{
				Title:   "validation error",
				Details: "user id should be empty",
			})
			return
		}

		// Fetch user details from the database
		user, err := s.CurdRepo.GetUser(r.Context(), userID)
		if err != nil {
			handler.ErrorResponse(w, http.StatusNotFound,
				handler.ErrResponse{
					Title:   "not found",
					Details: fmt.Sprintf("user not found: %v", err),
				},
			)
			return
		}

		// Send the user details in the response
		handler.SendResponse(w, user, http.StatusOK)
	}
}
