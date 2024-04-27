package repo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/satori/uuid"
)

const (
	addResult = "INSERT INTO results(id, question_id, voice) VALUES ($1, $2, $3);"

	getSurvey = "SELECT id, title, question_type, number, survey_id FROM questions WHERE survey_id = $1 ORDER BY number ASC;"
)

type SurveyRepo struct {
	db pgxtype.Querier
}

func CreateSurveyRepo(db pgxtype.Querier) *SurveyRepo {
	return &SurveyRepo{
		db: db,
	}
}

func (repo *SurveyRepo) GetSurvey(ctx context.Context, id uuid.UUID) ([]models.Question, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	result := make([]models.Question, 0)

	query, err := repo.db.Query(ctx, getSurvey, id)
	if err != nil {
		logger.Error(err.Error())
		return result, err
	}

	for query.Next() {
		var question models.Question
		if err := query.Scan(&question.Id, &question.Title, &question.QuestionType, &question.Number, &question.SurveyId); err != nil {
			logger.Error(err.Error())
			return result, fmt.Errorf("error occured while scanning questions: %w", err)
		}
		result = append(result, question)
	}

	logger.Info("success")
	return result, nil
}

func (repo *SurveyRepo) AddResult(ctx context.Context, res models.Result) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	_, err := repo.db.Exec(ctx, addResult, res.Id, res.QuestionId, res.Voice)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil

}
