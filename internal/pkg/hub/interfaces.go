package hub

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"

	"github.com/gorilla/websocket"
	"github.com/satori/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type HubInterface interface {
	StartCache(context.Context)
	WriteToCache(context.Context, models.CacheMessage)
	AddClient(context.Context, uuid.UUID, *websocket.Conn)
	Run(context.Context)
}
