package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/jwt"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/code"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/images"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/request"
	"io"
	"log/slog"
	"os"
	"path"
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
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	currentTime := time.Now().UTC()
	expTime := currentTime.Add(JWTLifeTime)

	newUser := models.User{
		Id:           uuid.NewV4(),
		Username:     data.Username,
		PasswordHash: request.GetHash(data.Password),
		ImagePath:    defaultImagePath,
		CreateTime:   currentTime,
		SecondFactor: "",
	}

	err := uc.repo.CreateUser(ctx, newUser)
	if err != nil {
		logger.Error(err.Error())
		return models.User{}, "", currentTime, auth.ErrCreatingUser
	}

	token, err := jwt.GenToken(newUser, JWTLifeTime)
	if err != nil {
		logger.Error("middleware.GenToken error: " + err.Error())
		return models.User{}, "", currentTime, err
	}

	logger.Info("success")
	return newUser, token, expTime, nil
}

func (uc *AuthUsecase) SignIn(ctx context.Context, data models.UserFormData) (models.User, string, time.Time, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	currentTime := time.Now().UTC()
	expTime := currentTime.Add(JWTLifeTime)

	user, err := uc.repo.GetUserByUsername(ctx, data.Username)
	if err != nil {
		// тут increase счетчика
		logger.Error(err.Error())
		return models.User{}, "", currentTime, auth.ErrUserNotFound
	}
	if user.PasswordHash != request.GetHash(data.Password) {
		// тут increase счетчика
		logger.Error("wrong password")
		return models.User{}, "", currentTime, auth.ErrWrongUserData
	}

	if user.SecondFactor != "" {
		if data.Code == "" {
			logger.Error(auth.ErrFirstFactorPassed.Error())
			return models.User{}, "", currentTime, auth.ErrFirstFactorPassed
		}

		err := code.CheckCode(data.Code, string(user.SecondFactor))
		if err != nil {
			// тут increase счетчика
			logger.Error(err.Error())
			return models.User{}, "", currentTime, auth.ErrWrongAuthCode
		}
	}

	token, err := jwt.GenToken(user, JWTLifeTime)
	if err != nil {
		logger.Error("middleware.GenToken error: " + err.Error())
		return models.User{}, "", currentTime, err
	}

	logger.Info("success")
	return user, token, expTime, nil
}

func (uc *AuthUsecase) CheckUser(ctx context.Context, id uuid.UUID) (models.User, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	userData, err := uc.repo.GetUserById(ctx, id)
	if err != nil {
		logger.Error(err.Error())
		return models.User{}, auth.ErrUserNotFound
	}

	logger.Info("success")
	return userData, nil
}

func (uc *AuthUsecase) UpdateProfile(ctx context.Context, userID uuid.UUID, payload models.ProfileUpdatePayload) (models.User, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	payload.Sanitize()

	user, err := uc.repo.GetUserById(ctx, userID)
	if err != nil {
		logger.Error(err.Error())
		return models.User{}, auth.ErrUserNotFound
	}

	if payload.Password.Old != "" && payload.Password.New != "" {
		if err := payload.Validate(); err != nil {
			logger.Error("validation error: " + err.Error())
			return models.User{}, err
		}

		if user.PasswordHash != request.GetHash(payload.Password.Old) {
			logger.Error("wrong password")
			return models.User{}, auth.ErrWrongPassword
		}

		user.PasswordHash = request.GetHash(payload.Password.New)
	}

	user.Description = payload.Description

	if err := uc.repo.UpdateProfile(ctx, user); err != nil {
		logger.Error(err.Error())
		return models.User{}, err
	}

	logger.Info("success")
	return user, nil
}

func (uc *AuthUsecase) UpdateProfileAvatar(ctx context.Context, userID uuid.UUID, avatar io.ReadSeeker, extension string) (models.User, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	user, err := uc.repo.GetUserById(ctx, userID)
	if err != nil {
		logger.Error(err.Error())
		return models.User{}, auth.ErrUserNotFound
	}

	imagesBasePath := os.Getenv("IMAGES_BASE_PATH")
	newImagePath := uuid.NewV4().String() + extension

	if err := images.WriteAvatarOnDisk(path.Join(imagesBasePath, newImagePath), avatar); err != nil {
		logger.Error("write on disk: " + err.Error())
		return models.User{}, err
	}

	if err := uc.repo.UpdateProfileAvatar(ctx, userID, newImagePath); err != nil {
		logger.Error(err.Error())
		return models.User{}, err
	}

	// удаление старой аватарки делаем только после успешного создания новой
	if user.ImagePath != "default.jpg" {
		if err := os.Remove(path.Join(imagesBasePath, user.ImagePath)); err != nil {
			logger.Error(err.Error())
		}
	}

	user.ImagePath = newImagePath

	logger.Info("success")
	return user, nil
}

func (uc *AuthUsecase) GenerateAndUpdateSecret(ctx context.Context, username string) ([]byte, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	secret := code.GenerateSecret()
	err := uc.repo.UpdateSecret(ctx, username, string(secret))
	if err != nil {
		logger.Error(err.Error())
		return []byte{}, err
	}

	logger.Info("success")
	return secret, nil
}
