package models

import (
	"time"

	"github.com/satori/uuid"
)

type Message struct {
	NoteId      uuid.UUID `json:"note_id"`
	Created     time.Time `json:"created"`
	MessageInfo string    `json:"message_info"`
	Type        string    `json:"type" default:"updated"`
}

type CacheMessage struct {
	Type        string    `json:"type" default:"updated"`
	NoteId      uuid.UUID `json:"note_id"`
	Username    string    `json:"username"`
	Created     time.Time `json:"created"`
	MessageInfo string    `json:"message_info"`
	SocketID    uuid.UUID `json:"socket_id"`
}

type JoinMessage struct {
	Type      string    `json:"type"`
	NoteId    uuid.UUID `json:"note_id"`
	UserId    uuid.UUID `json:"user_id"`
	Username  string    `json:"username,omitempty"`
	ImagePath string    `json:"image_path,omitempty"`
}

type InviteMessage struct {
	Type      string    `json:"type"`
	NoteId    uuid.UUID `json:"note_id"`
	NoteTitle string    `json:"note_title"`
	Owner     string    `json:"owner"`
	Created   time.Time `json:"created"`
}

type SocketIDMessage struct {
	Type     string    `json:"type"`
	SocketID uuid.UUID `json:"socket_id"`
}
