package repo

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/satori/uuid"
	"log/slog"
)

const (
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
