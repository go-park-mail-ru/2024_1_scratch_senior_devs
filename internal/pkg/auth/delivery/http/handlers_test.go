package http

import (
	"bytes"
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	mock_auth "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/mocks"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthHandler_SignUp(t *testing.T) {
	var tests = []struct {
		name           string
		requestBody    string
		username       string
		password       string
		usecaseErr     error
		expectedStatus int
	}{
		{
			name:           "AuthHandler_SignUp_Success",
			requestBody:    `{"username":"test_user_2","password":"12345678a"}`,
			username:       "test_user_2",
			password:       "12345678a",
			usecaseErr:     nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "AuthHandler_SignUp_Fail_1",
			requestBody:    `{"username":"test_user_2","password":"12345678a"`,
			username:       "test_user_2",
			password:       "12345678a",
			usecaseErr:     nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "AuthHandler_SignUp_Fail_2",
			requestBody:    `{"username":"test_user_2","password":"12345678"}`,
			username:       "test_user_2",
			password:       "12345678",
			usecaseErr:     nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "AuthHandler_SignUp_Fail_3",
			requestBody:    `{"username":"test_user_2","password":"12345678a"}`,
			username:       "test_user_2",
			password:       "12345678a",
			usecaseErr:     errors.New("registration failed"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_auth.NewMockAuthUsecase(ctrl)
			defer ctrl.Finish()

			if tt.name != "AuthHandler_SignUp_Fail_1" && tt.name != "AuthHandler_SignUp_Fail_2" {
				mockUsecase.EXPECT().SignUp(gomock.Any(), &models.UserFormData{
					Username: tt.username,
					Password: tt.password,
				}).Return(&models.User{
					Id:           uuid.NewV4(),
					Description:  nil,
					Username:     tt.username,
					PasswordHash: utils.GetHash(tt.password),
				}, "this_is_jwt_token", time.Now(), tt.usecaseErr)
			}

			req := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewBufferString(tt.requestBody))
			w := httptest.NewRecorder()

			handler := CreateAuthHandler(mockUsecase)
			handler.SignUp(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAuthHandler_SignIn(t *testing.T) {
	var tests = []struct {
		name           string
		requestBody    string
		username       string
		password       string
		usecaseErr     error
		expectedStatus int
	}{
		{
			name:           "AuthHandler_SignIn_Success",
			requestBody:    `{"username":"test_user_2","password":"12345678a"}`,
			username:       "test_user_2",
			password:       "12345678a",
			usecaseErr:     nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "AuthHandler_SignIn_Fail_1",
			requestBody:    `{"username":"test_user_2","password":"12345678a"`,
			username:       "test_user_2",
			password:       "12345678a",
			usecaseErr:     nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "AuthHandler_SignIn_Fail_2",
			requestBody:    `{"username":"test_user_2","password":"12345678a"}`,
			username:       "test_user_2",
			password:       "12345678a",
			usecaseErr:     errors.New("registration failed"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_auth.NewMockAuthUsecase(ctrl)
			defer ctrl.Finish()

			if tt.name != "AuthHandler_SignIn_Fail_1" {
				mockUsecase.EXPECT().SignIn(gomock.Any(), &models.UserFormData{
					Username: tt.username,
					Password: tt.password,
				}).Return(&models.User{
					Id:           uuid.NewV4(),
					Description:  nil,
					Username:     tt.username,
					PasswordHash: utils.GetHash(tt.password),
				}, "this_is_jwt_token", time.Now(), tt.usecaseErr)
			}

			req := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewBufferString(tt.requestBody))
			w := httptest.NewRecorder()

			handler := CreateAuthHandler(mockUsecase)
			handler.SignIn(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAuthHandler_LogOut(t *testing.T) {
	var tests = []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "AuthHandler_SignUp_Success",
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_auth.NewMockAuthUsecase(ctrl)
			defer ctrl.Finish()

			req := httptest.NewRequest("DELETE", "http://example.com/api/handler", nil)
			w := httptest.NewRecorder()

			handler := CreateAuthHandler(mockUsecase)
			handler.LogOut(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAuthHandler_CheckUser(t *testing.T) {
	var tests = []struct {
		name           string
		id             uuid.UUID
		username       string
		usecaseErr     error
		expectedStatus int
	}{
		{
			name:           "AuthHandler_CheckUser_Success",
			id:             uuid.NewV4(),
			username:       "test_user_2",
			usecaseErr:     nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "AuthHandler_CheckUser_Fail_1",
			id:             uuid.NewV4(),
			username:       "test_user_2",
			usecaseErr:     nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "AuthHandler_CheckUser_Fail_2",
			id:             uuid.NewV4(),
			username:       "test_user_2",
			usecaseErr:     errors.New("error in CheckUser"),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_auth.NewMockAuthUsecase(ctrl)
			defer ctrl.Finish()

			if tt.name != "AuthHandler_CheckUser_Fail_1" {
				mockUsecase.EXPECT().CheckUser(gomock.Any(), tt.id).Return(&models.User{
					Id:           tt.id,
					Description:  nil,
					Username:     tt.username,
					PasswordHash: utils.GetHash("fh9ch283c"),
				}, tt.usecaseErr)
			}

			req := httptest.NewRequest("GET", "http://example.com/api/handler", nil)
			w := httptest.NewRecorder()

			ctx := context.Background()
			if tt.name == "AuthHandler_CheckUser_Fail_1" {
				ctx = context.WithValue(req.Context(), "payload", models.Note{})
			} else {
				ctx = context.WithValue(req.Context(), "payload", models.JwtPayload{Id: tt.id, Username: tt.username})
			}
			req = req.WithContext(ctx)

			handler := CreateAuthHandler(mockUsecase)
			handler.CheckUser(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
