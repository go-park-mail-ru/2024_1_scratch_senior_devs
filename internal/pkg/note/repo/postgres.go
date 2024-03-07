package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

const (
	getAllNotes = "SELECT id, data, create_time, update_time, owner_id FROM notes WHERE owner_id = $1 AND LOWER(data->>'title') LIKE LOWER($2) LIMIT $3 OFFSET $4;"
	createNote  = "INSERT INTO notes(id, data, create_time, update_time, owner_id) VALUES ($1, $2::json, $3, $4, $5);"
)

type NoteRepo struct {
	db pgxtype.Querier
}

func CreateNoteRepo(db pgxtype.Querier) *NoteRepo {
	return &NoteRepo{db: db}
}

func (repo *NoteRepo) ReadAllNotes(ctx context.Context, userId uuid.UUID, count int64, offset int64, titleSubstr string) ([]models.Note, error) {
	result := make([]models.Note, 0, count)

	query, err := repo.db.Query(ctx, getAllNotes, userId, "%"+titleSubstr+"%", count, offset)
	if err != nil {
		return result, err
	}

	for query.Next() {
		var note models.Note
		err := query.Scan(&note.Id, &note.Data, &note.CreateTime, &note.UpdateTime, &note.OwnerId)
		if err != nil {
			return result, fmt.Errorf("error occured while scanning notes:%w", err)
		}
		result = append(result, note)
	}
	return result, nil
}

func (repo *NoteRepo) CreateNote(ctx context.Context, note models.Note) error {
	_, err := repo.db.Exec(ctx, createNote, note.Id, note.Data, note.CreateTime, note.UpdateTime, note.OwnerId)
	return err
}
