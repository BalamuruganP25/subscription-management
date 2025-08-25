package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"subscription-management/pkg/handler"
)

func CreateUser(s *handler.ProcessConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req handler.UserRequest

		// Decode the request body
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			handler.ErrorResponse(w, http.StatusBadRequest,
				handler.ErrResponse{
					Title:   "payload error",
					Details: fmt.Sprintf("invalid request : %v", err),
				},
			)
			return
		}

		if err := ValidateUserRequest(req); err != nil {
			handler.ErrorResponse(w, http.StatusBadRequest,
				handler.ErrResponse{
					Title:   "validation error",
					Details: err.Error(),
				},
			)
			return

		}

		user_id, err := s.CurdRepo.CreateUser(r.Context(), req.Name, req.Email_id, req.Phone_number)
		if err != nil {
			handler.ErrorResponse(w, http.StatusInternalServerError,
				handler.ErrResponse{
					Title:   "internal error",
					Details: fmt.Sprintf("failed to create user: %v", err),
				},
			)
			return
		}

		response := handler.UserResponse{
			ID:           user_id,
			Name:         req.Name,
			Email_id:     req.Email_id,
			Phone_number: req.Phone_number,
			Status:       "Active",
		}
		// Send the created user in the response
		handler.SendResponse(w, response, http.StatusCreated)
	}
}
