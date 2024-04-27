package usecase

import (
	"context"
	"errors"
	"log/slog"
	"time"

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

func (uc *SurveyUsecase) CreateSurvey(ctx context.Context, questions models.CreateSurveyRequest) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	if err := uc.repo.AddSurvey(ctx, uuid.NewV4(), time.Now().UTC()); err != nil {
		logger.Error(err.Error())
		return err
	}

	for i, question := range questions.Questions {
		if err := uc.repo.AddQuestion(ctx, models.Question{
			Id:           uuid.NewV4(),
			Title:        question.Title,
			QuestionType: question.QuestionType,
			Number:       i + 1,
		}); err != nil {
			logger.Error(err.Error())
			return err
		}
	}

	logger.Info("success")
	return nil
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
	if vote.Voice < 0 || vote.Voice > 10 {
		return errors.New("not acceptable")
	}
	newResult := models.Result{
		Id:         uuid.NewV4(),
		QuestionId: vote.QuestionId,
		Voice:      vote.Voice,
	}

	if err := uc.repo.AddResult(ctx, newResult); err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}

func (uc *SurveyUsecase) GetStats(ctx context.Context) ([]models.Stat, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	result, err := uc.repo.GetStats(ctx)
	if err != nil {
		logger.Error(err.Error())
		return []models.Stat{}, err
	}

	logger.Info("success")
	return result, nil
}
