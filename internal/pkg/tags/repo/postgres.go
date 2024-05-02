package repo

import (
	"context"
	"log/slog"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/satori/uuid"
)

const (
	addTag    = "INSERT INTO note_tag(note_id, tag_name) VALUES ($1, $2);"
	deleteTag = "DELETE FROM note_tag WHERE note_id = $1 AND tag_name = $2;"
)

type TagRepo struct {
	db pgxtype.Querier
}

func CreateTagRepo(db pgxtype.Querier) *TagRepo {
	return &TagRepo{
		db: db,
	}
}

func (repo *TagRepo) AddTag(ctx context.Context, tagName string, noteId uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	_, err := repo.db.Exec(ctx, addTag, noteId, tagName)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *TagRepo) DeleteTag(ctx context.Context, tagName string, noteId uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	_, err := repo.db.Exec(ctx, deleteTag, noteId, tagName)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}
