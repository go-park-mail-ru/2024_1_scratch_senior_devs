package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
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
	repo auth.AuthRepo
}

func CreateAuthUsecase(repo auth.AuthRepo) *AuthUsecase {
	return &AuthUsecase{
		repo: repo,
	}
}

func (uc *AuthUsecase) SignUp(ctx context.Context, data *models.UserFormData) (*models.User, string, time.Time, error) {
	currentTime := time.Now().UTC()
	expTime := currentTime.Add(JWTLifeTime)

	newUser := &models.User{
		Id:           uuid.NewV4(),
		Username:     data.Username,
		PasswordHash: utils.GetHash(data.Password),
		ImagePath:    defaultImagePath,
		CreateTime:   currentTime,
	}

	err := uc.repo.CreateUser(ctx, newUser)
	if err != nil {
		return &models.User{}, "", currentTime, err
	}

	token, _ := middleware.GenToken(newUser, JWTLifeTime)

	return newUser, token, expTime, nil
}

func (uc *AuthUsecase) SignIn(ctx context.Context, data *models.UserFormData) (*models.User, string, time.Time, error) {
	currentTime := time.Now().UTC()
	expTime := currentTime.Add(JWTLifeTime)

	user, err := uc.repo.CheckUserCredentials(ctx, data.Username, utils.GetHash(data.Password))
	if err != nil {
		return &models.User{}, "", currentTime, err
	}

	token, _ := middleware.GenToken(user, JWTLifeTime)

	return user, token, expTime, nil
}

func (uc *AuthUsecase) CheckUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	userData, err := uc.repo.GetUserById(ctx, id)
	if err != nil {
		return &models.User{}, err
	}

	return userData, nil
}
