package handler

import "subscription-management/pkg/repository"

type ProcessConfig struct {
	CurdRepo repository.CrudRepo
}

type UserRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone"`
}
type UserResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone"`
}
