package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"
	generatedNote "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/satori/uuid"
)

type GrpcNoteHandler struct {
	generatedNote.NoteServer
	uc note.NoteUsecase
}

func NewGrpcNoteHandler(uc note.NoteUsecase) *GrpcNoteHandler {
	return &GrpcNoteHandler{uc: uc}
}

func getNote(note models.Note) *generatedNote.NoteModel {
	children := make([]string, len(note.Children))
	for i, child := range note.Children {
		children[i] = child.String()
	}

	collaborators := make([]string, len(note.Collaborators))
	for i, collaborator := range note.Collaborators {
		collaborators[i] = collaborator.String()
	}

	tags := make([]string, len(note.Tags))
	copy(tags, note.Tags)

	return &generatedNote.NoteModel{
		Id:            note.Id.String(),
		Data:          note.Data,
		CreateTime:    note.CreateTime.String(),
		UpdateTime:    note.UpdateTime.String(),
		OwnerId:       note.OwnerId.String(),
		Parent:        note.Parent.String(),
		Children:      children,
		Tags:          tags,
		Collaborators: collaborators,
	}
}

func (h *GrpcNoteHandler) GetTags(ctx context.Context, in *generatedNote.GetTagsRequest) (*generatedNote.GetTagsResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	tags, err := h.uc.GetTags(ctx, uuid.FromStringOrNil(in.UserId))
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("success")
	return &generatedNote.GetTagsResponse{Tags: tags}, nil
}

func (h *GrpcNoteHandler) AddTag(ctx context.Context, in *generatedNote.TagRequest) (*generatedNote.GetNoteResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	result, err := h.uc.AddTag(ctx, in.TagName, uuid.FromStringOrNil(in.NoteId), uuid.FromStringOrNil(in.UserId))
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("success")
	return &generatedNote.GetNoteResponse{Note: getNote(result)}, nil
}

func (h *GrpcNoteHandler) DeleteTag(ctx context.Context, in *generatedNote.TagRequest) (*generatedNote.GetNoteResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	result, err := h.uc.DeleteTag(ctx, in.TagName, uuid.FromStringOrNil(in.NoteId), uuid.FromStringOrNil(in.UserId))
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("success")
	return &generatedNote.GetNoteResponse{Note: getNote(result)}, nil
}

func (h *GrpcNoteHandler) RememberTag(ctx context.Context, in *generatedNote.AllTagRequest) (*generatedNote.EmptyResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	if err := h.uc.RememberTag(ctx, in.TagName, uuid.FromStringOrNil(in.UserId)); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("success")
	return &generatedNote.EmptyResponse{}, nil
}

func (h *GrpcNoteHandler) ForgetTag(ctx context.Context, in *generatedNote.AllTagRequest) (*generatedNote.EmptyResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	if err := h.uc.ForgetTag(ctx, in.TagName, uuid.FromStringOrNil(in.UserId)); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("success")
	return &generatedNote.EmptyResponse{}, nil
}

func (h *GrpcNoteHandler) AddCollaborator(ctx context.Context, in *generatedNote.AddCollaboratorRequest) (*generatedNote.AddCollaboratorResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	if err := h.uc.AddCollaborator(ctx, uuid.FromStringOrNil(in.NoteId), uuid.FromStringOrNil(in.UserId), uuid.FromStringOrNil(in.GuestId)); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("success")
	return &generatedNote.AddCollaboratorResponse{}, nil
}

func (h *GrpcNoteHandler) GetAllNotes(ctx context.Context, in *generatedNote.GetAllRequest) (*generatedNote.GetAllResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	result, err := h.uc.GetAllNotes(ctx, uuid.FromStringOrNil(in.UserId), in.Count, in.Offset, in.Title, in.Tags)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("not found")
	}

	protoNotes := make([]*generatedNote.NoteModel, len(result))
	for i, item := range result {
		protoNotes[i] = getNote(item)
	}

	logger.Info("success")
	return &generatedNote.GetAllResponse{
		Notes: protoNotes,
	}, nil
}

func (h *GrpcNoteHandler) GetNote(ctx context.Context, in *generatedNote.GetNoteRequest) (*generatedNote.GetNoteResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	result, err := h.uc.GetNote(ctx, uuid.FromStringOrNil(in.Id), uuid.FromStringOrNil(in.UserId))
	if err != nil {
		logger.Error(err.Error())

		return nil, errors.New("not found")
	}

	logger.Info("success")
	return &generatedNote.GetNoteResponse{
		Note: getNote(result),
	}, nil
}

func (h *GrpcNoteHandler) AddNote(ctx context.Context, in *generatedNote.AddNoteRequest) (*generatedNote.AddNoteResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	result, err := h.uc.CreateNote(ctx, uuid.FromStringOrNil(in.UserId), in.Data)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("not found")
	}

	logger.Info("success")
	return &generatedNote.AddNoteResponse{
		Note: getNote(result),
	}, nil
}

func (h *GrpcNoteHandler) UpdateNote(ctx context.Context, in *generatedNote.UpdateNoteRequest) (*generatedNote.UpdateNoteResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	result, err := h.uc.UpdateNote(ctx, uuid.FromStringOrNil(in.Id), uuid.FromStringOrNil(in.UserId), in.Data)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("not found")
	}

	logger.Info("success")
	return &generatedNote.UpdateNoteResponse{
		Note: getNote(result),
	}, nil
}

func (h *GrpcNoteHandler) DeleteNote(ctx context.Context, in *generatedNote.DeleteNoteRequest) (*generatedNote.DeleteNoteResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	err := h.uc.DeleteNote(ctx, uuid.FromStringOrNil(in.Id), uuid.FromStringOrNil(in.UserId))
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("not found")
	}

	logger.Info("success")
	return &generatedNote.DeleteNoteResponse{}, nil
}

func (h *GrpcNoteHandler) CreateSubNote(ctx context.Context, in *generatedNote.CreateSubNoteRequest) (*generatedNote.CreateSubNoteResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	response, err := h.uc.CreateSubNote(ctx, uuid.FromStringOrNil(in.UserId), in.NoteData, uuid.FromStringOrNil(in.ParentId))
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("success")
	return &generatedNote.CreateSubNoteResponse{
		Note: getNote(response),
	}, nil
}

func (h *GrpcNoteHandler) CheckPermissions(ctx context.Context, in *generatedNote.CheckPermissionsRequest) (*generatedNote.CheckPermissionsResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	response, err := h.uc.CheckPermissions(ctx, uuid.FromStringOrNil(in.NoteId), uuid.FromStringOrNil(in.UserId))
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("success")
	return &generatedNote.CheckPermissionsResponse{
		Result: response,
	}, nil
}
