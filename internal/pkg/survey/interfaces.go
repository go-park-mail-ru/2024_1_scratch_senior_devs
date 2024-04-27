package survey

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type SurveyUsecase interface {
	GetSurvey(ctx context.Context) ([]models.Question, error)
	Vote(ctx context.Context, vote models.Vote) error
}

type SurveyRepo interface {
	GetSurvey(ctx context.Context) ([]models.Question, error)
	AddResult(ctx context.Context, res models.Result) error
}
