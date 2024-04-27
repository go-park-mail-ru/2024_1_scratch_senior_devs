package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/survey"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/satori/uuid"
	"log/slog"
)

type SurveyUsecase struct {
	repo survey.SurveyRepo
}

func CreateSurveyUsecase(repo survey.SurveyRepo) *SurveyUsecase {
	return &SurveyUsecase{
		repo: repo,
	}
}

func (uc *SurveyUsecase) GetSurvey(ctx context.Context, id uuid.UUID) ([]models.Question, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	result, err := uc.repo.GetSurvey(ctx, id)
	if err != nil {
		logger.Error(err.Error())
		return []models.Question{}, err
	}

	logger.Info("success")
	return result, nil
}
