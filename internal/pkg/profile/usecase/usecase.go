package usecase

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/profile"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
	"github.com/satori/uuid"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
)

type ProfileUsecase struct {
	repo   profile.ProfileRepo
	logger *slog.Logger
}

func CreateProfileUsecase(repo profile.ProfileRepo, logger *slog.Logger) *ProfileUsecase {
	return &ProfileUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *ProfileUsecase) UpdateProfile(ctx context.Context, userID uuid.UUID, payload models.ProfileUpdatePayload) (models.User, error) {
	logger := uc.logger.With(slog.String("ID", utils.GetRequestId(ctx)), slog.String("func", utils.GFN()))

	payload.Sanitize()

	user, err := uc.repo.GetUserById(ctx, userID)
	if err != nil {
		logger.Error(err.Error())
		return models.User{}, err
	}

	if payload.Password.Old != "" && payload.Password.New != "" {
		if err := payload.Validate(); err != nil {
			logger.Error("validation error: " + err.Error())
			return models.User{}, err
		}

		if user.PasswordHash != utils.GetHash(payload.Password.Old) {
			logger.Error("wrong password: " + err.Error())
			return models.User{}, errors.New("wrong password")
		}

		user.PasswordHash = utils.GetHash(payload.Password.New)
	}

	if payload.Description != "" {
		user.Description = payload.Description
	}

	if err := uc.repo.UpdateProfile(ctx, user); err != nil {
		logger.Error(err.Error())
		return models.User{}, err
	}

	logger.Info("success")
	return user, nil
}

func (uc *ProfileUsecase) UpdateProfileAvatar(ctx context.Context, userID uuid.UUID, avatar multipart.File) (models.User, error) {
	logger := uc.logger.With(slog.String("ID", utils.GetRequestId(ctx)), slog.String("func", utils.GFN()))

	user, err := uc.repo.GetUserById(ctx, userID)
	if err != nil {
		logger.Error(err.Error())
		return models.User{}, err
	}

	imagesBasePath := os.Getenv("IMAGES_BASE_PATH")

	if user.ImagePath != "default.jpg" {
		if err := os.Remove(imagesBasePath + user.ImagePath); err != nil {
			logger.Error(err.Error())
			return models.User{}, err
		}
	}

	newImagePath := uuid.NewV4().String()
	file, err := os.Create(imagesBasePath + newImagePath)
	if err != nil {
		logger.Error(err.Error())
		return models.User{}, err
	}
	defer file.Close()

	_, err = avatar.Seek(0, 0)
	if err != nil {
		logger.Error(err.Error())
		return models.User{}, err
	}
	_, err = io.Copy(file, avatar)
	if err != nil {
		logger.Error(err.Error())
		return models.User{}, err
	}

	if err := uc.repo.UpdateProfileAvatar(ctx, userID, newImagePath); err != nil {
		logger.Error(err.Error())
		return models.User{}, err
	}

	user.ImagePath = newImagePath

	logger.Info("success")
	return user, nil
}
