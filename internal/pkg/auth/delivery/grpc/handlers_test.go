package grpcw

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/delivery/grpc/gen"
	mock_auth "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/mocks"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/responses"
	"github.com/golang/mock/gomock"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_SignUp(t *testing.T) {
	var tests = []struct {
		name        string
		requestBody string
		username    string
		password    string
		usecaseErr  error
	}{
		{
			name:        "AuthHandler_SignUp_Success",
			requestBody: `{"username":"testuser2","password":"12345678a"}`,
			username:    "testuser2",
			password:    "12345678a",
			usecaseErr:  nil,
		},
		{
			name:        "AuthHandler_SignUp_Fail",
			requestBody: `{"username":"testuser2","password":"12345678a"}`,
			username:    "testuser2",
			password:    "12345678a",
			usecaseErr:  errors.New("registration failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_auth.NewMockAuthUsecase(ctrl)
			mockBlockerUsecase := mock_auth.NewMockBlockerUsecase(ctrl)
			defer ctrl.Finish()

			mockUsecase.EXPECT().SignUp(gomock.Any(), models.UserFormData{
				Username: tt.username,
				Password: tt.password,
			}).Return(models.User{
				Id:           uuid.NewV4(),
				Description:  "",
				Username:     tt.username,
				PasswordHash: responses.GetHash(tt.password),
			}, "this_is_jwt_token", time.Now(), tt.usecaseErr)

			handler := NewGrpcAuthHandler(mockUsecase, mockBlockerUsecase)
			_, err := handler.SignUp(context.Background(), &gen.UserFormData{
				Username: tt.username,
				Password: tt.password,
			})

			if tt.usecaseErr != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAuthHandler_SignIn(t *testing.T) {
	var tests = []struct {
		name        string
		requestBody string
		username    string
		password    string
		usecaseErr  error
	}{
		{
			name:        "AuthHandler_SignIn_Success",
			requestBody: `{"username":"testuser2","password":"12345678a"}`,
			username:    "testuser2",
			password:    "12345678a",
			usecaseErr:  nil,
		},
		{
			name:        "AuthHandler_SignIn_Fail",
			requestBody: `{"username":"testuser2","password":"12345678a"}`,
			username:    "testuser2",
			password:    "12345678a",
			usecaseErr:  errors.New("registration failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_auth.NewMockAuthUsecase(ctrl)
			mockBlockerUsecase := mock_auth.NewMockBlockerUsecase(ctrl)
			defer ctrl.Finish()

			mockUsecase.EXPECT().SignIn(gomock.Any(), models.UserFormData{
				Username: tt.username,
				Password: tt.password,
			}).Return(models.User{
				Id:           uuid.NewV4(),
				Description:  "",
				Username:     tt.username,
				PasswordHash: responses.GetHash(tt.password),
			}, "this_is_jwt_token", time.Now(), tt.usecaseErr)

			handler := NewGrpcAuthHandler(mockUsecase, mockBlockerUsecase)
			_, err := handler.SignIn(context.Background(), &gen.UserFormData{
				Username: tt.username,
				Password: tt.password,
			})

			if tt.usecaseErr != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAuthHandler_CheckUser(t *testing.T) {
	var tests = []struct {
		name       string
		id         uuid.UUID
		username   string
		usecaseErr error
	}{
		{
			name:       "AuthHandler_CheckUser_Success",
			id:         uuid.NewV4(),
			username:   "testuser2",
			usecaseErr: nil,
		},
		{
			name:       "AuthHandler_CheckUser_Fail",
			id:         uuid.NewV4(),
			username:   "testuser2",
			usecaseErr: errors.New("check failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_auth.NewMockAuthUsecase(ctrl)
			mockBlockerUsecase := mock_auth.NewMockBlockerUsecase(ctrl)
			defer ctrl.Finish()

			mockUsecase.EXPECT().CheckUser(gomock.Any(), tt.id).Return(models.User{
				Id:           uuid.NewV4(),
				Description:  "",
				Username:     tt.username,
				PasswordHash: responses.GetHash("12345678a"),
			}, tt.usecaseErr)

			handler := NewGrpcAuthHandler(mockUsecase, mockBlockerUsecase)
			_, err := handler.CheckUser(context.Background(), &gen.CheckUserRequest{
				UserId: tt.id.String(),
			})

			if tt.usecaseErr != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAuthHandler_UpdateProfile(t *testing.T) {
	userID := uuid.FromStringOrNil("ac6966bc-3c26-45a0-963e-b168fc34fd79")
	username := "user"

	type args struct {
		userID   uuid.UUID
		username string
	}

	tests := []struct {
		args       args
		name       string
		usecaseErr error
	}{
		{
			args: args{
				userID:   userID,
				username: username,
			},
			name:       "Test_Success",
			usecaseErr: nil,
		},
		{
			args: args{
				userID:   userID,
				username: username,
			},
			name:       "Test_BadRequest",
			usecaseErr: errors.New("update error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_auth.NewMockAuthUsecase(ctrl)
			mockBlockerUsecase := mock_auth.NewMockBlockerUsecase(ctrl)
			defer ctrl.Finish()

			mockUsecase.EXPECT().UpdateProfile(gomock.Any(), userID, gomock.Any()).Return(models.User{
				Id:           userID,
				Description:  "slkakjckld",
				Username:     username,
				PasswordHash: responses.GetHash("12345678b"),
			}, tt.usecaseErr)

			handler := NewGrpcAuthHandler(mockUsecase, mockBlockerUsecase)
			_, err := handler.UpdateProfile(context.Background(), &gen.UpdateProfileRequest{
				UserId: tt.args.userID.String(),
				Payload: &gen.ProfileUpdatePayload{
					Description: "slkakjckld",
					Password: &gen.Passwords{
						Old: "12345678a",
						New: "12345678b",
					},
				},
			})

			if tt.usecaseErr != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAuthHandler_DisableSecondFactor(t *testing.T) {
	userID := uuid.FromStringOrNil("ac6966bc-3c26-45a0-963e-b168fc34fd79")
	username := "test"

	type args struct {
		userID   uuid.UUID
		username string
	}

	tests := []struct {
		name       string
		args       args
		usecaseErr error
	}{
		{
			name: "Test_Success",
			args: args{
				userID:   userID,
				username: username,
			},
			usecaseErr: nil,
		},
		{
			name: "Test_BadRequest",
			args: args{
				userID:   userID,
				username: username,
			},
			usecaseErr: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_auth.NewMockAuthUsecase(ctrl)
			mockBlockerUsecase := mock_auth.NewMockBlockerUsecase(ctrl)
			defer ctrl.Finish()

			mockUsecase.EXPECT().DeleteSecret(gomock.Any(), username).Return(tt.usecaseErr)

			handler := NewGrpcAuthHandler(mockUsecase, mockBlockerUsecase)
			_, err := handler.DeleteSecret(context.Background(), &gen.SecretRequest{
				Username: tt.args.username,
			})

			if tt.usecaseErr != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAuthHandler_GetQRCode(t *testing.T) {
	userID := uuid.FromStringOrNil("ac6966bc-3c26-45a0-963e-b168fc34fd79")
	username := "user2"

	type args struct {
		userID   uuid.UUID
		username string
	}

	tests := []struct {
		name       string
		args       args
		usecaseErr error
	}{
		{
			name: "Test_Success",
			args: args{
				userID:   userID,
				username: username,
			},
			usecaseErr: nil,
		},
		{
			name: "Test_BadRequest",
			args: args{
				userID:   userID,
				username: username,
			},
			usecaseErr: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_auth.NewMockAuthUsecase(ctrl)
			mockBlockerUsecase := mock_auth.NewMockBlockerUsecase(ctrl)
			defer ctrl.Finish()

			mockUsecase.EXPECT().GenerateAndUpdateSecret(gomock.Any(), username).Return([]byte{}, tt.usecaseErr)

			handler := NewGrpcAuthHandler(mockUsecase, mockBlockerUsecase)
			_, err := handler.GenerateAndUpdateSecret(context.Background(), &gen.SecretRequest{
				Username: tt.args.username,
			})

			if tt.usecaseErr != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
