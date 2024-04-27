package survey

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/satori/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type SurveyUsecase interface {
	GetSurvey(ctx context.Context, id uuid.UUID) ([]models.Question, error)
}

type SurveyRepo interface {
	GetSurvey(ctx context.Context, id uuid.UUID) ([]models.Question, error)
}
