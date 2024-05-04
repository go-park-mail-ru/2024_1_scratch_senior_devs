package models

import (
	"github.com/satori/uuid"
	"time"
)

type Message struct {
	NoteId      uuid.UUID `json:"note_id"`
	Created     time.Time `json:"created"`
	MessageInfo []byte    `json:"message_info"`
}

type CacheMessage struct {
	NoteId      uuid.UUID `json:"note_id"`
	Username    string    `json:"username"`
	Created     time.Time `json:"created"`
	MessageInfo []byte    `json:"message_info"`
}

type JoinMessage struct {
	Type      string    `json:"type"`
	NoteId    uuid.UUID `json:"note_id"`
	UserId    uuid.UUID `json:"user_id"`
	Username  string    `json:"username,omitempty"`
	ImagePath string    `json:"image_path,omitempty"`
}
