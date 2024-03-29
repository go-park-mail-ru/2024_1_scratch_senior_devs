package repo

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/jackc/pgtype/pgxtype"
	"log/slog"
)

const (
	createAttach = "INSERT INTO attaches(id, path, note_id) VALUES ($1, $2, $3);"
)

type AttachRepo struct {
	db     pgxtype.Querier
	logger *slog.Logger
}

func CreateAttachRepo(db pgxtype.Querier, logger *slog.Logger) *AttachRepo {
	return &AttachRepo{
		db:     db,
		logger: logger,
	}
}

func (repo *AttachRepo) AddAttach(ctx context.Context, attach models.Attach) error {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	_, err := repo.db.Exec(ctx, createAttach, attach.Id, attach.Path, attach.NoteId)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}
