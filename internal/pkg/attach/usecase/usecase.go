package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/attach"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/sources"
	"github.com/satori/uuid"
	"io"
	"log/slog"
	"os"
	"path"
)

type AttachUsecase struct {
	repo   attach.AttachRepo
	logger *slog.Logger
}

func CreateAttachUsecase(repo attach.AttachRepo, logger *slog.Logger) *AttachUsecase {
	return &AttachUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *AttachUsecase) AddAttach(ctx context.Context, noteID uuid.UUID, attach io.ReadSeeker, extension string) (models.Attach, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	attachBasePath := os.Getenv("ATTACHES_BASE_PATH")
	newAttachId := uuid.NewV4()
	newAttachPath := newAttachId.String() + extension

	if err := sources.WriteFileOnDisk(path.Join(attachBasePath, newAttachPath), attach); err != nil {
		logger.Error("write on disk: " + err.Error())
		return models.Attach{}, err
	}

	newAttach := models.Attach{
		Id:     newAttachId,
		Path:   newAttachPath,
		NoteId: noteID,
	}

	if err := uc.repo.AddAttach(ctx, newAttach); err != nil {
		logger.Error(err.Error())
		return models.Attach{}, err
	}

	logger.Info("success")
	return newAttach, nil
}
