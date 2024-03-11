package repo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
)

const (
	getAllNotes = "SELECT id, data, create_time, update_time, owner_id FROM notes WHERE owner_id = $1 AND LOWER(data->>'title') LIKE LOWER($2) ORDER BY COALESCE(update_time, create_time) DESC LIMIT $3 OFFSET $4;"
	getNote     = "SELECT id, data, create_time, update_time, owner_id FROM notes WHERE id = $1;"
	createNote  = "INSERT INTO notes(id, data, create_time, update_time, owner_id) VALUES ($1, $2::json, $3, $4, $5);"
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
	result := make([]models.Note, 0, count)
	repo.logger.Info(utils.GFN())

	query, err := repo.db.Query(ctx, getAllNotes, userId, "%"+titleSubstr+"%", count, offset)
	if err != nil {
		repo.logger.Error(err.Error())
		return result, err
	}

	for query.Next() {
		var note models.Note
		err := query.Scan(&note.Id, &note.Data, &note.CreateTime, &note.UpdateTime, &note.OwnerId)
		if err != nil {
			repo.logger.Error(err.Error())
			return result, fmt.Errorf("error occured while scanning notes:%w", err)
		}
		result = append(result, note)
	}
	return result, nil
}

func (repo *NoteRepo) ReadNote(ctx context.Context, noteId uuid.UUID) (models.Note, error) {
	resultNote := models.Note{}
	repo.logger.Info(utils.GFN())

	err := repo.db.QueryRow(ctx, getNote, noteId).Scan(
		&resultNote.Id,
		&resultNote.Data,
		&resultNote.CreateTime,
		&resultNote.UpdateTime,
		&resultNote.OwnerId,
	)

	if err != nil {
		repo.logger.Error(err.Error())
		return models.Note{}, err
	}

	return resultNote, nil
}

func (repo *NoteRepo) CreateNote(ctx context.Context, note models.Note) error {
	repo.logger.Info(utils.GFN())
	_, err := repo.db.Exec(ctx, createNote, note.Id, note.Data, note.CreateTime, note.UpdateTime, note.OwnerId)
	if err != nil {
		repo.logger.Error(err.Error())
	}
	return err
}
