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

func TestDeleteUserById(t *testing.T) {
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
					Status:       "active",
				}, nil)
				repo.On("DeleteUser", mock.Anything, "1").Return(nil)
			},
			wantCode: http.StatusOK,
			wantResponse: handler.DeleteUserResponse{
				Message: "user deleted successfully",
			},
		},
		{
			name:     "missing id",
			userID:   "",
			mockFunc: nil,
			wantCode: http.StatusBadRequest,
			wantResponse: handler.ErrResponse{
				Title:   "invalid request",
				Details: "user id should be empty",
			},
		},
		{
			name:   "user not found",
			userID: "2",
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
			name:   "delete error",
			userID: "1",
			mockFunc: func(repo *mocks.CrudRepo) {
				repo.On("GetUser", mock.Anything, "1").Return(repository.UserResponse{
					ID:           "1",
					Name:         "Alice",
					Email_id:     "alice@example.com",
					Phone_number: "1234567890",
				}, nil)
				repo.On("DeleteUser", mock.Anything, "1").Return(errors.New("delete failed"))
			},
			wantCode: http.StatusInternalServerError,
			wantResponse: handler.ErrResponse{
				Title:   "delete error",
				Details: "delete failed",
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
			handlerFunc := user.DeleteUserById(&config)

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/user/%s", tt.userID), nil)
			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("id", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
			rec := httptest.NewRecorder()

			handlerFunc(rec, req)

			assert.Equal(t, tt.wantCode, rec.Code)

			if tt.wantCode == http.StatusOK {
				var got handler.DeleteUserResponse
				err := json.NewDecoder(rec.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResponse, got)
			} else {
				var got handler.ErrResponse
				_ = json.NewDecoder(rec.Body).Decode(&got)
				wantErrResp := tt.wantResponse.(handler.ErrResponse)
				assert.Equal(t, wantErrResp.Title, got.Title)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
