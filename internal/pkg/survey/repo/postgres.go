package repo

import (
	"context"
	"log/slog"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/jackc/pgtype/pgxtype"
)

const (
	addResult = "INSERT INTO results(id, question_id, voice) VALUES ($1, $2, $3);"
)

type SurveyRepo struct {
	db pgxtype.Querier
}

func CreateSurveyRepo(db pgxtype.Querier) *SurveyRepo {
	return &SurveyRepo{
		db: db,
	}
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
