package usecase

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/profile"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
	"github.com/satori/uuid"
)

type ProfileUsecase struct {
	repo profile.ProfileRepo
}

func CreateProfileUsecase(repo profile.ProfileRepo) *ProfileUsecase {
	return &ProfileUsecase{
		repo: repo,
	}
}

func (uc *ProfileUsecase) UpdateProfile(ctx context.Context, userID uuid.UUID, payload models.ProfileUpdatePayload) (models.User, error) {
	payload.Sanitize()
	if err := payload.Validate(); err != nil {
		return models.User{}, err
	}

	user, err := uc.repo.GetUserById(ctx, userID)
	if err != nil {
		return models.User{}, err
	}

	if payload.Password.Old != "" && payload.Password.New != "" {
		if user.PasswordHash != utils.GetHash(payload.Password.Old) {
			return models.User{}, errors.New("wrong password")
		} else {
			user.PasswordHash = utils.GetHash(payload.Password.New)
		}
	}

	if payload.Description != "" {
		user.Description = payload.Description
	}

	err = uc.repo.UpdateProfile(ctx, user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
