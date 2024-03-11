package usecase

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/profile"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
	"github.com/satori/uuid"
	"log/slog"
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
		} else {
			user.PasswordHash = utils.GetHash(payload.Password.New)
		}
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
