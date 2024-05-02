package tags

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/satori/uuid"
)

type TagRepo interface {
	AddTag(ctx context.Context, tagName string, noteId uuid.UUID) error
	DeleteTag(ctx context.Context, tagName string, noteId uuid.UUID) error
}

type TagUsecase interface {
	AddTag(ctx context.Context, tagName string, noteId uuid.UUID, userId uuid.UUID) (models.Note, error)
	DeleteTag(ctx context.Context, tagName string, noteId uuid.UUID, userId uuid.UUID) (models.Note, error)
}
