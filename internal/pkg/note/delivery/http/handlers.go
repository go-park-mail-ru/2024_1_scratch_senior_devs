package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
)

const (
	defaultCount  = 10
	defaultOffset = 0
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
	count := defaultCount
	offset := defaultOffset

	h.logger.Info(utils.GFN())
	var err error

	strCount := r.URL.Query().Get("count")
	if strCount != "" {
		count, err = strconv.Atoi(strCount)
		if err != nil {
			h.logger.Error(err.Error())
			utils.WriteErrorMessage(w, http.StatusBadRequest, "incorrect count parameter")
			return
		}
		if count <= 0 {
			count = defaultCount
		}
	}

	strOffset := r.URL.Query().Get("offset")
	if strOffset != "" {
		offset, err = strconv.Atoi(strOffset)
		if err != nil {
			h.logger.Error(err.Error())
			utils.WriteErrorMessage(w, http.StatusBadRequest, "incorrect offset parameter")
			return
		}
		if offset < 0 {
			offset = defaultOffset
		}
	}

	titleSubstr := r.URL.Query().Get("title")

	payload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		h.logger.Info("Problem while getting jwt payload from context")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := h.uc.GetAllNotes(r.Context(), payload.Id, int64(count), int64(offset), titleSubstr)
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	err = utils.WriteResponseData(w, data, http.StatusOK)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error in GetAllNotes handler: %s", err)
		return
	}
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
	h.logger.Info(utils.GFN())
	noteIdString := mux.Vars(r)["id"]
	noteId, err := uuid.FromString(noteIdString)
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteErrorMessage(w, http.StatusBadRequest, "note id must be a type of uuid")
		return
	}

	payload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		h.logger.Info("Problem while getting jwt payload from context")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	resultNote, err := h.uc.GetNote(r.Context(), noteId, payload.Id)
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteErrorMessage(w, http.StatusNotFound, err.Error())
		return
	}

	err = utils.WriteResponseData(w, resultNote, http.StatusOK)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
	h.logger.Info(utils.GFN())
	jwtPayload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		h.logger.Info("Problem while getting jwt payload from context")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newNote, err := h.uc.CreateNote(r.Context(), jwtPayload.Id)
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteErrorMessage(w, http.StatusBadRequest, "invalid query")
		return
	}

	err = utils.WriteResponseData(w, newNote, http.StatusCreated)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error in AddNote handler: %s", err)
		return
	}
}
