package auth

import (
	"context"
	"io"
	"time"

	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type AuthUsecase interface {
	SignUp(context.Context, models.UserFormData) (models.User, string, time.Time, error)
	SignIn(context.Context, models.UserFormData) (models.User, string, time.Time, error)
	CheckUser(context.Context, uuid.UUID) (models.User, error)
	UpdateProfile(context.Context, uuid.UUID, models.ProfileUpdatePayload) (models.User, error)
	UpdateProfileAvatar(context.Context, uuid.UUID, io.ReadSeeker, string) (models.User, error)
}

type AuthRepo interface {
	CreateUser(context.Context, models.User) error
	GetUserById(context.Context, uuid.UUID) (models.User, error)
	GetUserByUsername(context.Context, string) (models.User, error)
	UpdateProfile(context.Context, models.User) error
	UpdateProfileAvatar(context.Context, uuid.UUID, string) error
}
