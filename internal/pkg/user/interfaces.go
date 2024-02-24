package user

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

type UsersRepo interface {
	CreateUser(context.Context, *models.User) error
}
