package survey

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type SurveyUsecase interface {
}

type SurveyRepo interface {
	AddResult(ctx context.Context, res models.Result) error
}
