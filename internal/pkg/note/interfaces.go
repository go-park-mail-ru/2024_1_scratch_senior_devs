package note

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/gofrs/uuid"
)

type NoteUsecase interface {
	GetAllNotes(context.Context, uuid.UUID, int64, int64) ([]models.Note, error) //здесь передаем айди юзера
	GetNote(context.Context, uuid.UUID) (models.Note, error)                     //здесь передаем айди заметки
}

type NoteRepo interface {
	ReadAllNotes(context.Context, uuid.UUID, int64, int64) ([]models.Note, error) //тут айди юзера передаем
	ReadNote(context.Context, uuid.UUID) (models.Note, error)                     //тут айди заметки передаем
}
