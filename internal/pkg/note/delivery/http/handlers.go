package http

import (
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/paging"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/request"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/response"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"
	"log/slog"
	"net/http"

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
// @Failure		400		{object}	response.ErrorResponse	true	"error"
// @Failure		401
// @Router		/api/note/get_all [get]
func (h *NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

	count, offset, err := paging.GetParams(r)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		response.WriteErrorMessage(w, http.StatusBadRequest, "invalid parameters")
		return
	}

	titleSubstr := r.URL.Query().Get("title")

	payload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, response.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := h.uc.GetAllNotes(r.Context(), payload.Id, int64(count), int64(offset), titleSubstr)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		response.WriteErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := response.WriteResponseData(w, data, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, response.WriteBodyError+err.Error())
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
// @Produce		json
// @Param		id		path		string					true	"note id"
// @Success		200		{object}	models.Note				true	"note"
// @Failure		400		{object}	response.ErrorResponse	true	"incorrect id"
// @Failure		401
// @Failure		404		{object}	response.ErrorResponse	true	"note not found"
// @Router		/api/note/{id} [get]
func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

	noteIdString := mux.Vars(r)["id"]
	noteId, err := uuid.FromString(noteIdString)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, "incorrect id parameter"+err.Error())
		response.WriteErrorMessage(w, http.StatusBadRequest, "note id must be a type of uuid")
		return
	}

	payload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, response.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	resultNote, err := h.uc.GetNote(r.Context(), noteId, payload.Id)
	if err != nil {
		log.LogHandlerError(logger, http.StatusNotFound, err.Error())
		response.WriteErrorMessage(w, http.StatusNotFound, err.Error())
		return
	}

	if err := response.WriteResponseData(w, resultNote, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, response.WriteBodyError+err.Error())
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
// @Param    	data  	body    	object          		true  	"note data"
// @Success		200		{object}	models.Note				true	"note"
// @Failure		400		{object}	response.ErrorResponse	true	"error"
// @Failure		401
// @Router		/api/note/add [post]
func (h *NoteHandler) AddNote(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

	jwtPayload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, response.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	noteData, err := request.ValidateRequestData(r)
	//noteData = []byte(html.EscapeString(string(noteData)))
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, response.ParseBodyError+err.Error())
		response.WriteErrorMessage(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	newNote, err := h.uc.CreateNote(r.Context(), jwtPayload.Id, noteData)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		response.WriteErrorMessage(w, http.StatusBadRequest, "invalid query")
		return
	}

	if err := response.WriteResponseData(w, newNote, http.StatusCreated); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, response.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}
