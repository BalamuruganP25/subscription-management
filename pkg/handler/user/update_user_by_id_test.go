package user_test

import (
	"bytes"
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

func TestUpdateUserById(t *testing.T) {
	tests := []struct {
		name         string
		userID       string
		body         string
		mockFunc     func(repo *mocks.CrudRepo)
		wantCode     int
		wantResponse any
	}{
		{
			name:   "success",
			userID: "1",
			body:   `{"phone_number":"9876543210"}`,
			mockFunc: func(repo *mocks.CrudRepo) {
				repo.On("GetUser", mock.Anything, "1").Return(repository.UserResponse{
					ID:           "1",
					Name:         "Alice",
					Email_id:     "alice@example.com",
					Phone_number: "1234567890",
				}, nil)
				repo.On("UpdateUser", mock.Anything, "1", "9876543210").Return(nil)
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
			name:         "missing id",
			userID:       "",
			body:         `{"phone_number":"9876543210"}`,
			mockFunc:     nil,
			wantCode:     http.StatusBadRequest,
			wantResponse: handler.ErrResponse{Title: "validation error", Details: "user id should be empty"},
		},
		{
			name:         "invalid json",
			userID:       "1",
			body:         `{"phone":}`,
			mockFunc:     nil,
			wantCode:     http.StatusBadRequest,
			wantResponse: handler.ErrResponse{Title: "invalid request"},
		},
		{
			name:   "user not found",
			userID: "2",
			body:   `{"phone_number":"9876543210"}`,
			mockFunc: func(repo *mocks.CrudRepo) {
				repo.On("GetUser", mock.Anything, "2").Return(repository.UserResponse{}, errors.New("not found"))
			},
			wantCode: http.StatusInternalServerError,
			wantResponse: handler.ErrResponse{
				Title:   "internal error",
				Details: "failed to fetch user: not found",
			},
		},
		{
			name:   "update error",
			userID: "1",
			body:   `{"phone_number":"9876543210"}`,
			mockFunc: func(repo *mocks.CrudRepo) {
				repo.On("GetUser", mock.Anything, "1").Return(repository.UserResponse{
					ID:           "1",
					Name:         "Alice",
					Email_id:     "alice@example.com",
					Phone_number: "1234567890",
				}, nil)
				repo.On("UpdateUser", mock.Anything, "1", "9876543210").Return(errors.New("update failed"))
			},
			wantCode: http.StatusInternalServerError,
			wantResponse: handler.ErrResponse{
				Title:   "update error",
				Details: "update failed",
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
			handlerFunc := user.UpdateUserById(&config)

			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/user/%s", tt.userID), bytes.NewBufferString(tt.body))
			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("id", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
			rec := httptest.NewRecorder()

			handlerFunc(rec, req)

			assert.Equal(t, tt.wantCode, rec.Code)

			if tt.wantCode == http.StatusOK {
				var got repository.UserResponse
				err := json.NewDecoder(rec.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResponse, got)
			} else {
				var got handler.ErrResponse
				_ = json.NewDecoder(rec.Body).Decode(&got)
				wantErrResp, ok := tt.wantResponse.(handler.ErrResponse)
				if assert.True(t, ok, "wantResponse should be of type handler.ErrResponse") {
					assert.Equal(t, wantErrResp.Title, got.Title)
				}
			}

			mockRepo.AssertExpectations(t)

		})
	}
}
