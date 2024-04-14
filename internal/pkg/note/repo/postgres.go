package repo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

const (
	getAllNotes = "SELECT id, note_data, created_at, updated_at, owner_id FROM notes WHERE owner_id = $1 AND LOWER(note_data->>'title') LIKE LOWER($2) ORDER BY update_time DESC LIMIT $3 OFFSET $4;"
	getNote     = "SELECT id, note_data, created_at, updated_at, owner_id FROM notes WHERE id = $1;"
	createNote  = "INSERT INTO notes(id, note_data, created_at, updated_at, owner_id) VALUES ($1, $2, $3, $4, $5);"
	updateNote  = "UPDATE notes SET note_data = $1, updated_at = $2 WHERE id = $3; "
	deleteNote  = "DELETE FROM notes WHERE id = $1;"
)

type NoteRepo struct {
	db     pgxtype.Querier
	logger *slog.Logger
}

func CreateNoteRepo(db pgxtype.Querier, logger *slog.Logger) *NoteRepo {
	return &NoteRepo{
		db:     db,
		logger: logger,
	}
}

func (repo *NoteRepo) ReadAllNotes(ctx context.Context, userId uuid.UUID, count int64, offset int64, titleSubstr string) ([]models.Note, error) {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	result := make([]models.Note, 0, count)

	query, err := repo.db.Query(ctx, getAllNotes, userId, "%"+titleSubstr+"%", count, offset)
	if err != nil {
		logger.Error(err.Error())
		return result, err
	}

	for query.Next() {
		var note models.Note
		if err := query.Scan(&note.Id, &note.Data, &note.CreateTime, &note.UpdateTime, &note.OwnerId); err != nil {
			logger.Error(err.Error())
			return result, fmt.Errorf("error occured while scanning notes:%w", err)
		}
		result = append(result, note)
	}

	logger.Info("success")
	return result, nil
}

func (repo *NoteRepo) ReadNote(ctx context.Context, noteId uuid.UUID) (models.Note, error) {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	resultNote := models.Note{}

	err := repo.db.QueryRow(ctx, getNote, noteId).Scan(
		&resultNote.Id,
		&resultNote.Data,
		&resultNote.CreateTime,
		&resultNote.UpdateTime,
		&resultNote.OwnerId,
	)

	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	logger.Info("success")
	return resultNote, nil
}

func (repo *NoteRepo) CreateNote(ctx context.Context, note models.Note) error {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	_, err := repo.db.Exec(ctx, createNote, note.Id, string(note.Data), note.CreateTime, note.UpdateTime, note.OwnerId)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NoteRepo) UpdateNote(ctx context.Context, note models.Note) error {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	_, err := repo.db.Exec(ctx, updateNote, note.Data, note.UpdateTime, note.Id)

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil

}
func (repo *NoteRepo) DeleteNote(ctx context.Context, id uuid.UUID) error {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	_, err := repo.db.Exec(ctx, deleteNote, id)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}
