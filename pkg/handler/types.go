package handler

import "subscription-management/pkg/repository"

type ProcessConfig struct {
	CurdRepo repository.CrudRepo
}

type UserRequest struct {
	Name         string `json:"name"`
	Email_id     string `json:"email_id"`
	Phone_number string `json:"phone_number"`
}
type UserResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email_id     string `json:"email_id"`
	Phone_number string `json:"phone_number"`
}

type UpdateUserRequest struct {
	Phone_number string `json:"phone_number"`
}

type CreateCustomerRequest struct {
	Name         string `json:"name"`
	Email_id     string `json:"email_id"`
	Phone_number string `json:"phone_number"`
}

type CreateCustomerResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email_id     string `json:"email_id"`
	Phone_number string `json:"phone_number"`
}
