package models

import (
	"html"
	"time"

	"github.com/satori/uuid"
)

type Note struct {
	Id         uuid.UUID   `json:"id"`
	Data       []byte      `json:"data,omitempty"`
	CreateTime time.Time   `json:"create_time"`
	UpdateTime time.Time   `json:"update_time"`
	OwnerId    uuid.UUID   `json:"owner_id"`
	Parent     uuid.UUID   `json:"parent"`
	Children   []uuid.UUID `json:"children"`
	Tags       []string    `json:"tags"`
}

func Sanitize(noteData []byte) []byte {
	return []byte(html.EscapeString(string(noteData)))
}

type ElasticNote struct {
	Id            uuid.UUID   `json:"id"`
	Data          string      `json:"data,omitempty"`
	CreateTime    time.Time   `json:"create_time"`
	UpdateTime    time.Time   `json:"update_time"`
	OwnerId       uuid.UUID   `json:"owner_id"`
	Parent        uuid.UUID   `json:"parent"`
	Children      []uuid.UUID `json:"children"`
	Collaborators []uuid.UUID `json:"collaborators"`
	//Tags          []string    `json:"tags"`
}

type NoteUpdate struct {
	Doc Note `json:"doc"`
}

type UpsertNoteRequest struct {
	Data interface{} `json:"data"`
}

type AddCollaboratorRequest struct {
	Username string `json:"username"`
}

// ================================================================
// only swagger examples

type NoteDataForSwagger struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type NoteForSwagger struct {
	Id         uuid.UUID          `json:"id"`
	Data       NoteDataForSwagger `json:"data,omitempty"`
	CreateTime time.Time          `json:"create_time"`
	UpdateTime time.Time          `json:"update_time,omitempty"`
	OwnerId    uuid.UUID          `json:"owner_id"`
	Parent     uuid.UUID          `json:"parent"`
	Children   []uuid.UUID        `json:"children"`
}

type UpsertNoteRequestForSwagger struct {
	Data NoteDataForSwagger `json:"data"`
}
