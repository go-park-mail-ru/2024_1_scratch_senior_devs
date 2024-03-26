package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/delivery"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/paging"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"
)

type NoteHandler struct {
	uc     note.NoteUsecase
	logger *slog.Logger
}

func CreateNotesHandler(uc note.NoteUsecase, logger *slog.Logger) *NoteHandler {
	return &NoteHandler{
		uc:     uc,
		logger: logger,
	}
}

const (
	incorrectIdErr = "incorrect id paramter"
)

// GetAllNotes godoc
// @Summary		Get all notes
// @Description	Get a list of notes of current user
// @Tags 		note
// @ID			get-all-notes
// @Produce		json
// @Param		count	query		int							false	"notes count"
// @Param		offset	query		int							false	"notes offset"
// @Param		title	query		string						false	"notes title substring"
// @Success		200		{object}	[]models.NoteForSwagger		true	"notes"
// @Failure		400		{object}	response.ErrorResponse		true	"error"
// @Failure		401
// @Router		/api/note/get_all [get]
func (h *NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

	count, offset, err := paging.GetParams(r)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		delivery.WriteErrorMessage(w, http.StatusBadRequest, "invalid parameters")
		return
	}

	titleSubstr := r.URL.Query().Get("title")

	payload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, delivery.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := h.uc.GetAllNotes(r.Context(), payload.Id, int64(count), int64(offset), titleSubstr)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		delivery.WriteErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := delivery.WriteResponseData(w, data, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, delivery.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

// GetNote godoc
// @Summary		Get one note
// @Description	Get one of notes of current user
// @Tags 		note
// @ID			get-note
// @Param		id		path		string					true	"note id"
// @Success		200		{object}	models.NoteForSwagger	true	"note"
// @Failure		400		{object}	response.ErrorResponse	true	"incorrect id"
// @Failure		401
// @Failure		404		{object}	response.ErrorResponse	true	"note not found"
// @Router		/api/note/{id} [get]
func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

	noteIdString := mux.Vars(r)["id"]
	noteId, err := uuid.FromString(noteIdString)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, incorrectIdErr+err.Error())
		delivery.WriteErrorMessage(w, http.StatusBadRequest, "note id must be a type of uuid")
		return
	}

	payload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, delivery.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	resultNote, err := h.uc.GetNote(r.Context(), noteId, payload.Id)
	if err != nil {
		log.LogHandlerError(logger, http.StatusNotFound, err.Error())
		delivery.WriteErrorMessage(w, http.StatusNotFound, err.Error())
		return
	}

	if err := delivery.WriteResponseData(w, resultNote, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, delivery.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

// AddNote godoc
// @Summary		Add note
// @Description	Create new note to current user
// @Tags 		note
// @ID			add-note
// @Accept		json
// @Produce		json
// @Param    	payload		body    	models.UpsertNoteRequestForSwagger  	true  	"note data"
// @Success		200			{object}	models.NoteForSwagger					true	"note"
// @Failure		400			{object}	response.ErrorResponse					true	"error"
// @Failure		401
// @Router		/api/note/add [post]
func (h *NoteHandler) AddNote(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

	jwtPayload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, delivery.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	payload := models.UpsertNoteRequest{}
	if err := delivery.GetRequestData(r, &payload); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, delivery.ParseBodyError+err.Error())
		delivery.WriteErrorMessage(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	noteData, err := json.Marshal(payload.Data)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, delivery.ParseBodyError+err.Error())
		delivery.WriteErrorMessage(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	newNote, err := h.uc.CreateNote(r.Context(), jwtPayload.Id, noteData)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		delivery.WriteErrorMessage(w, http.StatusBadRequest, "invalid query")
		return
	}

	if err := delivery.WriteResponseData(w, newNote, http.StatusCreated); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, delivery.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

// UpdateNote godoc
// @Summary		Update note
// @Description	Create new note to current user
// @Tags 		note
// @ID			update-note
// @Accept		json
// @Produce		json
// @Param		id			path		string								true	"note id"
// @Param    	payload		body    	models.UpsertNoteRequestForSwagger	true  	"note data"
// @Success		200			{object}	models.NoteForSwagger				true	"note"
// @Failure		400			{object}	response.ErrorResponse				true	"error"
// @Failure		401
// @Router		/api/note/{id}/edit [post]
func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

	noteIdString := mux.Vars(r)["id"]
	noteId, err := uuid.FromString(noteIdString)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, incorrectIdErr+err.Error())
		delivery.WriteErrorMessage(w, http.StatusBadRequest, "note id must be a type of uuid")
		return
	}

	jwtPayload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, delivery.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	payload := models.UpsertNoteRequest{}
	if err := delivery.GetRequestData(r, &payload); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, delivery.ParseBodyError+err.Error())
		delivery.WriteErrorMessage(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	noteData, err := json.Marshal(payload.Data)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, delivery.ParseBodyError+err.Error())
		delivery.WriteErrorMessage(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	updatedNote, err := h.uc.UpdateNote(r.Context(), noteId, jwtPayload.Id, noteData)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		delivery.WriteErrorMessage(w, http.StatusBadRequest, "note not found")
		return
	}

	if err := delivery.WriteResponseData(w, updatedNote, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, delivery.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

// DeleteNote godoc
// @Summary		Delete note
// @Description	Delete selected note of current user
// @Tags 		note
// @ID			delete-note
// @Param		id		path		string					true	"note id"
// @Success		204
// @Failure		400		{object}	response.ErrorResponse	true	"incorrect id"
// @Failure		401
// @Failure		404		{object}	response.ErrorResponse	true	"note not found"
// @Router		/api/note/{id}/delete [delete]
func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

	noteIdString := mux.Vars(r)["id"]
	noteId, err := uuid.FromString(noteIdString)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, incorrectIdErr+err.Error())
		delivery.WriteErrorMessage(w, http.StatusBadRequest, "note id must be a type of uuid")
		return
	}

	jwtPayload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, delivery.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = h.uc.DeleteNote(r.Context(), noteId, jwtPayload.Id)
	if err != nil {
		log.LogHandlerError(logger, http.StatusNotFound, err.Error())
		delivery.WriteErrorMessage(w, http.StatusNotFound, "note not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.LogHandlerInfo(logger, http.StatusNoContent, "success")
}
