package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/survey/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/responses"
)

type SurveyHandler struct {
	client gen.StatClient
}

func CreateSurveyHandler(client gen.StatClient) *SurveyHandler {
	return &SurveyHandler{
		client: client,
	}
}

func (h *SurveyHandler) Vote(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	payload := models.Vote{}
	if err := responses.GetRequestData(r, &payload); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, responses.ParseBodyError+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("incorrect data format"))
		return
	}

	_, err := h.client.Vote(r.Context(), &gen.VoteRequest{
		QuestionId: payload.QuestionId.String(),
		Vote:       int32(payload.Voice),
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("invalid query"))
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}
