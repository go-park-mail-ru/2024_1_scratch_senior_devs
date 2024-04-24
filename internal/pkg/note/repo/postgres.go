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
	getAllNotes   = "SELECT id, data, create_time, update_time, owner_id, parent, children FROM notes WHERE parent = '00000000-0000-0000-0000-000000000000' AND owner_id = $1 ORDER BY update_time DESC LIMIT $2 OFFSET $3;"
	getNote       = "SELECT id, data, create_time, update_time, owner_id, parent, children FROM notes WHERE id = $1;"
	createNote    = "INSERT INTO notes(id, data, create_time, update_time, owner_id, parent, children) VALUES ($1, $2::json, $3, $4, $5, $6, $7::UUID[]);"
	updateNote    = "UPDATE notes SET data = $1, update_time = $2 WHERE id = $3; "
	deleteNote    = "DELETE FROM notes CASCADE WHERE id = $1;"
	addSubNote    = "UPDATE notes SET children = array_append(children, $1) WHERE id = $2;"
	removeSubNote = "UPDATE notes SET children = array_remove(children, $1) WHERE id = $2;"
)

type NotePostgres struct {
	db pgxtype.Querier
}

func CreateNotePostgres(db pgxtype.Querier) *NotePostgres {
	return &NotePostgres{
		db: db,
	}
}

func (repo *NotePostgres) ReadAllNotes(ctx context.Context, userId uuid.UUID, count int64, offset int64) ([]models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	result := make([]models.Note, 0, count)

	query, err := repo.db.Query(ctx, getAllNotes, userId, count, offset)
	if err != nil {
		logger.Error(err.Error())
		return result, err
	}

	for query.Next() {
		var note models.Note
		if err := query.Scan(&note.Id, &note.Data, &note.CreateTime, &note.UpdateTime, &note.OwnerId, &note.Parent, &note.Children); err != nil {
			logger.Error(err.Error())
			return result, fmt.Errorf("error occured while scanning notes: %w", err)
		}
		result = append(result, note)
	}

	logger.Info("success")
	return result, nil
}

func (repo *NotePostgres) ReadNote(ctx context.Context, noteId uuid.UUID) (models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	resultNote := models.Note{}

	err := repo.db.QueryRow(ctx, getNote, noteId).Scan(
		&resultNote.Id,
		&resultNote.Data,
		&resultNote.CreateTime,
		&resultNote.UpdateTime,
		&resultNote.OwnerId,
		&resultNote.Parent,
		&resultNote.Children,
	)

	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	logger.Info("success")
	return resultNote, nil
}

func (repo *NotePostgres) CreateNote(ctx context.Context, note models.Note) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	_, err := repo.db.Exec(ctx, createNote, note.Id, note.Data, note.CreateTime, note.UpdateTime, note.OwnerId, note.Parent, note.Children)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NotePostgres) UpdateNote(ctx context.Context, note models.Note) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	_, err := repo.db.Exec(ctx, updateNote, note.Data, note.UpdateTime, note.Id)

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil

}

func (repo *NotePostgres) DeleteNote(ctx context.Context, id uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	_, err := repo.db.Exec(ctx, deleteNote, id)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NotePostgres) AddSubNote(ctx context.Context, id uuid.UUID, childID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	_, err := repo.db.Exec(ctx, addSubNote, childID, id)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NotePostgres) RemoveSubNote(ctx context.Context, id uuid.UUID, childID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	_, err := repo.db.Exec(ctx, removeSubNote, childID, id)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}
