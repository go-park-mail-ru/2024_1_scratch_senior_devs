package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/survey/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/responses"
	"github.com/samber/lo"
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
		MinMark:      int(question.MinMark),
		Skip:         int(question.Skip),
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
			MinMark:      int64(question.MinMark),
			Skip:         int64(question.Skip),
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
		QuestionId:   uuid.FromStringOrNil(stat.QuestionId),
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

func getCSAT(data []models.Stat) (map[int]int, float64) {
	result := make(map[int]int)
	var summary int
	var total int

	data = lo.Filter(data, func(item models.Stat, i int) bool {
		return item.Voice > 0 && item.Voice <= 5
	})

	for _, item := range data {
		result[item.Voice] = item.Count
		summary += item.Voice * item.Count
		total += item.Count
	}
	return result, float64(summary) / float64(total)
}
func getNPS(data []models.Stat) (map[string]int, float64) {
	result := make(map[string]int)
	// for _, item := range data {
	// 	if item.Voice < 6{
	// 		result["d"] = item.Count
	// 	}

	// }

	lo.ForEach(data, func(d models.Stat, i int) {
		if d.Voice <= 6 {
			data[i].Type = "detractor"
		} else if d.Voice >= 6 && d.Voice <= 8 {
			data[i].Type = "passive"
		} else {
			data[i].Type = "promouter"
		}
	})
	var total int
	for _, v := range data {
		result[v.Type] += v.Count
		total += v.Count
	}

	d, _ := result["detractor"]
	p, _ := result["promouter"]

	var r float64
	if total != 0 {
		r = float64(p-d) / float64(total)
	}
	return result, r
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

	var groupedData map[string][]models.Stat

	groupedData = lo.GroupBy(data, func(item models.Stat) string {
		return string(item.QuestionId.String())

	})

	var resp []models.StatResponse

	for _, item := range groupedData {
		if len(item) == 0 {
			continue
		}

		respItem := models.StatResponse{
			QuestionId:   item[0].QuestionId,
			Title:        item[0].Title,
			QuestionType: item[0].QuestionType,
		}
		switch item[0].QuestionType {
		case "CSAT":
			respItem.Stats, respItem.Value = getCSAT(item)
		case "NPS":
			respItem.Stats, respItem.Value = getNPS(item)

		}
		resp = append(resp, respItem)
	}

	if err := responses.WriteResponseData(w, resp, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}
