package hub

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"

	"github.com/satori/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type HubInterface interface {
	StartCache(context.Context)
	StartCacheMain(ctx context.Context)
	WriteToCache(context.Context, models.CacheMessage)
	WriteToCacheMain(context.Context, uuid.UUID, models.InviteMessage)
	AddClient(context.Context, uuid.UUID, *CustomClient)
	AddClientMain(context.Context, uuid.UUID, *CustomClient)
	Run(context.Context)
}
