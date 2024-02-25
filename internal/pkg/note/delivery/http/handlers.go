package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"
	"github.com/gofrs/uuid"
)

type NoteHandler struct {
	uc note.NoteUsecase
}

func NewNoteHandler(uc note.NoteUsecase) *NoteHandler {
	return &NoteHandler{
		uc: uc,
	}
}

func (h *NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.URL.Query().Get("count"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if count == 0 {
		count = 10
	}
	if offset == 0 {
		offset = 10
	}
	payload := r.Context().Value("payload").(models.JwtPayload)
	data, err := h.uc.GetAllNotes(r.Context(), uuid.UUID(payload.Id), int64(count), int64(offset))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	marshaledData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Write(marshaledData)

}
