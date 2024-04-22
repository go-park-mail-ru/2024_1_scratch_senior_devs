package grpc

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"
	generatedNote "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc/gen"
	"github.com/satori/uuid"
)

type GrpcNoteHandler struct {
	generatedNote.NoteServer
	uc note.NoteUsecase
}

func NewGrpcNoteHandler(uc note.NoteUsecase) *GrpcNoteHandler {
	return &GrpcNoteHandler{uc: uc}
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

func (h *GrpcNoteHandler) GetNote(ctx context.Context, in *generatedNote.GetNoteRequest) (*generatedNote.GetNoteResponse, error) {
	payload, ok := ctx.Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		return nil, errors.New("cant parse payload")
	}

	result, err := h.uc.GetNote(ctx, uuid.FromStringOrNil(in.Id), payload.Id)
	if err != nil {
		return nil, errors.New("not found")
	}

	protoNote := generatedNote.NoteModel{
		Id:         result.Id.String(),
		Data:       string(result.Data),
		CreateTime: result.CreateTime.String(),
		UpdateTime: result.UpdateTime.String(),
		OwnerId:    result.OwnerId.String(),
	}
	return &generatedNote.GetNoteResponse{
		Note: &protoNote,
	}, nil
}

func (h *GrpcNoteHandler) AddNote(ctx context.Context, in *generatedNote.AddNoteRequest) (*generatedNote.AddNoteResponse, error) {
	payload, ok := ctx.Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		return nil, errors.New("cant parse payload")
	}

	result, err := h.uc.CreateNote(ctx, payload.Id, []byte(in.Data))
	if err != nil {
		return nil, errors.New("not found")
	}

	protoNote := generatedNote.NoteModel{
		Id:         result.Id.String(),
		Data:       string(result.Data),
		CreateTime: result.CreateTime.String(),
		UpdateTime: result.UpdateTime.String(),
		OwnerId:    result.OwnerId.String(),
	}
	return &generatedNote.AddNoteResponse{
		Note: &protoNote,
	}, nil
}

func (h *GrpcNoteHandler) UpdateNote(ctx context.Context, in *generatedNote.UpdateNoteRequest) (*generatedNote.UpdateNoteResponse, error) {
	payload, ok := ctx.Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		return nil, errors.New("cant parse payload")
	}

	result, err := h.uc.UpdateNote(ctx, uuid.FromStringOrNil(in.Id), payload.Id, []byte(in.Data))
	if err != nil {
		return nil, errors.New("not found")
	}

	protoNote := generatedNote.NoteModel{
		Id:         result.Id.String(),
		Data:       string(result.Data),
		CreateTime: result.CreateTime.String(),
		UpdateTime: result.UpdateTime.String(),
		OwnerId:    result.OwnerId.String(),
	}
	return &generatedNote.UpdateNoteResponse{
		Note: &protoNote,
	}, nil
}
func (h *GrpcNoteHandler) DeleteNote(ctx context.Context, in *generatedNote.DeleteNoteRequest) (*generatedNote.DeleteNoteResponse, error) {
	payload, ok := ctx.Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		return nil, errors.New("cant parse payload")
	}

	err := h.uc.DeleteNote(ctx, uuid.FromStringOrNil(in.Id), payload.Id)
	if err != nil {
		return nil, errors.New("not found")
	}

	return &generatedNote.DeleteNoteResponse{}, nil
}