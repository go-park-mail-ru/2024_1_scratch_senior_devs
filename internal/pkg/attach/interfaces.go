package attach

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/satori/uuid"
	"io"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type AttachUsecase interface {
	AddAttach(ctx context.Context, noteID uuid.UUID, attach io.ReadSeeker, extension string) (models.Attach, error)
}

type AttachRepo interface {
	AddAttach(ctx context.Context, attach models.Attach) error
}
