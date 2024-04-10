package models

import "github.com/satori/uuid"

type Attach struct {
	Id     uuid.UUID `json:"id"`
	Path   string    `json:"path"`
	NoteId uuid.UUID `json:"note_id"`
}
