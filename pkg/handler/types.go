package handler

import "subscription-management/pkg/repository"

type ProcessConfig struct {
	CurdRepo  repository.CrudRepo
	StripeKey string
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

type CreateSubscriptionRequest struct {
	CustomerID string `json:"customer_id"`
	PriceID    string `json:"price_id"`
	PromoCode  string `json:"promo_code"`
}

type CreateSubscriptionResponse struct {
	CustomerID         string `json:"customer_id"`
	PriceID            string `json:"price_id"`
	PromoCode          string `json:"promo_code"`
	SubscriptionID     string `json:"subscription_id"`
	SubscriptionStatus string `json:"subscription_status"`
}
