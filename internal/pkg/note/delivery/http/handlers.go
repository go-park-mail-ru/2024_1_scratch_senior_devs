package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/usecase"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
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

	var data []models.Note
	var err error
	payload := r.Context().Value("payload")
	if castedPayload, ok := payload.(models.JwtPayload); ok {
		data, err = h.uc.GetAllNotes(r.Context(), castedPayload.Id, int64(count), int64(offset))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = utils.WriteResponseData(w, data, http.StatusAccepted)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error in GetAllNotes handler: %s", err)
		return
	}
}
