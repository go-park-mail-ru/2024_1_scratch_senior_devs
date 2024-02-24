package usecase

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	user "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/repo"
)

type AuthUsecase struct {
	repo user.AuthRepo
}

func CreateAuthUsecase(repo user.AuthRepo) *AuthUsecase {
	return &AuthUsecase{
		repo: repo,
	}
}

func getHash(data string) string {
	hasher := md5.New()
	hasher.Write([]byte(data))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func (uc *AuthUsecase) SignUp(ctx context.Context, data *models.UserFormData) (*models.User, error) {
	newUser := &models.User{
		Id:           uuid.NewV4(),
		Username:     data.Username,
		PasswordHash: getHash(data.Password),
		ImagePath:    "default.jpg",
		CreateTime:   time.Now().UTC(),
	}

	err := uc.repo.CreateUser(ctx, newUser)
	if err != nil {
		err = fmt.Errorf("error creating user: %w", err)
		return &models.User{}, err
	}

	return newUser, nil
}
