package http

import (
	"bytes"
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	mock_auth "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/mocks"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/request"
	"github.com/golang/mock/gomock"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var testLogger *slog.Logger

func init() {
	testLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

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
			requestBody:    `{"username":"testuser2","password":"12345678a"}`,
			username:       "testuser2",
			password:       "12345678a",
			usecaseErr:     nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "AuthHandler_SignUp_Fail_1",
			requestBody:    `{"username":"testuser2","password":"12345678a"`,
			username:       "testuser2",
			password:       "12345678a",
			usecaseErr:     nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "AuthHandler_SignUp_Fail_2",
			requestBody:    `{"username":"test_user_2","password":"12345678"}`,
			username:       "testuser2",
			password:       "12345678",
			usecaseErr:     nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "AuthHandler_SignUp_Fail_3",
			requestBody:    `{"username":"testuser2","password":"12345678a"}`,
			username:       "testuser2",
			password:       "12345678a",
			usecaseErr:     errors.New("registration failed"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_auth.NewMockAuthUsecase(ctrl)
			mockBlocker := mock_auth.NewMockBlockerUsecase(ctrl)
			defer ctrl.Finish()

			if tt.name != "AuthHandler_SignUp_Fail_1" && tt.name != "AuthHandler_SignUp_Fail_2" {
				mockUsecase.EXPECT().SignUp(gomock.Any(), models.UserFormData{
					Username: tt.username,
					Password: tt.password,
				}).Return(models.User{
					Id:           uuid.NewV4(),
					Description:  "",
					Username:     tt.username,
					PasswordHash: request.GetHash(tt.password),
				}, "this_is_jwt_token", time.Now(), tt.usecaseErr)
			}

			req := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewBufferString(tt.requestBody))
			w := httptest.NewRecorder()

			handler := CreateAuthHandler(mockUsecase, mockBlocker, testLogger)
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
			requestBody:    `{"username":"testuser2","password":"12345678a"}`,
			username:       "testuser2",
			password:       "12345678a",
			usecaseErr:     nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "AuthHandler_SignIn_Fail_1",
			requestBody:    `{"username":"testuser2","password":"12345678a"`,
			username:       "testuser2",
			password:       "12345678a",
			usecaseErr:     nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "AuthHandler_SignIn_Fail_2",
			requestBody:    `{"username":"testuser2","password":"12345678a"}`,
			username:       "testuser2",
			password:       "12345678a",
			usecaseErr:     errors.New("registration failed"),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_auth.NewMockAuthUsecase(ctrl)
			mockBlocker := mock_auth.NewMockBlockerUsecase(ctrl)
			defer ctrl.Finish()

			if tt.name != "AuthHandler_SignIn_Fail_1" {
				mockUsecase.EXPECT().SignIn(gomock.Any(), models.UserFormData{
					Username: tt.username,
					Password: tt.password,
				}).Return(models.User{
					Id:           uuid.NewV4(),
					Description:  "",
					Username:     tt.username,
					PasswordHash: request.GetHash(tt.password),
				}, "this_is_jwt_token", time.Now(), tt.usecaseErr)
			}

			mockBlocker.EXPECT().CheckLoginAttempts(gomock.Any(), gomock.Any()).Return(nil)

			req := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewBufferString(tt.requestBody))
			w := httptest.NewRecorder()

			handler := CreateAuthHandler(mockUsecase, mockBlocker, testLogger)
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
			mockBlocker := mock_auth.NewMockBlockerUsecase(ctrl)
			defer ctrl.Finish()

			req := httptest.NewRequest("DELETE", "http://example.com/api/handler", nil)
			w := httptest.NewRecorder()

			handler := CreateAuthHandler(mockUsecase, mockBlocker, testLogger)
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
			username:       "testuser2",
			usecaseErr:     nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "AuthHandler_CheckUser_Fail_1",
			id:             uuid.NewV4(),
			username:       "testuser2",
			usecaseErr:     nil,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_auth.NewMockAuthUsecase(ctrl)
			mockBlocker := mock_auth.NewMockBlockerUsecase(ctrl)
			defer ctrl.Finish()

			req := httptest.NewRequest("GET", "http://example.com/api/handler", nil)
			w := httptest.NewRecorder()

			ctx := context.WithValue(req.Context(), models.PayloadContextKey, models.JwtPayload{Id: tt.id, Username: tt.username})
			if tt.name == "AuthHandler_CheckUser_Fail_1" {
				ctx = context.WithValue(req.Context(), models.PayloadContextKey, models.Note{})
			}
			req = req.WithContext(ctx)

			handler := CreateAuthHandler(mockUsecase, mockBlocker, testLogger)
			handler.CheckUser(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAuthHandler_GetProfile(t *testing.T) {
	var tests = []struct {
		name           string
		id             uuid.UUID
		username       string
		usecaseErr     error
		expectedStatus int
	}{
		{
			name:           "AuthHandler_GetProfile_Success",
			id:             uuid.NewV4(),
			username:       "testuser2",
			usecaseErr:     nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "AuthHandler_GetProfile_Fail_1",
			id:             uuid.NewV4(),
			username:       "testuser2",
			usecaseErr:     nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "AuthHandler_GetProfile_Fail_2",
			id:             uuid.NewV4(),
			username:       "testuser2",
			usecaseErr:     errors.New("error in CheckUser"),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_auth.NewMockAuthUsecase(ctrl)
			mockBlocker := mock_auth.NewMockBlockerUsecase(ctrl)
			defer ctrl.Finish()

			if tt.name != "AuthHandler_GetProfile_Fail_1" {
				mockUsecase.EXPECT().CheckUser(gomock.Any(), tt.id).Return(models.User{
					Id:           tt.id,
					Description:  "",
					Username:     tt.username,
					PasswordHash: request.GetHash("fh9ch283c"),
				}, tt.usecaseErr)
			}

			req := httptest.NewRequest("GET", "http://example.com/api/handler", nil)
			w := httptest.NewRecorder()

			ctx := context.WithValue(req.Context(), models.PayloadContextKey, models.JwtPayload{Id: tt.id, Username: tt.username})
			if tt.name == "AuthHandler_GetProfile_Fail_1" {
				ctx = context.WithValue(req.Context(), models.PayloadContextKey, models.Note{})
			}
			req = req.WithContext(ctx)

			handler := CreateAuthHandler(mockUsecase, mockBlocker, testLogger)
			handler.GetProfile(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
