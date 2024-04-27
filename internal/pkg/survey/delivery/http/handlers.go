package http

import "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc/gen"

type SurveyHandler struct {
	client gen.NoteClient
}

func CreateSurveyHandler(client gen.NoteClient) *SurveyHandler {
	return &SurveyHandler{
		client: client,
	}
}
