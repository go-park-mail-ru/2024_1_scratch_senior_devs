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
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/authmw"
)

const (
	JWTLifeTime      = 24 * time.Hour
	defaultImagePath = "default.jpg"
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

func (uc *AuthUsecase) SignUp(ctx context.Context, data *models.UserFormData) (*models.User, string, error) {
	newUser := &models.User{
		Id:           uuid.NewV4(),
		Username:     data.Username,
		PasswordHash: getHash(data.Password),
		ImagePath:    defaultImagePath,
		CreateTime:   time.Now().UTC(),
	}

	err := uc.repo.CreateUser(ctx, newUser)
	if err != nil {
		err = fmt.Errorf("error creating user: %w", err)
		return &models.User{}, "", err
	}

	token, err := authmw.GenToken(newUser, JWTLifeTime)
	if err != nil {
		err = fmt.Errorf("error generating jwt: %w", err)
		return newUser, "", err
	}

	return newUser, token, nil
}

func (uc *AuthUsecase) SignIn(ctx context.Context, data *models.UserFormData) (*models.User, string, error) {
	user, err := uc.repo.CheckUserCredentials(ctx, data.Username, getHash(data.Password))
	if err != nil {
		return &models.User{}, "", fmt.Errorf("%w", err)
	}
	token, err := authmw.GenToken(user, JWTLifeTime)
	if err != nil {
		err = fmt.Errorf("error jenerating jwt: %w", err)
		return user, "", err
	}
	return user, token, nil

}
