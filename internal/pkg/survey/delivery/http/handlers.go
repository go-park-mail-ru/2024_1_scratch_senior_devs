package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/survey/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/responses"
	"github.com/satori/uuid"
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

func getQuestion(question *gen.Question) models.Question {
	return models.Question{
		Id:           uuid.FromStringOrNil(question.Id),
		Title:        question.Title,
		QuestionType: question.QuestionType,
		Number:       int(question.Number),
		SurveyId:     uuid.FromStringOrNil(question.SurveyId),
	}
}

func (h *SurveyHandler) CreateSurvey(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	_, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	payload := models.CreateSurveyRequest{}
	if err := responses.GetRequestData(r, &payload); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, responses.ParseBodyError+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("incorrect data format"))
		return
	}

	questions := make([]*gen.CreateQuestionRequest, len(payload.Questions))
	for i, question := range payload.Questions {
		questions[i] = &gen.CreateQuestionRequest{
			Title:        question.Title,
			QuestionType: question.QuestionType,
		}
	}
	if _, err := h.client.CreateSurvey(r.Context(), &gen.CreateSurveyRequest{
		Questions: questions,
	}); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, responses.ParseBodyError+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("error"))
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
	w.WriteHeader(http.StatusNoContent)
}

func getStat(stat *gen.StatModel) models.Stat {
	val, _ := strconv.ParseInt(stat.Label, 10, 32)
	return models.Stat{
		Title:        stat.QuestionTitle,
		QuestionType: stat.QuestionType,
		Voice:        int(val),
		Count:        int(stat.Value),
	}
}
func (h *SurveyHandler) GetSurvey(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	protoData, err := h.client.GetSurvey(r.Context(), &gen.GetSurveyRequest{})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, err)
		return
	}
	data := make([]models.Question, len(protoData.Questions))

	for i, question := range protoData.Questions {
		data[i] = getQuestion(question)
	}
	if err := responses.WriteResponseData(w, data, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

func (h *SurveyHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	protoData, err := h.client.GetStats(r.Context(), &gen.GetStatsRequest{})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, err)
		return
	}
	data := make([]models.Stat, len(protoData.Stats))

	for i, s := range protoData.Stats {
		data[i] = getStat(s)
	}
	if err := responses.WriteResponseData(w, data, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}
