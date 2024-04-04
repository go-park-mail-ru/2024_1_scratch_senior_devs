package models

import (
	"time"

	"github.com/satori/uuid"
)

type Note struct {
	Id         uuid.UUID `json:"id"`
	Data       []byte    `json:"data,omitempty"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	OwnerId    uuid.UUID `json:"owner_id"`
}

type ElasticNote struct {
	Id          uuid.UUID   `json:"id"`
	Data        string      `json:"data,omitempty"`
	ElasticData interface{} `json:"elastic_data,omitempty"`
	CreateTime  time.Time   `json:"create_time"`
	UpdateTime  time.Time   `json:"update_time"`
	OwnerId     uuid.UUID   `json:"owner_id"`
}

type ElasticNoteUpdate struct {
	Doc ElasticNote `json:"doc"`
}

type UpsertNoteRequest struct {
	Data interface{} `json:"data"`
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
}

type UpsertNoteRequestForSwagger struct {
	Data NoteDataForSwagger `json:"data"`
}
