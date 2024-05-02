package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/tags"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/responses"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"
)

type TagHandler struct {
	uc tags.TagUsecase
}

func CreateTagHandler(uc tags.TagUsecase) *TagHandler {
	return &TagHandler{
		uc: uc,
	}
}

func (h *TagHandler) AddTag(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	jwtPayload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	noteIdString := mux.Vars(r)["id"]
	noteId, err := uuid.FromString(noteIdString)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("note id must be a type of uuid"))
		return
	}

	var tagName models.TagRequest
	if err := responses.GetRequestData(r, &tagName); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, responses.ParseBodyError+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("incorrect data format"))
		return
	}
	tagName.Sanitize()
	err = tagName.Validate()
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	note, err := h.uc.AddTag(r.Context(), tagName.TagName, noteId, jwtPayload.Id)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := responses.WriteResponseData(w, note, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

func (h *TagHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	jwtPayload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	noteIdString := mux.Vars(r)["id"]
	noteId, err := uuid.FromString(noteIdString)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("note id must be a type of uuid"))
		return
	}

	var tagName models.TagRequest
	if err := responses.GetRequestData(r, &tagName); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, responses.ParseBodyError+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("incorrect data format"))
		return
	}

	note, err := h.uc.DeleteTag(r.Context(), tagName.TagName, noteId, jwtPayload.Id)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := responses.WriteResponseData(w, note, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}
