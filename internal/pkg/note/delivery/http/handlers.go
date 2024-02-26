package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/usecase"
)

type NoteHandler struct {
	uc usecase.NotesUsecase
}

func CreateNotesHandler(uc usecase.NotesUsecase) *NoteHandler {
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
		offset = 0
	}
	var data []models.Note
	var err error
	payload := r.Context().Value("payload")
	if castedPayload, ok := payload.(models.JwtPayload); ok {
		data, err = h.uc.GetAllNotes(r.Context(), castedPayload.Id, int64(count), int64(offset))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	marshaledData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Write(marshaledData)

}
