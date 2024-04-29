package hub

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/satori/uuid"
)

type HubInterface interface {
	AddClient(context.Context, uuid.UUID, *websocket.Conn)
	Run(context.Context)
}
