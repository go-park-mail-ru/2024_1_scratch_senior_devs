package attach

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/satori/uuid"
	"io"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type AttachUsecase interface {
	AddAttach(ctx context.Context, noteID uuid.UUID, userID uuid.UUID, attach io.ReadSeeker, extension string) (models.Attach, error)
	DeleteAttach(ctx context.Context, attachID uuid.UUID, userID uuid.UUID) error
}

type AttachRepo interface {
	GetAttach(ctx context.Context, id uuid.UUID) (models.Attach, error)
	AddAttach(ctx context.Context, attach models.Attach) error
	DeleteAttach(ctx context.Context, id uuid.UUID) error
}
