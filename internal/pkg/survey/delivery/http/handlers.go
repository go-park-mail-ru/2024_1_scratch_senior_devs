package http

import (
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/survey/delivery/grpc/gen"
)

type SurveyHandler struct {
	client gen.StatClient
}

func CreateSurveyHandler(client gen.StatClient) *SurveyHandler {
	return &SurveyHandler{
		client: client,
	}
}
