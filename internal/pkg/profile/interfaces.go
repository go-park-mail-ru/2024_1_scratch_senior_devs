package profile

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/satori/uuid"
	"mime/multipart"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type ProfileUsecase interface {
	UpdateProfile(context.Context, uuid.UUID, models.ProfileUpdatePayload) (models.User, error)
	UpdateProfileAvatar(context.Context, uuid.UUID, multipart.File) (models.User, error)
}

type ProfileRepo interface {
	UpdateProfile(context.Context, models.User) error
	UpdateProfileAvatar(context.Context, uuid.UUID, string) error
	GetUserById(context.Context, uuid.UUID) (models.User, error)
}
