package grpc

import (
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/survey"
	generatedSurvey "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/survey/delivery/grpc/gen"
)

type GrpcSurveyHandler struct {
	generatedSurvey.StatServer
	uc survey.SurveyUsecase
}

func NewGrpcSurveyHandler(uc survey.SurveyUsecase) *GrpcSurveyHandler {
	return &GrpcSurveyHandler{uc: uc}
}
