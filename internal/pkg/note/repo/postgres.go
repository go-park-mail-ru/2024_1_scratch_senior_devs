package repo

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/satori/uuid"
)

const (
	getAllNotes = "SELECT data, create_time, update_time, owner_id FROM notes WHERE owner_id= $1 LIMIT $2 OFFSET $3;"
)

type NotesRepo struct {
	db pgxtype.Querier
}

func CreateNotesRepo(db pgxtype.Querier) *NotesRepo {
	return &NotesRepo{db: db}
}

func (repo *NotesRepo) ReadAllNotes(ctx context.Context, userId uuid.UUID, count int64, offset int64) (result []models.Note, err error) {

	query, err := repo.db.Query(ctx, getAllNotes, userId, count, offset)
	if err != nil {
		return result, fmt.Errorf("error occured while getting notes: %w", err)
	}
	for query.Next() {
		var note models.Note
		err := query.Scan(&note)
		if err != nil {
			return result, fmt.Errorf("error occured while scanning notes:%w", err)
		}
		result = append(result, note)

	}
	return result, nil

}
