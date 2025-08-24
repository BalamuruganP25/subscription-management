package user

import (
	"fmt"
	"net/http"
	"subscription-management/pkg/handler"

	"github.com/go-chi/chi/v5"
)

// soft delete user by setting status to false
func DeleteUserById(s *handler.ProcessConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get user ID from URL parameters
		userID := chi.URLParam(r, "id")
		if userID == "" {
			handler.ErrorResponse(w, http.StatusBadRequest, handler.ErrResponse{
				Title:   "invalid request",
				Details: "user id should not be empty",
			})
			return
		}

		// Fetch user details from the database
		_, err := s.CurdRepo.GetUser(r.Context(), userID)
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

		err = s.CurdRepo.DeleteUser(r.Context(), userID)
		if err != nil {
			handler.ErrorResponse(w, http.StatusInternalServerError,
				handler.ErrResponse{
					Title:   "delete error",
					Details: err.Error(),
				},
			)
			return
		}

		resp := handler.DeleteUserResponse{
			Message: "user deleted successfully",
		}
		// Send the user delete response
		handler.SendResponse(w, resp, http.StatusOK)
	}
}
