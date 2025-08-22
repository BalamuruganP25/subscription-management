package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"subscription-management/pkg/handler"

	"github.com/go-chi/chi/v5"
)

func UpdateUserById(s *handler.ProcessConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req handler.UpdateUserRequest

		// Get user ID from URL parameters
		userID := chi.URLParam(r, "id")
		if userID == "" {
			handler.ErrorResponse(w, http.StatusBadRequest, handler.ErrResponse{
				Title:   "validation error",
				Details: "user id should be empty",
			})
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handler.ErrorResponse(w, http.StatusBadRequest, handler.ErrResponse{
				Title:   "invalid request",
				Details: fmt.Sprintf("failed to parse request body: %v", err),
			})
			return
		}

		// Fetch user details from the database
		user, err := s.CurdRepo.GetUser(r.Context(), userID)
		if err != nil {
			if userID == "" {
				handler.ErrorResponse(w, http.StatusNotFound,
					handler.ErrResponse{
						Title:   "not found",
						Details: fmt.Sprintf("user not found: %v", err),
					},
				)
				return
			}

			handler.ErrorResponse(w, http.StatusInternalServerError,
				handler.ErrResponse{
					Title:   "internal error",
					Details: fmt.Sprintf("failed to fetch user: %v", err),
				},
			)
			return
		}

		err = s.CurdRepo.UpdateUser(r.Context(), userID, req.Phone_number)
		if err != nil {
			handler.ErrorResponse(w, http.StatusInternalServerError,
				handler.ErrResponse{
					Title:   "update error",
					Details: err.Error(),
				},
			)
			return
		}

		// Send the user details in the response
		handler.SendResponse(w, user, http.StatusOK)
	}
}
