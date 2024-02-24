package models

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type Note struct {
	Id         uuid.UUID       `json:"id"`
	Data       json.RawMessage `json:"data,omitempty"`
	CreateTime time.Time       `json:"create_time"`
	UpdateTime time.Time       `json:"update_time,omitempty"`
	OwnerId    int             `json:"owner_id"`
}
