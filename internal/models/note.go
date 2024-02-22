package models

import (
	"encoding/json"
	"time"
)

type Note struct {
	Id         int             `json:"id"`
	Data       json.RawMessage `json:"data"`
	CreateTime time.Time       `json:"create_time"`
	UpdateTime time.Time       `json:"update_time,omitempty"`
	OwnerId    int             `json:"owner_id"`
}
