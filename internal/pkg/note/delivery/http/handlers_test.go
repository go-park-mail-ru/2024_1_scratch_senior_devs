package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	mock_note "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
)

var testLogger *slog.Logger

func init() {
	testLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func TestNoteHandler_GetAllNotes(t *testing.T) {
	const successTestName = "Test Success"

	tests := []struct {
		name           string
		expectedStatus int
		id             uuid.UUID
		username       string
		expectedData   []models.Note
	}{
		{

			name:           successTestName,
			id:             uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
			username:       "testuser",
			expectedStatus: http.StatusOK,
			expectedData: []models.Note{
				{
					Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
					OwnerId:    uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
					UpdateTime: time.Time{},
					CreateTime: time.Time{},
					Data:       nil,
				},
				{
					Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b12e-b41604c56f95"),
					OwnerId:    uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
					UpdateTime: time.Time{},
					CreateTime: time.Time{},
					Data:       nil,
				},
			},
		},
		{

			name:           "Test Unauthorized",
			id:             uuid.FromStringOrNil(""),
			username:       "",
			expectedStatus: http.StatusUnauthorized,
			expectedData:   []models.Note{},
		},
		{

			name:           "Test Error",
			id:             uuid.FromStringOrNil(""),
			username:       "",
			expectedStatus: http.StatusBadRequest,
			expectedData:   []models.Note{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			req := httptest.NewRequest("GET", "http://example.com/api/handler", nil)
			w := httptest.NewRecorder()
			ctx := context.WithValue(req.Context(), config.PayloadContextKey, models.JwtPayload{Id: tt.id, Username: tt.username})

			if tt.name == "Test Unauthorized" {
				ctx = context.WithValue(req.Context(), config.PayloadContextKey, models.Note{})

			}
			if tt.name == successTestName {
				mockUsecase.EXPECT().GetAllNotes(ctx, tt.id, int64(10), int64(0), "").Return(tt.expectedData, nil)
			}
			if tt.name == "Test Error" {
				mockUsecase.EXPECT().GetAllNotes(ctx, tt.id, int64(10), int64(0), "").Return(tt.expectedData, errors.New("error"))

			}
			req = req.WithContext(ctx)

			h := CreateNotesHandler(mockUsecase, testLogger)
			h.GetAllNotes(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.name == successTestName {
				d, _ := json.Marshal(tt.expectedData)
				assert.Equal(t, w.Body.Bytes(), d)
			}
		})
	}
}

func TestNoteHandler_GetNote(t *testing.T) {
	const successTestName = "Test Success"

	tests := []struct {
		name string

		expectedStatus int
		userId         uuid.UUID
		noteId         uuid.UUID
		username       string
		expectedData   models.Note
	}{
		// TODO: Add test cases.
		{
			name:           successTestName,
			expectedStatus: 200,
			noteId:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
			userId:         uuid.FromStringOrNil(""),
			username:       "test_user",
			expectedData:   models.Note{},
		},
		{
			name:           "Test Unauthorized",
			expectedStatus: 401,
			noteId:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
			userId:         uuid.FromStringOrNil(""),
			username:       "test_user",
			expectedData:   models.Note{},
		},
		{
			name:           "Test Error",
			expectedStatus: http.StatusNotFound,
			noteId:         uuid.FromStringOrNil(""),
			userId:         uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
			username:       "test_user",
			expectedData:   models.Note{},
		},
		{
			name:           "Test Bad Request",
			expectedStatus: http.StatusBadRequest,
			noteId:         uuid.FromStringOrNil(""),
			userId:         uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
			username:       "test_user",
			expectedData:   models.Note{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			req := httptest.NewRequest("GET", "http://example.com/api/note/c80e3ea8-0813-4731-b6ee-b41604c56f95", nil)
			w := httptest.NewRecorder()
			ctx := context.WithValue(req.Context(), config.PayloadContextKey, models.JwtPayload{Id: tt.userId, Username: tt.username})
			if tt.name == "Test Unauthorized" {
				ctx = context.WithValue(req.Context(), config.PayloadContextKey, models.Note{})

			}
			if tt.name == successTestName {
				mockUsecase.EXPECT().GetNote(gomock.Any(), tt.noteId, tt.userId).Return(tt.expectedData, nil)
			}
			if tt.name == "Test Error" {
				mockUsecase.EXPECT().GetNote(gomock.Any(), tt.noteId, tt.userId).Return(tt.expectedData, errors.New("error"))

			}
			req = req.WithContext(ctx)
			if tt.name != "Test Bad Request" {
				req = mux.SetURLVars(req, map[string]string{"id": tt.noteId.String()})
			}
			h := CreateNotesHandler(mockUsecase, testLogger)
			h.GetNote(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.name == successTestName {
				d, _ := json.Marshal(tt.expectedData)
				assert.Equal(t, w.Body.Bytes(), d)
			}
		})
	}
}
