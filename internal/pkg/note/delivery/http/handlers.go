package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
)

const (
	defaultCount  = 10
	defaultOffset = 0
)

type NoteHandler struct {
	uc note.NoteUsecase
}

func CreateNotesHandler(uc note.NoteUsecase) *NoteHandler {
	return &NoteHandler{
		uc: uc,
	}
}

// GetAllNotes
// @Summary		Get all notes
// @Description	Get a list of notes of current user
// @Tags 		note
// @ID			get-all-notes
// @Produce		json
// @Success		200		{object}	[]models.Note			true	"notes"
// @Failure		400		{object}	utils.ErrorResponse		true	"error"
// @Failure		401
// @Router		/api/note/get_all [get]
func (h *NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	count := defaultCount
	offset := defaultOffset

	var err error

	strCount := r.URL.Query().Get("count")
	if strCount != "" {
		count, err = strconv.Atoi(strCount)
		if err != nil {
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
			utils.WriteErrorMessage(w, http.StatusBadRequest, "incorrect offset parameter")
			return
		}
		if offset < 0 {
			offset = defaultOffset
		}
	}

	payload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := h.uc.GetAllNotes(r.Context(), payload.Id, int64(count), int64(offset))
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	err = utils.WriteResponseData(w, data, http.StatusOK)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error in GetAllNotes handler: %s", err)
		return
	}
}
