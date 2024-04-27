package repo

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/jackc/pgtype/pgxtype"
)

const (
	addResult   = "INSERT INTO results(id, question_id, voice) VALUES ($1, $2, $3);"
	addSurvey   = "INSERT INTO surveys(id, created_at) VALUES ($1, $2);"
	addQuestion = "INSERT INTO questions(id, title, question_type, number, survey_id) VALUES ($1, $2, $3, $4, $5);"

	getSurvey = "SELECT id, title, question_type, number, survey_id FROM questions WHERE survey_id = (SELECT id FROM surveys ORDER BY created_at DESC LIMIT 1) ORDER BY number ASC;"

	getStats = `
	SELECT q.id, q.title, q.question_type, r.voice, COUNT(r.voice) from results r
	JOIN questions q on q.id = r.question_id 
	GROUP BY  r.voice, q.id;
	
	`
)

type SurveyRepo struct {
	db pgxtype.Querier
}

func CreateSurveyRepo(db pgxtype.Querier) *SurveyRepo {
	return &SurveyRepo{
		db: db,
	}
}

func (repo *SurveyRepo) AddSurvey(ctx context.Context, surveyID uuid.UUID, createTime time.Time) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	_, err := repo.db.Exec(ctx, addSurvey, surveyID, createTime)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *SurveyRepo) AddQuestion(ctx context.Context, question models.Question) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	_, err := repo.db.Exec(ctx, addQuestion, question.Id, question.Title, question.QuestionType, question.Number, question.SurveyId)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *SurveyRepo) GetSurvey(ctx context.Context) ([]models.Question, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	result := make([]models.Question, 0)

	query, err := repo.db.Query(ctx, getSurvey)
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

func (repo *SurveyRepo) GetStats(ctx context.Context) ([]models.Stat, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	result := make([]models.Stat, 0)

	query, err := repo.db.Query(ctx, getStats)
	if err != nil {
		logger.Error(err.Error())
		return result, err
	}

	for query.Next() {
		var stat models.Stat
		if err := query.Scan(&stat.QuestionId, &stat.Title, &stat.QuestionType, &stat.Voice, &stat.Count); err != nil {
			logger.Error(err.Error())
			return result, fmt.Errorf("error occured while scanning stats: %w", err)
		}
		result = append(result, stat)
	}
	logger.Info("success")
	return result, nil
}
