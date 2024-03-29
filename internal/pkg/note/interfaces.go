package note

import (
	"context"
	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type NoteUsecase interface {
	GetAllNotes(context.Context, uuid.UUID, int64, int64, string) ([]models.Note, error)
	GetNote(context.Context, uuid.UUID, uuid.UUID) (models.Note, error)
	CreateNote(context.Context, uuid.UUID, []byte) (models.Note, error)
	UpdateNote(context.Context, uuid.UUID, uuid.UUID, []byte) (models.Note, error)
	DeleteNote(context.Context, uuid.UUID, uuid.UUID) error
}

type NoteRepo interface {
	ReadAllNotes(context.Context, uuid.UUID, int64, int64, string) ([]models.Note, error)
	ReadNote(context.Context, uuid.UUID) (models.Note, error)
	CreateNote(context.Context, models.Note) error
	UpdateNote(context.Context, models.Note) error
	DeleteNote(context.Context, uuid.UUID) error
}
