package hub

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/gorilla/websocket"
	"github.com/satori/uuid"
	"log/slog"
	"sync"
	"time"
)

type Hub struct {
	connect       sync.Map
	currentOffset time.Time
	repo          note.NoteBaseRepo
	cfg           config.HubConfig
}

func NewHub(repo note.NoteBaseRepo, cfg config.HubConfig) *Hub {
	return &Hub{
		repo:          repo,
		currentOffset: time.Now().UTC(),
		cfg:           cfg,
	}
}

func (h *Hub) AddClient(ctx context.Context, noteID uuid.UUID, client *websocket.Conn) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	h.connect.Store(client, noteID)
	h.connect.Range(func(key, value interface{}) bool {
		logger.Info("add client, note_id = " + value.(uuid.UUID).String())
		return true
	})

	go func() {
		for {
			_, _, err := client.NextReader()
			if err != nil {
				_ = client.Close()
				return
			}
		}
	}()

	client.SetCloseHandler(func(code int, text string) error {
		h.connect.Delete(client)
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

				logger.Info("run, note_id = " + noteID.String())

				messages, err := h.repo.GetUpdates(ctx, noteID, h.currentOffset)
				if err != nil {
					logger.Error(err.Error())
				}

				logger.Info(fmt.Sprintf("messages found: %+v\n", messages))

				for _, message := range messages {
					err := connect.WriteJSON(message)
					if err != nil {
						continue
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
