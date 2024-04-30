package hub

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/satori/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type HubInterface interface {
	AddClient(context.Context, uuid.UUID, *websocket.Conn)
	Run(context.Context)
}
