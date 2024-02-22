package models

import (
	"time"
)

type Note struct {
	Id         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content,omitempty"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	OwnerId    int       `json:"owner_id"`
}
