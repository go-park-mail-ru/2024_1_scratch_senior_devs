package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/user"
	"github.com/google/uuid"
)

type AuthUsecase struct {
	repo user.UsersRepo
}

func CreateAuthUsecase(repo user.UsersRepo) *AuthUsecase {
	return &AuthUsecase{
		repo: repo,
	}
}

func (uc *AuthUsecase) SignUp(ctx context.Context, data *models.SignUpForm) (*models.User, error) {
	newUser := &models.User{
		Id:           uuid.New(),
		Username:     data.Username,
		PasswordHash: data.Password + "abc", // типа захэшировал
		ImagePath:    "default.jpg",
		CreateTime:   time.Now(),
	}

	err := uc.repo.CreateUser(ctx, newUser)
	if err != nil {
		err = fmt.Errorf("error creating user: %w", err)
		return &models.User{}, err
	}

	return newUser, nil
}
