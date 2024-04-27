package usecase

import (
	"context"
	"log/slog"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/survey"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/satori/uuid"
)

type SurveyUsecase struct {
	repo survey.SurveyRepo
}

func CreateSurveyUsecase(repo survey.SurveyRepo) *SurveyUsecase {
	return &SurveyUsecase{
		repo: repo,
	}
}

func (uc *SurveyUsecase) GetSurvey(ctx context.Context) ([]models.Question, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	result, err := uc.repo.GetSurvey(ctx)
	if err != nil {
		logger.Error(err.Error())
		return []models.Question{}, err
	}

	logger.Info("success")
	return result, nil
}

func (uc *SurveyUsecase) Vote(ctx context.Context, vote models.Vote) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	newResult := models.Result{
		Id: uuid.NewV4(),
	}

	if err := uc.repo.AddResult(ctx, newResult); err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}
