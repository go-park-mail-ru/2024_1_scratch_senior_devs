package models

import (
	"html"
	"time"

	"github.com/satori/uuid"
)

type Note struct {
	Id            uuid.UUID   `json:"id"`
	Data          string      `json:"data,omitempty"`
	CreateTime    time.Time   `json:"create_time"`
	UpdateTime    time.Time   `json:"update_time"`
	OwnerId       uuid.UUID   `json:"owner_id"`
	Parent        uuid.UUID   `json:"parent"`
	Children      []uuid.UUID `json:"children"`
	Tags          []string    `json:"tags"`
	Collaborators []uuid.UUID `json:"collaborators"`
	Icon          string      `json:"icon"`
	Header        string      `json:"header"`
	Favorite      bool        `json:"favorite"`
}

func (note *Note) Sanitize() {
	note.Data = html.EscapeString(note.Data)
}

type NoteUpdate struct {
	Doc Note `json:"doc"`
}

type UpsertNoteRequest struct {
	Data     interface{} `json:"data"`
	SocketID uuid.UUID   `json:"socket_id,omitempty"`
}

type AddCollaboratorRequest struct {
	Username string `json:"username"`
}

type SetIconRequest struct {
	Icon string `json:"icon"`
}

type SetHeaderRequest struct {
	Header string `json:"header"`
}

// ================================================================
// only swagger examples

type NoteDataForSwagger struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type NoteForSwagger struct {
	Id            uuid.UUID          `json:"id"`
	Data          NoteDataForSwagger `json:"data,omitempty"`
	CreateTime    time.Time          `json:"create_time"`
	UpdateTime    time.Time          `json:"update_time,omitempty"`
	OwnerId       uuid.UUID          `json:"owner_id"`
	Parent        uuid.UUID          `json:"parent"`
	Children      []uuid.UUID        `json:"children"`
	Tags          []string           `json:"tags"`
	Collaborators []uuid.UUID        `json:"collaborators"`
	Icon          string             `json:"icon"`
	Header        string             `json:"header"`
	Favorite      bool               `json:"favorite"`
}

type UpsertNoteRequestForSwagger struct {
	Data NoteDataForSwagger `json:"data"`
}
