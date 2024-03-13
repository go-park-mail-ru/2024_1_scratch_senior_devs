package http

import (
	"github.com/gorilla/mux"
	"github.com/satori/uuid"
	"log/slog"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
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

// GetAllNotes godoc
// @Summary		Get all notes
// @Description	Get a list of notes of current user
// @Tags 		note
// @ID			get-all-notes
// @Produce		json
// @Param		count	query		int						false	"notes count"
// @Param		offset	query		int						false	"notes offset"
// @Param		title	query		string					false	"notes title substring"
// @Success		200		{object}	[]models.Note			true	"notes"
// @Failure		400		{object}	utils.ErrorResponse		true	"error"
// @Failure		401
// @Router		/api/note/get_all [get]
func (h *NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", utils.GetRequestId(r.Context())), slog.String("func", utils.GFN()))

	count, offset, err := utils.GetParams(r)
	if err != nil {
		utils.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		utils.WriteErrorMessage(w, http.StatusBadRequest, "invalid parameters")
		return
	}

	titleSubstr := r.URL.Query().Get("title")

	payload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		utils.LogHandlerError(logger, http.StatusUnauthorized, utils.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := h.uc.GetAllNotes(r.Context(), payload.Id, int64(count), int64(offset), titleSubstr)
	if err != nil {
		utils.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		utils.WriteErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := utils.WriteResponseData(w, data, http.StatusOK); err != nil {
		utils.LogHandlerError(logger, http.StatusInternalServerError, utils.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.LogHandlerInfo(logger, http.StatusOK, "success")
}

// GetNote godoc
// @Summary		Get one note
// @Description	Get one of notes of current user
// @Tags 		note
// @ID			get-note
// @Produce		json
// @Param		id		path		string					true	"note id"
// @Success		200		{object}	models.Note				true	"note"
// @Failure		400		{object}	utils.ErrorResponse		true	"incorrect id"
// @Failure		401
// @Failure		404		{object}	utils.ErrorResponse		true	"note not found"
// @Router		/api/note/{id} [get]
func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", utils.GetRequestId(r.Context())), slog.String("func", utils.GFN()))

	noteIdString := mux.Vars(r)["id"]
	noteId, err := uuid.FromString(noteIdString)
	if err != nil {
		utils.LogHandlerError(logger, http.StatusBadRequest, "incorrect id parameter"+err.Error())
		utils.WriteErrorMessage(w, http.StatusBadRequest, "note id must be a type of uuid")
		return
	}

	payload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		utils.LogHandlerError(logger, http.StatusUnauthorized, utils.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	resultNote, err := h.uc.GetNote(r.Context(), noteId, payload.Id)
	if err != nil {
		utils.LogHandlerError(logger, http.StatusNotFound, err.Error())
		utils.WriteErrorMessage(w, http.StatusNotFound, err.Error())
		return
	}

	if err := utils.WriteResponseData(w, resultNote, http.StatusOK); err != nil {
		utils.LogHandlerError(logger, http.StatusInternalServerError, utils.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.LogHandlerInfo(logger, http.StatusOK, "success")
}

// AddNote godoc
// @Summary		Add note
// @Description	Create new note to current user
// @Tags 		note
// @ID			add-note
// @Produce		json
// @Success		200		{object}	models.Note				true	"note"
// @Failure		400		{object}	utils.ErrorResponse		true	"error"
// @Failure		401
// @Router		/api/note/add [post]
func (h *NoteHandler) AddNote(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", utils.GetRequestId(r.Context())), slog.String("func", utils.GFN()))

	jwtPayload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		utils.LogHandlerError(logger, http.StatusUnauthorized, utils.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newNote, err := h.uc.CreateNote(r.Context(), jwtPayload.Id)
	if err != nil {
		utils.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		utils.WriteErrorMessage(w, http.StatusBadRequest, "invalid query")
		return
	}

	if err := utils.WriteResponseData(w, newNote, http.StatusCreated); err != nil {
		utils.LogHandlerError(logger, http.StatusInternalServerError, utils.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.LogHandlerInfo(logger, http.StatusOK, "success")
}
