package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"subscription-management/pkg/handler"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
)

type ConstructEventFun func(payload []byte, header, secret string) (stripe.Event, error)

var ConstructEventFunc ConstructEventFun = webhook.ConstructEvent

func WebhookHandler(config *handler.ProcessConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const MaxBodyBytes = int64(65536)
		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			handler.ErrorResponse(w, http.StatusServiceUnavailable,
				handler.ErrResponse{
					Title:   "payload error",
					Details: fmt.Sprintf("req payload read error: %v", err),
				},
			)
			return

		}

		endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
		sigHeader := r.Header.Get("Stripe-Signature")

		event, err := ConstructEventFunc(payload, sigHeader, endpointSecret)
		if err != nil {
			handler.ErrorResponse(w, http.StatusBadRequest,
				handler.ErrResponse{
					Title:   "payload error",
					Details: fmt.Sprintf("Webhook signature verification failed: %v", err),
				},
			)
			return
		}
		// Handle the event
		// Use a retry mechanism
		handleSubscription := func(subscription stripe.Subscription) error {
			maxRetries := 3
			retryDelay := 500 * time.Millisecond
			for i := 0; i < maxRetries; i++ {
				err := config.CurdRepo.UpdateSubscription(r.Context(), subscription.ID, string(subscription.Status))
				if err == nil {
					return nil
				}
				time.Sleep(retryDelay)
				retryDelay *= 2 // Exponential backoff
			}
			return fmt.Errorf("update failed after %d retries", maxRetries)
		}

		switch event.Type {
		case "customer.subscription.updated", "customer.subscription.deleted":
			var subscription stripe.Subscription
			if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
				handler.ErrorResponse(w, http.StatusBadRequest,
					handler.ErrResponse{
						Title:   "payload error",
						Details: "Error parsing webhook JSON",
					},
				)
				return
			}

			err := handleSubscription(subscription)
			if err != nil {
				handler.ErrorResponse(w, http.StatusInternalServerError,
					handler.ErrResponse{
						Title:   "database error",
						Details: err.Error(),
					},
				)
				return
			}

		default:
			handler.ErrorResponse(w, http.StatusNotFound,
				handler.ErrResponse{
					Title:   "mismatch event type",
					Details: fmt.Sprintf("Unhandled event type: %s", event.Type),
				},
			)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
