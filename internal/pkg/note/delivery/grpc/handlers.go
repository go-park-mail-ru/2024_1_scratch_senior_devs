package grpc

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"
	generatedNote "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc/gen"
)

type GrpcNoteHandler struct {
	generatedNote.NoteServer
	uc note.NoteUsecase
}

func (h *GrpcNoteHandler) GetAll(ctx context.Context, in *generatedNote.GetAllRequest) (*generatedNote.GetAllResponse, error) {

	payload, ok := ctx.Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		return nil, errors.New("cant parse payload")
	}

	result, err := h.uc.GetAllNotes(ctx, payload.Id, in.Count, in.Offset, in.Title)
	if err != nil {
		return nil, errors.New("not found")
	}
	protoNotes := make([]*generatedNote.NoteModel, len(result))
	for i, item := range result {
		protoNotes[i] = &generatedNote.NoteModel{
			Id:         item.Id.String(),
			Data:       string(item.Data),
			CreateTime: item.CreateTime.String(),
			UpdateTime: item.UpdateTime.String(),
			OwnerId:    item.OwnerId.String(),
		}
	}
	return &generatedNote.GetAllResponse{
		Notes: protoNotes,
	}, nil

}
