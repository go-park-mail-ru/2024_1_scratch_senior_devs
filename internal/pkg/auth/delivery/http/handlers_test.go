package http

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	mock_auth "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/mocks"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/responses"
	"github.com/golang/mock/gomock"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
)

var testLogger *slog.Logger

func init() {
	testLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func TestAuthHandler_SignUp(t *testing.T) {
	testConfig := config.Config{
		AuthHandler: config.AuthHandlerConfig{
			QrIssuer:              "YouNote",
			AvatarMaxFormDataSize: 31457280,
			AvatarFileTypes: map[string]string{
				"image/jpeg": ".jpeg",
				"image/png":  ".png",
			},
			Jwt: config.JwtConfig{
				JwtCookie: "YouNoteJWT",
			},
			Csrf: config.CsrfConfig{
				CsrfCookie:   "YouNoteCSRF",
				CSRFLifeTime: time.Duration(24 * time.Hour),
			},
		},
		Validation: config.ValidationConfig{
			MinUsernameLength:    4,
			MaxUsernameLength:    12,
			MinPasswordLength:    8,
			MaxPasswordLength:    20,
			PasswordAllowedExtra: "$%&#",
			SecretLength:         6,
		},
	}
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
					PasswordHash: responses.GetHash(tt.password),
				}, "this_is_jwt_token", time.Now(), tt.usecaseErr)
			}

			req := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewBufferString(tt.requestBody))
			w := httptest.NewRecorder()

			handler := CreateAuthHandler(mockUsecase, mockBlocker, testLogger, testConfig.AuthHandler, testConfig.Validation)
			handler.SignUp(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAuthHandler_SignIn(t *testing.T) {
	testConfig := config.Config{
		AuthHandler: config.AuthHandlerConfig{
			QrIssuer:              "YouNote",
			AvatarMaxFormDataSize: 31457280,
			AvatarFileTypes: map[string]string{
				"image/jpeg": ".jpeg",
				"image/png":  ".png",
			},
			Jwt: config.JwtConfig{
				JwtCookie: "YouNoteJWT",
			},
			Csrf: config.CsrfConfig{
				CsrfCookie:   "YouNoteCSRF",
				CSRFLifeTime: time.Duration(24 * time.Hour),
			},
		},
		Validation: config.ValidationConfig{
			MinUsernameLength:    4,
			MaxUsernameLength:    12,
			MinPasswordLength:    8,
			MaxPasswordLength:    20,
			PasswordAllowedExtra: "$%&#",
			SecretLength:         6,
		},
	}
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
					PasswordHash: responses.GetHash(tt.password),
				}, "this_is_jwt_token", time.Now(), tt.usecaseErr)
			}

			mockBlocker.EXPECT().CheckLoginAttempts(gomock.Any(), gomock.Any()).Return(nil)

			req := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewBufferString(tt.requestBody))
			w := httptest.NewRecorder()

			handler := CreateAuthHandler(mockUsecase, mockBlocker, testLogger, testConfig.AuthHandler, testConfig.Validation)
			handler.SignIn(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAuthHandler_LogOut(t *testing.T) {
	testConfig := config.Config{
		AuthHandler: config.AuthHandlerConfig{
			QrIssuer:              "YouNote",
			AvatarMaxFormDataSize: 31457280,
			AvatarFileTypes: map[string]string{
				"image/jpeg": ".jpeg",
				"image/png":  ".png",
			},
			Jwt: config.JwtConfig{
				JwtCookie: "YouNoteJWT",
			},
			Csrf: config.CsrfConfig{
				CsrfCookie:   "YouNoteCSRF",
				CSRFLifeTime: time.Duration(24 * time.Hour),
			},
		},
		Validation: config.ValidationConfig{
			MinUsernameLength:    4,
			MaxUsernameLength:    12,
			MinPasswordLength:    8,
			MaxPasswordLength:    20,
			PasswordAllowedExtra: "$%&#",
			SecretLength:         6,
		},
	}
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

			handler := CreateAuthHandler(mockUsecase, mockBlocker, testLogger, testConfig.AuthHandler, testConfig.Validation)
			handler.LogOut(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAuthHandler_CheckUser(t *testing.T) {
	testConfig := config.Config{
		AuthHandler: config.AuthHandlerConfig{
			QrIssuer:              "YouNote",
			AvatarMaxFormDataSize: 31457280,
			AvatarFileTypes: map[string]string{
				"image/jpeg": ".jpeg",
				"image/png":  ".png",
			},
			Jwt: config.JwtConfig{
				JwtCookie: "YouNoteJWT",
			},
			Csrf: config.CsrfConfig{
				CsrfCookie:   "YouNoteCSRF",
				CSRFLifeTime: time.Duration(24 * time.Hour),
			},
		},
		Validation: config.ValidationConfig{
			MinUsernameLength:    4,
			MaxUsernameLength:    12,
			MinPasswordLength:    8,
			MaxPasswordLength:    20,
			PasswordAllowedExtra: "$%&#",
			SecretLength:         6,
		},
	}
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

			ctx := context.WithValue(req.Context(), config.PayloadContextKey, models.JwtPayload{Id: tt.id, Username: tt.username})
			if tt.name == "AuthHandler_CheckUser_Fail_1" {
				ctx = context.WithValue(req.Context(), config.PayloadContextKey, models.Note{})
			}
			req = req.WithContext(ctx)

			handler := CreateAuthHandler(mockUsecase, mockBlocker, testLogger, testConfig.AuthHandler, testConfig.Validation)
			handler.CheckUser(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAuthHandler_GetProfile(t *testing.T) {
	testConfig := config.Config{
		AuthHandler: config.AuthHandlerConfig{
			QrIssuer:              "YouNote",
			AvatarMaxFormDataSize: 31457280,
			AvatarFileTypes: map[string]string{
				"image/jpeg": ".jpeg",
				"image/png":  ".png",
			},
			Jwt: config.JwtConfig{
				JwtCookie: "YouNoteJWT",
			},
			Csrf: config.CsrfConfig{
				CsrfCookie:   "YouNoteCSRF",
				CSRFLifeTime: time.Duration(24 * time.Hour),
			},
		},
		Validation: config.ValidationConfig{
			MinUsernameLength:    4,
			MaxUsernameLength:    12,
			MinPasswordLength:    8,
			MaxPasswordLength:    20,
			PasswordAllowedExtra: "$%&#",
			SecretLength:         6,
		},
	}
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
					PasswordHash: responses.GetHash("fh9ch283c"),
				}, tt.usecaseErr)
			}

			req := httptest.NewRequest("GET", "http://example.com/api/handler", nil)
			w := httptest.NewRecorder()

			ctx := context.WithValue(req.Context(), config.PayloadContextKey, models.JwtPayload{Id: tt.id, Username: tt.username})
			if tt.name == "AuthHandler_GetProfile_Fail_1" {
				ctx = context.WithValue(req.Context(), config.PayloadContextKey, models.Note{})
			}
			req = req.WithContext(ctx)

			handler := CreateAuthHandler(mockUsecase, mockBlocker, testLogger, testConfig.AuthHandler, testConfig.Validation)
			handler.GetProfile(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAuthHandler_UpdateProfile(t *testing.T) {
	testConfig := config.Config{
		AuthHandler: config.AuthHandlerConfig{
			QrIssuer:              "YouNote",
			AvatarMaxFormDataSize: 31457280,
			AvatarFileTypes: map[string]string{
				"image/jpeg": ".jpeg",
				"image/png":  ".png",
			},
			Jwt: config.JwtConfig{
				JwtCookie: "YouNoteJWT",
			},
			Csrf: config.CsrfConfig{
				CsrfCookie:   "YouNoteCSRF",
				CSRFLifeTime: time.Duration(24 * time.Hour),
			},
		},
		Validation: config.ValidationConfig{
			MinUsernameLength:    4,
			MaxUsernameLength:    12,
			MinPasswordLength:    8,
			MaxPasswordLength:    20,
			PasswordAllowedExtra: "$%&#",
			SecretLength:         6,
		},
	}
	userID := uuid.FromStringOrNil("ac6966bc-3c26-45a0-963e-b168fc34fd79")
	username := "user"
	type args struct {
		userID   uuid.UUID
		username string
	}
	tests := []struct {
		args       args
		wantStatus int
		name       string
		payload    string
		ucMocker   func(ctx context.Context, uc *mock_auth.MockAuthUsecase, blockerUc *mock_auth.MockBlockerUsecase)
	}{
		{
			name: "Test_Update_Unauthorized",
			ucMocker: func(ctx context.Context, uc *mock_auth.MockAuthUsecase, blockerUc *mock_auth.MockBlockerUsecase) {
			},
			wantStatus: http.StatusUnauthorized,
		},

		{
			args: args{
				userID:   userID,
				username: username,
			},
			name:    "Test_BadRequest",
			payload: "",
			ucMocker: func(ctx context.Context, uc *mock_auth.MockAuthUsecase, blockerUc *mock_auth.MockBlockerUsecase) {
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			args: args{
				userID:   userID,
				username: username,
			},
			name: "Test_Success",
			payload: `
			{
				"description": "slkakjckld",
				"password": {
				    "old": "12345678a",
				    "new": "12345678b"
				}
			 }`,
			ucMocker: func(ctx context.Context, uc *mock_auth.MockAuthUsecase, blockerUc *mock_auth.MockBlockerUsecase) {
				uc.EXPECT().UpdateProfile(ctx, userID, gomock.Any()).Return(models.User{
					Id:           userID,
					Description:  "slkakjckld",
					Username:     username,
					PasswordHash: responses.GetHash("12345678b"),
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			args: args{
				userID:   userID,
				username: username,
			},
			name: "Test_BadRequest_2",
			payload: `
			{
				"description": "slkakjckld",
				"password": {
				    "old": "12345678a",
				    "new": "12345678b"
				}
			 }`,
			ucMocker: func(ctx context.Context, uc *mock_auth.MockAuthUsecase, blockerUc *mock_auth.MockBlockerUsecase) {
				uc.EXPECT().UpdateProfile(ctx, userID, gomock.Any()).Return(models.User{}, errors.New("error"))
			},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			uc := mock_auth.NewMockAuthUsecase(ctrl)
			blockerUc := mock_auth.NewMockBlockerUsecase(ctrl)
			defer ctrl.Finish()
			req := httptest.NewRequest("POST", "http://example.com/api/handler/", bytes.NewBufferString(tt.payload))
			w := httptest.NewRecorder()
			ctx := context.WithValue(req.Context(), config.PayloadContextKey, models.JwtPayload{Id: userID, Username: username})
			req = req.WithContext(ctx)
			if tt.name == "Test_Update_Unauthorized" {
				req = req.WithContext(context.Background())
			}

			tt.ucMocker(req.Context(), uc, blockerUc)

			h := CreateAuthHandler(uc, blockerUc, testLogger, testConfig.AuthHandler, testConfig.Validation)

			h.UpdateProfile(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestAuthHandler_DisableSecondFactor(t *testing.T) {
	testConfig := config.Config{
		AuthHandler: config.AuthHandlerConfig{
			QrIssuer:              "YouNote",
			AvatarMaxFormDataSize: 31457280,
			AvatarFileTypes: map[string]string{
				"image/jpeg": ".jpeg",
				"image/png":  ".png",
			},
			Jwt: config.JwtConfig{
				JwtCookie: "YouNoteJWT",
			},
			Csrf: config.CsrfConfig{
				CsrfCookie:   "YouNoteCSRF",
				CSRFLifeTime: time.Duration(24 * time.Hour),
			},
		},
		Validation: config.ValidationConfig{
			MinUsernameLength:    4,
			MaxUsernameLength:    12,
			MinPasswordLength:    8,
			MaxPasswordLength:    20,
			PasswordAllowedExtra: "$%&#",
			SecretLength:         6,
		},
	}

	userID := uuid.FromStringOrNil("ac6966bc-3c26-45a0-963e-b168fc34fd79")
	username := "test"
	type args struct {
		userID   uuid.UUID
		username string
	}
	tests := []struct {
		wantStatus int
		name       string
		ucMocker   func(ctx context.Context, uc *mock_auth.MockAuthUsecase, blockerUc *mock_auth.MockBlockerUsecase)
		args       args
	}{
		{
			name: "Test_Success",
			ucMocker: func(ctx context.Context, uc *mock_auth.MockAuthUsecase, blockerUc *mock_auth.MockBlockerUsecase) {
				uc.EXPECT().DeleteSecret(ctx, username).Return(nil)
			},
			args: args{
				userID:   userID,
				username: username,
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "Test_BadRequest",
			ucMocker: func(ctx context.Context, uc *mock_auth.MockAuthUsecase, blockerUc *mock_auth.MockBlockerUsecase) {
				uc.EXPECT().DeleteSecret(ctx, username).Return(errors.New("error"))
			},
			args: args{
				userID:   userID,
				username: username,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Test_Unauthorized",
			ucMocker: func(ctx context.Context, uc *mock_auth.MockAuthUsecase, blockerUc *mock_auth.MockBlockerUsecase) {

			},
			args: args{
				userID:   userID,
				username: username,
			},
			wantStatus: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			uc := mock_auth.NewMockAuthUsecase(ctrl)
			blockerUc := mock_auth.NewMockBlockerUsecase(ctrl)
			defer ctrl.Finish()
			req := httptest.NewRequest("GET", "http://example.com/api/handler/", bytes.NewBufferString(""))
			w := httptest.NewRecorder()
			ctx := context.WithValue(req.Context(), config.PayloadContextKey, models.JwtPayload{Id: userID, Username: username})
			req = req.WithContext(ctx)
			if tt.name == "Test_Unauthorized" {
				req = req.WithContext(context.Background())
			}

			tt.ucMocker(req.Context(), uc, blockerUc)

			h := CreateAuthHandler(uc, blockerUc, testLogger, testConfig.AuthHandler, testConfig.Validation)

			h.DisableSecondFactor(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestAuthHandler_GetQRCode(t *testing.T) {
	testConfig := config.Config{
		AuthHandler: config.AuthHandlerConfig{
			QrIssuer:              "YouNote",
			AvatarMaxFormDataSize: 31457280,
			AvatarFileTypes: map[string]string{
				"image/jpeg": ".jpeg",
				"image/png":  ".png",
			},
			Jwt: config.JwtConfig{
				JwtCookie: "YouNoteJWT",
			},
			Csrf: config.CsrfConfig{
				CsrfCookie:   "YouNoteCSRF",
				CSRFLifeTime: time.Duration(24 * time.Hour),
			},
		},
		Validation: config.ValidationConfig{
			MinUsernameLength:    4,
			MaxUsernameLength:    12,
			MinPasswordLength:    8,
			MaxPasswordLength:    20,
			PasswordAllowedExtra: "$%&#",
			SecretLength:         6,
		},
	}
	userID := uuid.FromStringOrNil("ac6966bc-3c26-45a0-963e-b168fc34fd79")
	username := "user2"
	type args struct {
		userID   uuid.UUID
		username string
	}
	tests := []struct {
		wantStatus int
		name       string
		ucMocker   func(ctx context.Context, uc *mock_auth.MockAuthUsecase, blockerUc *mock_auth.MockBlockerUsecase)
		args       args
	}{
		{
			name: "Test_Success",
			ucMocker: func(ctx context.Context, uc *mock_auth.MockAuthUsecase, blockerUc *mock_auth.MockBlockerUsecase) {
				uc.EXPECT().GenerateAndUpdateSecret(ctx, username).Return([]byte{}, nil)
			},
			args: args{
				userID:   userID,
				username: username,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Test_BadRequest",
			ucMocker: func(ctx context.Context, uc *mock_auth.MockAuthUsecase, blockerUc *mock_auth.MockBlockerUsecase) {
				uc.EXPECT().GenerateAndUpdateSecret(ctx, username).Return([]byte{}, errors.New(""))
			},
			args: args{
				userID:   userID,
				username: username,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Test_Unauthorized",
			ucMocker: func(ctx context.Context, uc *mock_auth.MockAuthUsecase, blockerUc *mock_auth.MockBlockerUsecase) {

			},
			args: args{
				userID:   userID,
				username: username,
			},
			wantStatus: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			uc := mock_auth.NewMockAuthUsecase(ctrl)
			blockerUc := mock_auth.NewMockBlockerUsecase(ctrl)
			defer ctrl.Finish()
			req := httptest.NewRequest("GET", "http://example.com/api/handler/", bytes.NewBufferString(""))
			w := httptest.NewRecorder()
			ctx := context.WithValue(req.Context(), config.PayloadContextKey, models.JwtPayload{Id: userID, Username: username})
			req = req.WithContext(ctx)
			if tt.name == "Test_Unauthorized" {
				req = req.WithContext(context.Background())
			}

			tt.ucMocker(req.Context(), uc, blockerUc)

			h := CreateAuthHandler(uc, blockerUc, testLogger, testConfig.AuthHandler, testConfig.Validation)

			h.GetQRCode(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
