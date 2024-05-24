package hub

import (
	"context"
	"encoding/json"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/metrics"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/gorilla/websocket"
	"github.com/jellydator/ttlcache/v3"
	"github.com/satori/uuid"
)

const ErrHubWrite = "can`t write hub`s message: "

type Hub struct {
	connect       sync.Map // wsConnection : noteID
	currentOffset time.Time
	cache         *ttlcache.Cache[uuid.UUID, models.CacheMessage] // noteID : message

	connectMain       sync.Map // wsConnection : userID
	currentOffsetMain time.Time
	cacheMain         *ttlcache.Cache[uuid.UUID, []models.InviteMessage] // userID : []message

	repo note.NoteBaseRepo
	cfg  config.HubConfig
	metr metrics.WSMetrics
}

func NewHub(repo note.NoteBaseRepo, cfg config.HubConfig, metr metrics.WSMetrics) *Hub {
	return &Hub{
		currentOffset: time.Now().UTC(),
		cache:         ttlcache.New[uuid.UUID, models.CacheMessage](ttlcache.WithTTL[uuid.UUID, models.CacheMessage](cfg.CacheTtl)),

		currentOffsetMain: time.Now().UTC(),
		cacheMain:         ttlcache.New[uuid.UUID, []models.InviteMessage](ttlcache.WithTTL[uuid.UUID, []models.InviteMessage](cfg.CacheTtl)),

		repo: repo,
		cfg:  cfg,
		metr: metr,
	}
}

func (h *Hub) StartCache(ctx context.Context) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	h.cache.Start()

	logger.Info("hub cache started")
}

func (h *Hub) StartCacheMain(ctx context.Context) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	h.cacheMain.Start()

	logger.Info("hub cacheMain started")
}

func (h *Hub) WriteToCache(ctx context.Context, message models.CacheMessage) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	h.cache.Set(message.NoteId, message, h.cfg.CacheTtl)
	logger.Info("cache - new message")
}

func (h *Hub) WriteToCacheMain(ctx context.Context, invitedID uuid.UUID, message models.InviteMessage) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	h.cacheMain.Set(invitedID, append(h.cacheMain.Get(invitedID).Value(), message), h.cfg.CacheTtl)

	logger.Info("cacheMain - new message")
}

type CustomClient struct {
	*websocket.Conn
	SocketID uuid.UUID
}

func NewCustomClient(connection *websocket.Conn) *CustomClient {
	return &CustomClient{
		Conn:     connection,
		SocketID: uuid.UUID{},
	}
}

func (h *Hub) AddClient(ctx context.Context, noteID uuid.UUID, client *CustomClient) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	client.SocketID = uuid.NewV4()
	h.connect.Store(client, noteID)
	h.metr.IncreaseConnections()

	go func() {
		if err := client.WriteJSON(models.SocketIDMessage{
			Type:     "info",
			SocketID: client.SocketID,
		}); err != nil {
			logger.Error(ErrHubWrite + err.Error())
		}

		for {
			messageType, reader, err := client.NextReader()
			if err != nil {
				_ = client.Close()
				h.metr.DecreaseConnections()
				return
			}

			switch messageType {
			case websocket.TextMessage:
				messageBytes, err := io.ReadAll(reader)
				if err != nil {
					logger.Error("failed to read message: " + err.Error())
					continue
				}

				var message models.JoinMessage
				if err := json.Unmarshal(messageBytes, &message); err != nil {
					logger.Error("incorrect message format: " + err.Error())
					continue
				}

				h.connect.Range(func(key, value interface{}) bool {
					connect := key.(*CustomClient)
					noteId := value.(uuid.UUID)

					if noteId == message.NoteId {
						if err := connect.WriteJSON(message); err != nil {
							logger.Error(ErrHubWrite + err.Error())
						}
					}

					return true
				})

			default:
				logger.Error("received unsupported message type")
			}
		}
	}()

	client.SetCloseHandler(func(code int, text string) error {
		h.connect.Delete(client)
		h.metr.DecreaseConnections()
		return nil
	})
}

func (h *Hub) AddClientMain(ctx context.Context, invitedID uuid.UUID, client *CustomClient) {
	client.SocketID = uuid.NewV4()
	h.connectMain.Store(client, invitedID)
	h.metr.IncreaseConnections()

	go func() {
		for {
			_, _, err := client.NextReader()
			if err != nil {
				_ = client.Close()
				h.metr.DecreaseConnections()
				return
			}
		}
	}()

	client.SetCloseHandler(func(code int, text string) error {
		h.connect.Delete(client)
		h.metr.DecreaseConnections()
		return nil
	})
}

func (h *Hub) Run(ctx context.Context) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	t := time.NewTicker(h.cfg.Period)
	defer t.Stop()

	tMain := time.NewTicker(h.cfg.Period)
	defer tMain.Stop()

	for {
		select {
		case <-t.C:
			h.connect.Range(func(key, value interface{}) bool {
				connect := key.(*CustomClient)
				noteID := value.(uuid.UUID)

				if h.cache.Has(noteID) {
					message := h.cache.Get(noteID).Value()

					if !message.Created.Before(h.currentOffset) {
						message.Type = "updated"

						if err := connect.WriteJSON(message); err != nil {
							logger.Error(ErrHubWrite + err.Error())
						}
					}
				} else {
					messages, err := h.repo.GetUpdates(ctx, noteID, h.currentOffset)
					if err != nil {
						logger.Error(err.Error())
					}

					for _, message := range messages {
						message.Type = "updated"
						if err := connect.WriteJSON(message); err != nil {
							continue
						}
					}
				}

				return true
			})

			h.currentOffset = h.currentOffset.Add(h.cfg.Period)

		case <-tMain.C:
			h.connectMain.Range(func(key, value interface{}) bool {
				connect := key.(*CustomClient)
				invitedID := value.(uuid.UUID)

				if h.cacheMain.Has(invitedID) {
					inviteMessages := h.cacheMain.Get(invitedID).Value()
					for _, message := range inviteMessages {
						if !message.Created.Before(h.currentOffsetMain) {
							if err := connect.WriteJSON(message); err != nil {
								logger.Error(ErrHubWrite + err.Error())
							}
						}
					}
				}

				return true
			})

			h.currentOffsetMain = h.currentOffsetMain.Add(h.cfg.Period)

		case <-ctx.Done():
			return
		}
	}
}
