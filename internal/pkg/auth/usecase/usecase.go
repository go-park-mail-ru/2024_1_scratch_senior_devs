package usecase

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
	"log/slog"
	"time"

	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth"
)

const (
	JWTLifeTime      = 24 * time.Hour
	defaultImagePath = "default.jpg"
)

type AuthUsecase struct {
	repo   auth.AuthRepo
	logger *slog.Logger
}

func CreateAuthUsecase(repo auth.AuthRepo, logger *slog.Logger) *AuthUsecase {
	return &AuthUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *AuthUsecase) SignUp(ctx context.Context, data models.UserFormData) (models.User, string, time.Time, error) {
	logger := uc.logger.With(slog.String("ID", utils.GetRequestId(ctx)), slog.String("func", utils.GFN()))

	currentTime := time.Now().UTC()
	expTime := currentTime.Add(JWTLifeTime)

	newUser := models.User{
		Id:           uuid.NewV4(),
		Username:     data.Username,
		PasswordHash: utils.GetHash(data.Password),
		ImagePath:    defaultImagePath,
		CreateTime:   currentTime,
	}

	err := uc.repo.CreateUser(ctx, newUser)
	if err != nil {
		logger.Error(err.Error())
		return models.User{}, "", currentTime, err
	}

	token, err := middleware.GenToken(newUser, JWTLifeTime)
	if err != nil {
		logger.Error("middleware.GenToken error: " + err.Error())
		return models.User{}, "", currentTime, err
	}

	logger.Info("success")
	return newUser, token, expTime, nil
}

func (uc *AuthUsecase) SignIn(ctx context.Context, data models.UserFormData) (models.User, string, time.Time, error) {
	logger := uc.logger.With(slog.String("ID", utils.GetRequestId(ctx)), slog.String("func", utils.GFN()))

	currentTime := time.Now().UTC()
	expTime := currentTime.Add(JWTLifeTime)

	user, err := uc.repo.GetUserByUsername(ctx, data.Username)
	if err != nil {
		logger.Error(err.Error())
		return models.User{}, "", currentTime, err
	}
	if user.PasswordHash != utils.GetHash(data.Password) {
		logger.Error("wrong password: " + err.Error())
		return models.User{}, "", currentTime, errors.New("wrong username or password")
	}

	token, err := middleware.GenToken(user, JWTLifeTime)
	if err != nil {
		logger.Error("middleware.GenToken error: " + err.Error())
		return models.User{}, "", currentTime, err
	}

	logger.Info("success")
	return user, token, expTime, nil
}

func (uc *AuthUsecase) CheckUser(ctx context.Context, id uuid.UUID) (models.User, error) {
	logger := uc.logger.With(slog.String("ID", utils.GetRequestId(ctx)), slog.String("func", utils.GFN()))

	userData, err := uc.repo.GetUserById(ctx, id)
	if err != nil {
		logger.Error(err.Error())
		return models.User{}, err
	}

	logger.Info("success")
	return userData, nil
}
