package note

import (
	"context"
	"time"

	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type NoteUsecase interface {
	GetAllNotes(context.Context, uuid.UUID, int64, int64, string, []string) ([]models.Note, error)
	GetNote(context.Context, uuid.UUID, uuid.UUID) (models.Note, error)
	CreateNote(context.Context, uuid.UUID, []byte) (models.Note, error)
	UpdateNote(context.Context, uuid.UUID, uuid.UUID, []byte) (models.Note, error)
	DeleteNote(context.Context, uuid.UUID, uuid.UUID) error

	CreateSubNote(context.Context, uuid.UUID, []byte, uuid.UUID) (models.Note, error)

	CheckCollaborator(context.Context, uuid.UUID, uuid.UUID) (bool, error)
	AddCollaborator(context.Context, uuid.UUID, uuid.UUID, string) error

	AddTag(ctx context.Context, tagName string, noteId uuid.UUID, userId uuid.UUID) (models.Note, error)
	DeleteTag(ctx context.Context, tagName string, noteId uuid.UUID, userId uuid.UUID) (models.Note, error)
	GetTags(ctx context.Context, userID uuid.UUID) ([]string, error)
}

type NoteBaseRepo interface {
	ReadAllNotes(context.Context, uuid.UUID, int64, int64, []string) ([]models.Note, error)
	ReadNote(context.Context, uuid.UUID) (models.Note, error)
	ReadAllNotesNoTags(context.Context, uuid.UUID, int64, int64) ([]models.Note, error)
	CreateNote(context.Context, models.Note) error
	UpdateNote(context.Context, models.Note) error
	DeleteNote(context.Context, uuid.UUID) error

	AddSubNote(context.Context, uuid.UUID, uuid.UUID) error
	RemoveSubNote(context.Context, uuid.UUID, uuid.UUID) error

	GetUpdates(context.Context, uuid.UUID, time.Time) ([]models.Message, error)

	CheckCollaborator(context.Context, uuid.UUID, uuid.UUID) (bool, error)
	AddCollaborator(context.Context, uuid.UUID, string) error

	AddTag(ctx context.Context, tagName string, noteId uuid.UUID) error
	DeleteTag(ctx context.Context, tagName string, noteId uuid.UUID) error
	GetTags(ctx context.Context, userID uuid.UUID) ([]string, error)
}

type NoteSearchRepo interface {
	SearchNotes(context.Context, uuid.UUID, int64, int64, string, []string) ([]models.Note, error)
	ReadNote(context.Context, uuid.UUID) (models.ElasticNote, error)
	CreateNote(context.Context, models.ElasticNote) error
	UpdateNote(context.Context, models.ElasticNote) error
	DeleteNote(context.Context, uuid.UUID) error

	AddSubNote(context.Context, uuid.UUID, uuid.UUID) error
	RemoveSubNote(context.Context, uuid.UUID, uuid.UUID) error

	AddCollaborator(context.Context, uuid.UUID, uuid.UUID) error

	AddTag(ctx context.Context, tagName string, noteID uuid.UUID) error
	DeleteTag(ctx context.Context, tagName string, noteID uuid.UUID) error
}
