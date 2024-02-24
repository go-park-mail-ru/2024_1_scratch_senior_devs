package auth

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"

	"github.com/satori/uuid"
)

type AuthUsecase interface {
	SignUp(context.Context, *models.UserFormData) (*models.User, error)
}

type AuthRepo interface {
	CreateUser(context.Context, *models.User) error
	GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}
