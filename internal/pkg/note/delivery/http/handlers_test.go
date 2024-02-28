package http

import (
	"net/http"
	"testing"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"
)

func TestNoteHandler_GetAllNotes(t *testing.T) {
	type fields struct {
		uc note.NoteUsecase
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &NoteHandler{
				uc: tt.fields.uc,
			}
			h.GetAllNotes(tt.args.w, tt.args.r)
		})
	}
}
