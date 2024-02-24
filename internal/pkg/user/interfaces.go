package user

import (
	"context"
	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

type UsersRepo interface {
	CreateUser(context.Context, *models.User) error
	GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}
