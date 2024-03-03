package auth

import (
	"context"
	"time"

	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type AuthUsecase interface {
	SignUp(context.Context, models.UserFormData) (models.User, string, time.Time, error)
	SignIn(context.Context, models.UserFormData) (models.User, string, time.Time, error)
	CheckUser(context.Context, uuid.UUID) (models.User, error)
}

type AuthRepo interface {
	CreateUser(context.Context, models.User) error
	GetUserById(context.Context, uuid.UUID) (models.User, error)
	GetUserByUsername(context.Context, string) (models.User, error)
	CheckUserCredentials(context.Context, string, string) (models.User, error)
}
