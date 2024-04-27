package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/survey"
	generatedSurvey "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/survey/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/satori/uuid"
)

type GrpcSurveyHandler struct {
	generatedSurvey.StatServer
	uc survey.SurveyUsecase
}

func NewGrpcSurveyHandler(uc survey.SurveyUsecase) *GrpcSurveyHandler {
	return &GrpcSurveyHandler{uc: uc}
}

func (h *GrpcSurveyHandler) AddNote(ctx context.Context, in *generatedSurvey.VoteRequest) (*generatedSurvey.VoteResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))
	vote := in.Vote
	qid := uuid.FromStringOrNil(in.QuestionId)
	err := h.uc.Vote(ctx, models.Vote{
		QuestionId: qid,
		Voice:      int(vote),
	})
	if err != nil {
		logger.Error(err.Error())
		return &generatedSurvey.VoteResponse{}, errors.New("not found")
	}

	logger.Info("success")
	return &generatedSurvey.VoteResponse{}, nil
}
