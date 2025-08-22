package repository

type UserResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email_id     string `json:"email_id"`
	Phone_number string `json:"phone_number"`
}
