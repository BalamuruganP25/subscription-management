package user_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"subscription-management/pkg/handler"
	"subscription-management/pkg/handler/user"
	"subscription-management/pkg/mocks"
	"subscription-management/pkg/repository"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUser(t *testing.T) {
	tests := []struct {
		name         string
		userID       string
		mockFunc     func(repo *mocks.CrudRepo)
		wantCode     int
		wantResponse any
	}{
		{
			name:   "success",
			userID: "1",
			mockFunc: func(repo *mocks.CrudRepo) {
				repo.On("GetUser", mock.Anything, "1").Return(repository.UserResponse{
					ID:           "1",
					Name:         "Alice",
					Email_id:     "alice@example.com",
					Phone_number: "1234567890",
				}, nil)
			},
			wantCode: http.StatusOK,
			wantResponse: repository.UserResponse{
				ID:           "1",
				Name:         "Alice",
				Email_id:     "alice@example.com",
				Phone_number: "1234567890",
			},
		},
		{
			name:   "not found",
			userID: "2",
			mockFunc: func(repo *mocks.CrudRepo) {
				repo.On("GetUser", mock.Anything, "2").Return(repository.UserResponse{}, errors.New("not found"))
			},
			wantCode: http.StatusNotFound,
			wantResponse: handler.ErrResponse{
				Title:   "not found",
				Details: "user not found: not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.CrudRepo)
			if tt.mockFunc != nil {
				tt.mockFunc(mockRepo)
			}
			config := handler.ProcessConfig{CurdRepo: mockRepo}
			handlerFunc := user.GetUserById(&config)

			router := chi.NewRouter()
			// Register both routes
			router.Get("/v1/user/{id}", handlerFunc)

			req := httptest.NewRequest("GET", fmt.Sprintf("/v1/user/%s", tt.userID), nil)
			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("id", tt.userID)
			req = req.WithContext(
				context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
			)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantCode, rec.Code)

			// Decode and check response body
			if tt.wantCode == http.StatusOK {
				var got repository.UserResponse
				err := json.NewDecoder(rec.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResponse, got)
			} else {
				var got handler.ErrResponse
				err := json.NewDecoder(rec.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResponse, got)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
