package hub

import (
	"context"
	"encoding/json"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	mock_metrics "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/metrics/mocks"
	mock_note "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewHub(t *testing.T) {
	hubConfig := config.HubConfig{
		Period:   100 * time.Millisecond,
		CacheTtl: 1 * time.Minute,
	}

	t.Run("create hub test", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mock_note.NewMockNoteBaseRepo(ctrl)
		mockMetrics := mock_metrics.NewMockWSMetrics(ctrl)

		hub := NewHub(mockRepo, hubConfig, mockMetrics)
		assert.NotEqual(t, time.Time{}, hub.currentOffset)
	})
}

func TestHub_StartCache(t *testing.T) {
	hubConfig := config.HubConfig{
		Period:   100 * time.Millisecond,
		CacheTtl: 1 * time.Minute,
	}

	t.Run("hub --> start cache test", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mock_note.NewMockNoteBaseRepo(ctrl)
		mockMetrics := mock_metrics.NewMockWSMetrics(ctrl)

		hub := NewHub(mockRepo, hubConfig, mockMetrics)
		go hub.StartCache(context.Background())

		noteID := uuid.NewV4()
		hub.cache.Set(noteID, models.CacheMessage{
			NoteId:      noteID,
			Username:    "test",
			Created:     time.Now().UTC(),
			MessageInfo: []byte("{}"),
			Type:        "updated",
		}, hubConfig.CacheTtl)

		message := hub.cache.Get(noteID).Value()
		assert.Equal(t, "test", message.Username)

		hub.cache.Stop()
	})
}

func TestHub_WriteToCache(t *testing.T) {
	hubConfig := config.HubConfig{
		Period:   100 * time.Millisecond,
		CacheTtl: 1 * time.Minute,
	}

	t.Run("hub --> write to cache test", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mock_note.NewMockNoteBaseRepo(ctrl)
		mockMetrics := mock_metrics.NewMockWSMetrics(ctrl)

		hub := NewHub(mockRepo, hubConfig, mockMetrics)
		go hub.StartCache(context.Background())

		noteID := uuid.NewV4()
		hub.WriteToCache(context.Background(), models.CacheMessage{
			NoteId:      noteID,
			Username:    "test",
			Created:     time.Now().UTC(),
			MessageInfo: []byte("{}"),
			Type:        "updated",
		})

		message := hub.cache.Get(noteID).Value()
		assert.Equal(t, "test", message.Username)

		hub.cache.Stop()
	})
}

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

//func TestHub_AddClient(t *testing.T) {
//	hubConfig := config.HubConfig{
//		Period:   100 * time.Millisecond,
//		CacheTtl: 1 * time.Minute,
//	}
//
//	t.Run("hub --> add client test", func(t *testing.T) {
//		ctrl := gomock.NewController(t)
//		defer ctrl.Finish()
//		mockRepo := mock_note.NewMockNoteBaseRepo(ctrl)
//		mockMetrics := mock_metrics.NewMockWSMetrics(ctrl)
//
//		mockMetrics.EXPECT().IncreaseConnections().Return()
//		mockMetrics.EXPECT().DecreaseConnections().Return()
//
//		s := httptest.NewServer(http.HandlerFunc(echo))
//		defer s.Close()
//
//		u := "ws" + strings.TrimPrefix(s.URL, "http")
//		connection, res, err := websocket.DefaultDialer.Dial(u, nil)
//		if err != nil {
//			t.Fatalf("%v", err)
//		}
//		defer res.Body.Close()
//		defer connection.Close()
//
//		hub := NewHub(mockRepo, hubConfig, mockMetrics)
//		go hub.StartCache(context.Background())
//		defer hub.cache.Stop()
//
//		noteID := uuid.NewV4()
//		hub.AddClient(context.Background(), noteID, connection)
//
//		userID := uuid.NewV4()
//		joinMessage := models.JoinMessage{
//			Type:      "opened",
//			NoteId:    noteID,
//			UserId:    userID,
//			Username:  "test",
//			ImagePath: "default.jpg",
//		}
//
//		if err := connection.WriteJSON(joinMessage); err != nil {
//			t.Fatalf("%v", err)
//		}
//
//		_, byteMessage, err := connection.ReadMessage()
//		if err != nil {
//			t.Fatalf("%v", err)
//		}
//
//		received := models.JoinMessage{}
//		if err := json.Unmarshal(byteMessage, &received); err != nil {
//			t.Fatalf("%v", err)
//		}
//
//		assert.Equal(t, joinMessage, received)
//
//		// ==============================================================
//
//		if err := connection.WriteMessage(websocket.BinaryMessage, []byte("hello")); err != nil {
//			t.Fatalf("%v", err)
//		}
//	})
//}

func TestHub_Run(t *testing.T) {
	hubConfig := config.HubConfig{
		Period:   100 * time.Millisecond,
		CacheTtl: 1 * time.Minute,
	}

	t.Run("hub --> run test", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mock_note.NewMockNoteBaseRepo(ctrl)
		mockMetrics := mock_metrics.NewMockWSMetrics(ctrl)

		s := httptest.NewServer(http.HandlerFunc(echo))
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http")
		connection, res, err := websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			t.Fatalf("%v", err)
		}
		defer res.Body.Close()
		defer connection.Close()

		hub := NewHub(mockRepo, hubConfig, mockMetrics)
		go hub.StartCache(context.Background())
		defer hub.cache.Stop()

		noteID := uuid.NewV4()
		hub.connect.Store(connection, noteID)

		currentTime := time.Now().UTC()
		updateCacheMessage := models.CacheMessage{
			Type:     "updated",
			NoteId:   noteID,
			Username: "test",
			Created:  currentTime,
		}
		updateMessage := models.Message{
			Type:    "updated",
			NoteId:  noteID,
			Created: currentTime,
		}

		mockRepo.EXPECT().GetUpdates(gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Message{updateMessage}, nil).Times(1)

		go hub.Run(context.Background())

		hub.WriteToCache(context.Background(), updateCacheMessage)

		_, byteMessage, err := connection.ReadMessage()
		if err != nil {
			t.Fatalf("%v", err)
		}

		received := models.CacheMessage{}
		if err := json.Unmarshal(byteMessage, &received); err != nil {
			t.Fatalf("%v", err)
		}

		assert.Equal(t, updateCacheMessage.MessageInfo, received.MessageInfo)

		hub.cache.Delete(noteID)

		// =====================================================================

		_, byteMessage2, err := connection.ReadMessage()
		if err != nil {
			t.Fatalf("%v", err)
		}

		received2 := models.Message{}
		if err := json.Unmarshal(byteMessage2, &received2); err != nil {
			t.Fatalf("%v", err)
		}

		assert.Equal(t, updateMessage.MessageInfo, received2.MessageInfo)
	})
}
