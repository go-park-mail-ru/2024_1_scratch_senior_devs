package auth

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

type AuthUsecase interface {
	SignUp(context.Context, *models.UserFormData) (*models.User, error)
}
