package survey

import (
	"context"
	"time"

	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type SurveyUsecase interface {
	GetSurvey(ctx context.Context) ([]models.Question, error)
	Vote(ctx context.Context, vote models.Vote) error
	CreateSurvey(ctx context.Context, questions []models.Question) error
	GetStats(ctx context.Context) ([]models.Stat, error)
}

type SurveyRepo interface {
	GetSurvey(ctx context.Context) ([]models.Question, error)
	AddResult(ctx context.Context, res models.Result) error
	AddSurvey(ctx context.Context, surveyID uuid.UUID, createTime time.Time) error
	AddQuestion(ctx context.Context, question models.Question) error
	GetStats(ctx context.Context) ([]models.Stat, error)
}
