package repo

import (
	"context"
	"log/slog"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/satori/uuid"
)

const (
	getAttach    = "SELECT id, path, note_id FROM attaches WHERE id = $1;"
	createAttach = "INSERT INTO attaches(id, path, note_id) VALUES ($1, $2, $3);"
	deleteAttach = "DELETE FROM attaches WHERE id = $1;"
)

type AttachRepo struct {
	db pgxtype.Querier
}

func CreateAttachRepo(db pgxtype.Querier) *AttachRepo {
	return &AttachRepo{
		db: db,
	}
}

func (repo *AttachRepo) GetAttach(ctx context.Context, id uuid.UUID) (models.Attach, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	attach := models.Attach{}

	err := repo.db.QueryRow(ctx, getAttach, id).Scan(
		&attach.Id,
		&attach.Path,
		&attach.NoteId,
	)
	if err != nil {
		logger.Error(err.Error())
		return models.Attach{}, err
	}

	logger.Info("success")
	return attach, nil
}

func (repo *AttachRepo) AddAttach(ctx context.Context, attach models.Attach) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	_, err := repo.db.Exec(ctx, createAttach, attach.Id, attach.Path, attach.NoteId)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *AttachRepo) DeleteAttach(ctx context.Context, id uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	_, err := repo.db.Exec(ctx, deleteAttach, id)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}
