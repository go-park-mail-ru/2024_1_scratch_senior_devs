package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	mock_note "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/mocks"
	"github.com/golang/mock/gomock"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNoteHandler_GetAllNotes(t *testing.T) {

	tests := []struct {
		name           string
		expectedStatus int
		id             uuid.UUID
		username       string
		expectedData   []models.Note
	}{
		{

			name:           "Test Success",
			id:             uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
			username:       "testuser",
			expectedStatus: http.StatusOK,
			expectedData: []models.Note{
				{
					Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
					OwnerId:    uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
					UpdateTime: nil,
					CreateTime: time.Time{},
					Data:       nil,
				},
				{
					Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b12e-b41604c56f95"),
					OwnerId:    uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
					UpdateTime: nil,
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
			ctx := context.WithValue(req.Context(), models.PayloadContextKey, models.JwtPayload{Id: tt.id, Username: tt.username})

			if tt.name == "Test Unauthorized" {
				ctx = context.WithValue(req.Context(), models.PayloadContextKey, models.Note{})

			}
			if tt.name == "Test Success" {
				mockUsecase.EXPECT().GetAllNotes(ctx, tt.id, gomock.Any(), gomock.Any()).Return(tt.expectedData, nil)
			}
			if tt.name == "Test Error" {
				mockUsecase.EXPECT().GetAllNotes(ctx, tt.id, gomock.Any(), gomock.Any()).Return(tt.expectedData, errors.New("error"))

			}
			req = req.WithContext(ctx)

			h := CreateNotesHandler(mockUsecase)
			h.GetAllNotes(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			//d, _ := json.Marshal(tt.expectedData)
			//assert.Equal(t, w.Body.Bytes(), d)
		})
	}
}
