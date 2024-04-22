package grpcw

import (
	"bytes"
	"context"
	"log/slog"

	generatedAuth "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth"
	"github.com/satori/uuid"
)

func getUser(user models.User) *generatedAuth.User {
	return &generatedAuth.User{
		Id:           user.Id.String(),
		Description:  user.Description,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		CreateTime:   user.CreateTime.String(),
		ImagePath:    user.ImagePath,
		SecondFactor: string(user.SecondFactor),
	}
}

type GrpcAuthHandler struct {
	generatedAuth.AuthServer
	uc auth.AuthUsecase
}

func NewGrpcAuthHandler(uc auth.AuthUsecase) *GrpcAuthHandler {
	return &GrpcAuthHandler{
		uc: uc,
	}
}

func (h *GrpcAuthHandler) SignIn(ctx context.Context, in *generatedAuth.UserFormData) (*generatedAuth.SignInResponse, error) {
	payload := models.UserFormData{
		Username: in.Username,
		Password: in.Password,
		Code:     in.Code,
	}

	user, jwtToken, expTime, err := h.uc.SignIn(ctx, payload)
	if err != nil {
		return nil, err
	}

	return &generatedAuth.SignInResponse{
		User:    getUser(user),
		Token:   jwtToken,
		Expires: expTime.String(),
	}, nil
}

func (h *GrpcAuthHandler) SignUp(ctx context.Context, in *generatedAuth.UserFormData) (*generatedAuth.SignUpResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	payload := models.UserFormData{
		Username: in.Username,
		Password: in.Password,
		Code:     in.Code,
	}

	user, jwtToken, expTime, err := h.uc.SignUp(ctx, payload)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("success")
	return &generatedAuth.SignUpResponse{
		User:    getUser(user),
		Token:   jwtToken,
		Expires: expTime.String(),
	}, nil
}

func (h *GrpcAuthHandler) CheckUser(ctx context.Context, in *generatedAuth.CheckUserRequest) (*generatedAuth.User, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	userID := uuid.FromStringOrNil(in.UserId)

	user, err := h.uc.CheckUser(ctx, userID)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("success")
	return getUser(user), nil
}

func (h *GrpcAuthHandler) UpdateProfile(ctx context.Context, in *generatedAuth.UpdateProfileRequest) (*generatedAuth.User, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	userID := uuid.FromStringOrNil(in.UserId)

	user, err := h.uc.UpdateProfile(ctx, userID, models.ProfileUpdatePayload{
		Description: in.Payload.Description,
		Password: models.Passwords{
			Old: in.Payload.Password.Old,
			New: in.Payload.Password.New,
		},
	})
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("success")
	return getUser(user), nil
}

func (h *GrpcAuthHandler) UpdateProfileAvatar(ctx context.Context, in *generatedAuth.UpdateProfileAvatarRequest) (*generatedAuth.User, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	userID := uuid.FromStringOrNil(in.UserId)
	avatarReadSeeker := bytes.NewReader(in.Avatar)

	user, err := h.uc.UpdateProfileAvatar(ctx, userID, avatarReadSeeker, in.Extension)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("success")
	return getUser(user), nil
}

func (h *GrpcAuthHandler) GenerateAndUpdateSecret(ctx context.Context, in *generatedAuth.SecretRequest) (*generatedAuth.GenerateAndUpdateSecretResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	secret, err := h.uc.GenerateAndUpdateSecret(ctx, in.Username)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("success")
	return &generatedAuth.GenerateAndUpdateSecretResponse{
		Secret: secret,
	}, nil
}

func (h *GrpcAuthHandler) DeleteSecret(ctx context.Context, in *generatedAuth.SecretRequest) (*generatedAuth.EmptyMessage, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	if err := h.uc.DeleteSecret(ctx, in.Username); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("success")
	return &generatedAuth.EmptyMessage{}, nil
}
