package note

import (
	"context"
	"time"

	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

const (
	ErrTooDeep              = "too deep"
	ErrTooManyTags          = "too many tags"
	ErrTooManySubnotes      = "too many subnotes"
	ErrTooManyCollaborators = "too many collaborators"
	ErrAlreadyCollaborator  = "already a collaborator"
)

type NoteUsecase interface {
	GetAllNotes(context.Context, uuid.UUID, int64, int64, string, []string) ([]models.Note, error)
	GetNote(context.Context, uuid.UUID, uuid.UUID) (models.Note, error)
	GetPublicNote(ctx context.Context, noteId uuid.UUID) (models.Note, error)
	CreateNote(context.Context, uuid.UUID, string) (models.Note, error)
	UpdateNote(context.Context, uuid.UUID, uuid.UUID, string) (models.Note, error)
	DeleteNote(context.Context, uuid.UUID, uuid.UUID) error

	CreateSubNote(context.Context, uuid.UUID, string, uuid.UUID) (models.Note, error)

	AddCollaborator(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) (string, error)

	AddTag(ctx context.Context, tagName string, noteId uuid.UUID, userId uuid.UUID) (models.Note, error)
	DeleteTag(ctx context.Context, tagName string, noteId uuid.UUID, userId uuid.UUID) (models.Note, error)
	GetTags(ctx context.Context, userID uuid.UUID) ([]string, error)

	UpdateTag(ctx context.Context, oldTag string, newTag string, userID uuid.UUID) error
	RememberTag(ctx context.Context, tagName string, userID uuid.UUID) error
	ForgetTag(ctx context.Context, tagName string, userID uuid.UUID) error

	SetIcon(ctx context.Context, noteID uuid.UUID, icon string, userID uuid.UUID) (models.Note, error)
	SetHeader(ctx context.Context, noteID uuid.UUID, header string, userID uuid.UUID) (models.Note, error)

	AddFav(ctx context.Context, noteID uuid.UUID, userID uuid.UUID) (models.Note, error)
	DelFav(ctx context.Context, noteID uuid.UUID, userID uuid.UUID) (models.Note, error)

	SetPublic(ctx context.Context, noteID uuid.UUID, userID uuid.UUID) (models.Note, error)
	SetPrivate(ctx context.Context, noteID uuid.UUID, userID uuid.UUID) (models.Note, error)

	GetAttachList(ctx context.Context, noteID uuid.UUID, userID uuid.UUID) ([]string, error)
	GetSharedAttachList(ctx context.Context, noteID uuid.UUID) ([]string, error)

	CheckPermissions(ctx context.Context, noteID uuid.UUID, userID uuid.UUID) (bool, error)
}

type NoteBaseRepo interface {
	ReadAllNotes(context.Context, uuid.UUID, int64, int64, []string) ([]models.Note, error)
	ReadNote(context.Context, uuid.UUID, uuid.UUID) (models.Note, error)
	ReadPublicNote(context.Context, uuid.UUID) (models.Note, error)
	CreateNote(context.Context, models.Note) error
	UpdateNote(context.Context, models.Note) error
	DeleteNote(context.Context, uuid.UUID) error

	AddSubNote(context.Context, uuid.UUID, uuid.UUID) error
	RemoveSubNote(context.Context, uuid.UUID, uuid.UUID) error

	GetUpdates(context.Context, uuid.UUID, time.Time) ([]models.Message, error)

	AddCollaborator(context.Context, uuid.UUID, uuid.UUID) (string, error)

	AddTag(ctx context.Context, tagName string, noteId uuid.UUID) error
	DeleteTag(ctx context.Context, tagName string, noteId uuid.UUID) error

	GetTags(ctx context.Context, userID uuid.UUID) ([]string, error)
	RememberTag(ctx context.Context, tagName string, userID uuid.UUID) error
	ForgetTag(ctx context.Context, tagName string, userID uuid.UUID) error
	DeleteTagFromAllNotes(ctx context.Context, tagName string, userID uuid.UUID) error
	UpdateTag(ctx context.Context, oldTag string, newTag string, userID uuid.UUID) error
	UpdateTagOnAllNotes(ctx context.Context, oldTag string, newTag string, userID uuid.UUID) error

	SetIcon(ctx context.Context, noteID uuid.UUID, icon string) error
	SetHeader(ctx context.Context, noteID uuid.UUID, header string) error

	AddFav(ctx context.Context, noteID uuid.UUID, userID uuid.UUID) error
	DelFav(ctx context.Context, noteID uuid.UUID, userID uuid.UUID) error

	SetPublic(ctx context.Context, noteID uuid.UUID) error
	SetPrivate(ctx context.Context, noteID uuid.UUID) error

	GetAttachList(ctx context.Context, noteID uuid.UUID) ([]string, error)
}

type NoteSearchRepo interface {
	SearchNotes(context.Context, uuid.UUID, int64, int64, string, []string) ([]models.Note, error)
	CreateNote(context.Context, models.Note) error
	UpdateNote(context.Context, models.Note) error
	DeleteNote(context.Context, uuid.UUID) error

	AddSubNote(context.Context, uuid.UUID, uuid.UUID) error
	RemoveSubNote(context.Context, uuid.UUID, uuid.UUID) error

	AddCollaborator(context.Context, uuid.UUID, uuid.UUID) error

	AddTag(ctx context.Context, tagName string, noteID uuid.UUID) error
	DeleteTag(ctx context.Context, tagName string, noteID uuid.UUID) error
	DeleteTagFromAllNotes(ctx context.Context, tagName string, userID uuid.UUID) error
	UpdateTagOnAllNotes(ctx context.Context, oldTag string, newTag string, userID uuid.UUID) error

	SetIcon(ctx context.Context, noteID uuid.UUID, icon string) error
	SetHeader(ctx context.Context, noteID uuid.UUID, header string) error

	ChangeFlag(ctx context.Context, noteID uuid.UUID, flag bool) error

	SetPublic(ctx context.Context, noteID uuid.UUID) error
	SetPrivate(ctx context.Context, noteID uuid.UUID) error
}
