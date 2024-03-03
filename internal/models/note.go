package models

import (
	"time"

	"github.com/satori/uuid"
)

type Note struct {
	Id         uuid.UUID  `json:"id"`
	Data       []byte     `json:"data,omitempty"`
	CreateTime time.Time  `json:"create_time"`
	UpdateTime *time.Time `json:"update_time,omitempty"`
	OwnerId    uuid.UUID  `json:"owner_id"`
}
