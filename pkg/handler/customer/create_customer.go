package customer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"subscription-management/pkg/handler"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
)

func CreateCustomer(s *handler.ProcessConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implementation for creating a customer goes here
		var req handler.CreateCustomerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handler.ErrorResponse(w, http.StatusBadRequest, handler.ErrResponse{
				Title:   "validation error",
				Details: "invalid request payload",
			})
			return
		}

		if err := ValidateCustomerRequest(req); err != nil {
			handler.ErrorResponse(w, http.StatusBadRequest,
				handler.ErrResponse{
					Title:   "validation error",
					Details: err.Error(),
				},
			)
			return

		}
		// soon i will add stripe key from env
		// stripe.Key = os.Getenv("STRIPE_API_KEY")
		// Create customer in Stripe

		params := &stripe.CustomerParams{
			Name:  stripe.String(req.Name),
			Email: stripe.String(req.Email_id),
			Phone: stripe.String(req.Phone_number),
		}
		c, err := customer.New(params)
		if err != nil {
			handler.ErrorResponse(w, http.StatusInternalServerError, handler.ErrResponse{
				Title:   "internal error",
				Details: fmt.Sprintf("failed to create customer in Stripe: %v", err),
			})
			return
		}

		// Create the customer in the database
		err = s.CurdRepo.CreateCustomer(r.Context(), c.ID, req.Name, req.Email_id, req.Phone_number)
		if err != nil {
			handler.ErrorResponse(w, http.StatusInternalServerError, handler.ErrResponse{
				Title:   "internal error",
				Details: fmt.Sprintf("failed to create customer: %v", err),
			})
			return
		}
		customer := handler.CreateCustomerResponse{
			ID:           c.ID,
			Name:         req.Name,
			Email_id:     req.Email_id,
			Phone_number: req.Phone_number,
		}

		// Send the created customer in the response
		handler.SendResponse(w, customer, http.StatusCreated)
	}
}
