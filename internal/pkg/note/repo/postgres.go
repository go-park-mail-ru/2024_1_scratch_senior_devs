package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

const (
	getAllNotes = "SELECT id, data, create_time, update_time, owner_id FROM notes WHERE owner_id= $1 LIMIT $2 OFFSET $3;"
)

type NoteRepo struct {
	db pgxtype.Querier
}

func CreateNoteRepo(db pgxtype.Querier) *NoteRepo {
	return &NoteRepo{db: db}
}

func (repo *NoteRepo) ReadAllNotes(ctx context.Context, userId uuid.UUID, count int64, offset int64) (result []models.Note, err error) {
	query, err := repo.db.Query(ctx, getAllNotes, userId, count, offset)
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
