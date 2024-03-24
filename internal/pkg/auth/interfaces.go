package auth

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

var (
	ErrCreatingUser      = errors.New("this username is already taken")
	ErrFirstFactorPassed = errors.New("first factor passed")
	ErrIncorrectPayload  = errors.New("incorrect data format")
	ErrUserNotFound      = errors.New("user not found")
	ErrWrongAuthCode     = errors.New("wrong code")
	ErrWrongFileFormat   = errors.New("incorrect file format")
	ErrWrongFilesNumber  = errors.New("incorrect form data: must be exactly one file")
	ErrWrongPassword     = errors.New("wrong password")
	ErrWrongUserData     = errors.New("wrong username or password")
)

type AuthUsecase interface {
	SignUp(context.Context, models.UserFormData) (models.User, string, time.Time, error)
	SignIn(context.Context, models.UserFormData) (models.User, string, time.Time, error)
	CheckUser(context.Context, uuid.UUID) (models.User, error)
	UpdateProfile(context.Context, uuid.UUID, models.ProfileUpdatePayload) (models.User, error)
	UpdateProfileAvatar(context.Context, uuid.UUID, io.ReadSeeker, string) (models.User, error)
	GenerateAndUpdateSecret(context.Context, string) ([]byte, error)
}

type AuthRepo interface {
	CreateUser(context.Context, models.User) error
	GetUserById(context.Context, uuid.UUID) (models.User, error)
	GetUserByUsername(context.Context, string) (models.User, error)
	UpdateProfile(context.Context, models.User) error
	UpdateProfileAvatar(context.Context, uuid.UUID, string) error
	UpdateSecret(context.Context, string, string) error
}

type BlockerUsecase interface {
	CheckLoginAttempts(context.Context, string) error
}

type BlockerRepo interface {
	GetLoginAttempts(context.Context, string) (int, error)
	IncreaseLoginAttempts(context.Context, string) error
}
