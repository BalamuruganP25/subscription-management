package webhook_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"subscription-management/pkg/handler"
	"subscription-management/pkg/handler/webhook"
	stripe_webhook "subscription-management/pkg/handler/webhook"
	"subscription-management/pkg/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go/v78"
)

func TestWebhookHandler_SubscriptionUpdated(t *testing.T) {
	mockRepo := new(mocks.CrudRepo)
	processConfig := &handler.ProcessConfig{CurdRepo: mockRepo}

	// Prepare dummy subscription and event
	sub := stripe.Subscription{
		ID:     "sub_123456789",
		Status: stripe.SubscriptionStatusActive,
	}
	subRaw, _ := json.Marshal(sub)

	event := stripe.Event{
		ID:   "evt_123",
		Type: "customer.subscription.updated",
		Data: &stripe.EventData{Raw: subRaw},
	}

	// Override constructEventFunc to bypass actual Stripe signature verification
	stripe_webhook.ConstructEventFunc = func(payload []byte, sigHeader, secret string) (stripe.Event, error) {
		return event, nil
	}
	defer func() { stripe_webhook.ConstructEventFunc = stripe_webhook.ConstructEventFunc }()

	// Mock DB call expectation
	mockRepo.On("UpdateSubscription", mock.Anything, sub.ID, string(sub.Status)).Return(nil)

	// Create HTTP request with dummy body and headers
	payload, _ := json.Marshal(event)
	req := httptest.NewRequest("POST", "/api/webhooks/stripe", bytes.NewBuffer(payload))
	req.Header.Set("Stripe-Signature", "dummy_signature")
	os.Setenv("STRIPE_WEBHOOK_SECRET", "dummy_secret")

	rec := httptest.NewRecorder()

	// Call the handler
	handlerFunc := webhook.WebhookHandler(processConfig)
	handlerFunc.ServeHTTP(rec, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestWebhookHandler_SubscriptionDeleted(t *testing.T) {
	mockRepo := new(mocks.CrudRepo)
	processConfig := &handler.ProcessConfig{CurdRepo: mockRepo}

	sub := stripe.Subscription{
		ID:     "sub_987654321",
		Status: stripe.SubscriptionStatusCanceled,
	}
	subRaw, _ := json.Marshal(sub)

	event := stripe.Event{
		ID:   "evt_456",
		Type: "customer.subscription.deleted",
		Data: &stripe.EventData{Raw: subRaw},
	}

	webhook.ConstructEventFunc = func(payload []byte, sigHeader, secret string) (stripe.Event, error) {
		return event, nil
	}
	defer func() { webhook.ConstructEventFunc = stripe_webhook.ConstructEventFunc }()

	mockRepo.On("UpdateSubscription", mock.Anything, sub.ID, string(sub.Status)).Return(nil)

	payload, _ := json.Marshal(event)
	req := httptest.NewRequest("POST", "/api/webhooks/stripe", bytes.NewBuffer(payload))
	req.Header.Set("Stripe-Signature", "dummy_signature")
	os.Setenv("STRIPE_WEBHOOK_SECRET", "dummy_secret")

	rec := httptest.NewRecorder()

	handlerFunc := webhook.WebhookHandler(processConfig)
	handlerFunc.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}
