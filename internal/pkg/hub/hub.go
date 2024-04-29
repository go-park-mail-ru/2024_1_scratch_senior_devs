package hub

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"
	"github.com/gorilla/websocket"
	"github.com/satori/uuid"
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

func (h *Hub) AddClient(noteID uuid.UUID, client *websocket.Conn) {
	h.connect.Store(client, noteID)

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
	t := time.NewTicker(h.cfg.Period)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			h.connect.Range(func(key, value interface{}) bool {
				connect := key.(*websocket.Conn)
				noteID := value.(uuid.UUID)

				messages, _ := h.repo.GetUpdates(ctx, noteID, h.currentOffset)
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
