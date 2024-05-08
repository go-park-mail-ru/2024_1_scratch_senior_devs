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

type Hub struct {
	connect       sync.Map
	currentOffset time.Time
	repo          note.NoteBaseRepo
	cache         *ttlcache.Cache[uuid.UUID, models.CacheMessage]
	cfg           config.HubConfig
	metr          metrics.WSMetrics
}

func NewHub(repo note.NoteBaseRepo, cfg config.HubConfig, metr metrics.WSMetrics) *Hub {
	return &Hub{
		repo:          repo,
		currentOffset: time.Now().UTC(),
		cache:         ttlcache.New[uuid.UUID, models.CacheMessage](ttlcache.WithTTL[uuid.UUID, models.CacheMessage](cfg.CacheTtl)),
		cfg:           cfg,
		metr:          metr,
	}
}

func (h *Hub) StartCache(ctx context.Context) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	h.cache.Start()
	logger.Info("hub cache started")
}

func (h *Hub) WriteToCache(ctx context.Context, message models.CacheMessage) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	h.cache.Set(message.NoteId, message, h.cfg.CacheTtl)
	logger.Info("cache - new message")
}

func (h *Hub) AddClient(ctx context.Context, noteID uuid.UUID, client *websocket.Conn) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	h.connect.Store(client, noteID)
	h.metr.IncreaseConnections()

	go func() {
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
					connect := key.(*websocket.Conn)
					noteId := value.(uuid.UUID)

					if noteId == message.NoteId {
						if err := connect.WriteJSON(message); err != nil {
							logger.Error("can`t write hub`s message: " + err.Error())
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

func (h *Hub) Run(ctx context.Context) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	t := time.NewTicker(h.cfg.Period)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			h.connect.Range(func(key, value interface{}) bool {
				connect := key.(*websocket.Conn)
				noteID := value.(uuid.UUID)

				if h.cache.Has(noteID) {
					message := h.cache.Get(noteID).Value()

					if !message.Created.Before(h.currentOffset) {
						message.Type = "updated"

						if err := connect.WriteJSON(message); err != nil {
							logger.Error("can`t write hub`s message: " + err.Error())
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

		case <-ctx.Done():
			return
		}
	}
}
