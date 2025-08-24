package customer

import (
	"encoding/json"
	"net/http"
	"subscription-management/pkg/handler"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/subscription"
)

func CreateSubscription(config *handler.ProcessConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req handler.CreateSubscriptionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Validate required fields
		if err := ValidateSubscriptionRequest(req); err != nil {
			handler.ErrorResponse(w, http.StatusBadRequest,
				handler.ErrResponse{
					Title:   "validation error",
					Details: err.Error(),
				},
			)
			return
		}
		stripe.Key = config.StripeKey

		params := &stripe.SubscriptionParams{
			Customer: stripe.String(req.CustomerID),
			Items: []*stripe.SubscriptionItemsParams{
				{Price: stripe.String(req.PriceID)},
			},
			PaymentBehavior: stripe.String("default_incomplete"),
			Expand:          []*string{stripe.String("latest_invoice.payment_intent")},
		}

		if req.PromoCode != "" {
			params.Discounts = []*stripe.SubscriptionDiscountParams{
				{PromotionCode: stripe.String(req.PromoCode)},
			}
		}
		sub, err := subscription.New(params)
		if err != nil {
			handler.ErrorResponse(w, http.StatusInternalServerError, handler.ErrResponse{
				Title:   "internal error",
				Details: err.Error(),
			})
			return
		}

		// Store subscription details in the database
		err = config.CurdRepo.CreateSubscription(r.Context(), req.CustomerID, req.PriceID, req.PromoCode, sub.ID, string(sub.Status))
		if err != nil {
			handler.ErrorResponse(w, http.StatusInternalServerError, handler.ErrResponse{
				Title:   "database operation error",
				Details: err.Error(),
			})
			return
		}

		subscription := handler.CreateSubscriptionResponse{
			CustomerID:         req.CustomerID,
			PriceID:            req.PriceID,
			PromoCode:          req.PromoCode,
			SubscriptionID:     sub.ID,
			SubscriptionStatus: string(sub.Status),
		}

		// Send the created subscription in the response
		handler.SendResponse(w, subscription, http.StatusCreated)
	}

}
