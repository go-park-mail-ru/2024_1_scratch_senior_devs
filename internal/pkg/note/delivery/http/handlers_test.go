package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc/gen"
	mock_grpc "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc/gen/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNoteHandler_GetAllNotes(t *testing.T) {
	const successTestName = "Test Success"

	tests := []struct {
		name           string
		expectedStatus int
		id             uuid.UUID
		username       string
		expectedData   []*gen.NoteModel
	}{
		{

			name:           successTestName,
			id:             uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
			username:       "testuser",
			expectedStatus: http.StatusOK,
			expectedData: []*gen.NoteModel{
				{
					Id:         "c80e3ea8-0813-4731-b6ee-b41604c56f95",
					OwnerId:    "a233ea8-0813-4731-b12e-b41604c56f95",
					UpdateTime: time.Time{}.String(),
					CreateTime: time.Time{}.String(),
					Data:       "",
				}, {

					Id:         "c80e3ea8-0813-4731-b12e-b41604c56f95",
					OwnerId:    "a233ea8-0813-4731-b12e-b41604c56f95",
					UpdateTime: time.Time{}.String(),
					CreateTime: time.Time{}.String(),
					Data:       "",
				},
			},
		},
		{

			name:           "Test Unauthorized",
			id:             uuid.FromStringOrNil(""),
			username:       "",
			expectedStatus: http.StatusUnauthorized,
			expectedData:   []*gen.NoteModel{},
		},
		{

			name:           "Test Error",
			id:             uuid.FromStringOrNil(""),
			username:       "",
			expectedStatus: http.StatusBadRequest,
			expectedData:   []*gen.NoteModel{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			//mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			mockUsecase := mock_grpc.NewMockNoteClient(ctrl)
			defer ctrl.Finish()
			req := httptest.NewRequest("GET", "http://example.com/api/handler", nil)
			w := httptest.NewRecorder()
			ctx := context.WithValue(req.Context(), config.PayloadContextKey, models.JwtPayload{Id: tt.id, Username: tt.username})

			if tt.name == "Test Unauthorized" {
				ctx = context.WithValue(req.Context(), config.PayloadContextKey, &gen.GetAllResponse{})

			}
			if tt.name == successTestName {
				mockUsecase.EXPECT().GetAllNotes(ctx, &gen.GetAllRequest{Count: int64(10), Offset: int64(0), Title: "", UserId: tt.id.String()}, gomock.Any()).Return(&gen.GetAllResponse{Notes: tt.expectedData}, nil)
			}
			if tt.name == "Test Error" {
				mockUsecase.EXPECT().GetAllNotes(ctx, &gen.GetAllRequest{Count: int64(10), Offset: int64(0), Title: "", UserId: tt.id.String()}, gomock.Any()).Return(&gen.GetAllResponse{Notes: tt.expectedData}, errors.New("error"))

			}
			req = req.WithContext(ctx)

			h := CreateNotesHandler(mockUsecase)
			h.GetAllNotes(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.name == successTestName {
				expectedResult := make([]models.Note, len(tt.expectedData))
				for i, note := range tt.expectedData {
					crTime, _ := time.Parse("2006-01-02 15:04:05 -0700 UTC", note.CreateTime)
					upTime, _ := time.Parse("2006-01-02 15:04:05 -0700 UTC", note.UpdateTime)
					expectedResult[i] = models.Note{
						Id:         uuid.FromStringOrNil(note.Id),
						Data:       []byte(note.Data),
						CreateTime: crTime,
						UpdateTime: upTime,
						OwnerId:    uuid.FromStringOrNil(note.OwnerId),
						Children:   []uuid.UUID{},
					}

				}
				d, _ := json.Marshal(expectedResult)
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
			expectedData: models.Note{
				Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
				OwnerId:    uuid.FromStringOrNil(""),
				UpdateTime: time.Time{},
				CreateTime: time.Time{},
				Data:       []byte(""),
			},
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
			//mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			mockUsecase := mock_grpc.NewMockNoteClient(ctrl)
			defer ctrl.Finish()
			req := httptest.NewRequest("GET", "http://example.com/api/note/c80e3ea8-0813-4731-b6ee-b41604c56f95", nil)
			w := httptest.NewRecorder()
			ctx := context.WithValue(req.Context(), config.PayloadContextKey, models.JwtPayload{Id: tt.userId, Username: tt.username})
			if tt.name == "Test Unauthorized" {
				ctx = context.WithValue(req.Context(), config.PayloadContextKey, "")

			}
			if tt.name == successTestName {
				mockUsecase.EXPECT().GetNote(gomock.Any(), &gen.GetNoteRequest{
					Id:     tt.noteId.String(),
					UserId: tt.userId.String(),
				}).Return(&gen.GetNoteResponse{Note: &gen.NoteModel{}}, nil)
			}
			if tt.name == "Test Error" {
				mockUsecase.EXPECT().GetNote(gomock.Any(), &gen.GetNoteRequest{
					Id:     tt.noteId.String(),
					UserId: tt.userId.String(),
				}).Return(&gen.GetNoteResponse{}, errors.New("error"))

			}
			req = req.WithContext(ctx)
			if tt.name != "Test Bad Request" {
				req = mux.SetURLVars(req, map[string]string{"id": tt.noteId.String()})
			}
			h := CreateNotesHandler(mockUsecase)
			h.GetNote(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.name == successTestName {
				d, _ := json.Marshal(tt.expectedData)
				assert.Equal(t, w.Body.Bytes(), d)
			}
		})
	}
}

func TestNoteHandler_AddNote(t *testing.T) {
	id := uuid.NewV4()
	currTime := time.Now().UTC()

	var tests = []struct {
		name           string
		requestBody    string
		noteData       []byte
		usecaseErr     bool
		expectedStatus int
	}{
		{
			name:           "Test_NoteHandler_AddNote_Success",
			requestBody:    `{"data":{"title": "my note"}}`,
			noteData:       []byte(`{"title": "my note"}`),
			usecaseErr:     false,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Test_NoteHandler_AddNote_Fail_1",
			requestBody:    `{"data":{"title": "my note"}`,
			noteData:       []byte(`{"title": "my note"}`),
			usecaseErr:     true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Test_NoteHandler_AddNote_Fail_2",
			requestBody:    `{"data":{"title": "my note"}}`,
			noteData:       []byte(`{"title": "my note"}`),
			usecaseErr:     true,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			//mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			mockUsecase := mock_grpc.NewMockNoteClient(ctrl)
			defer ctrl.Finish()

			if tt.name != "Test_NoteHandler_AddNote_Fail_1" {
				call := mockUsecase.EXPECT().AddNote(gomock.Any(), gomock.Any(), gomock.Any())
				if tt.usecaseErr {
					call.Return(models.Note{}, errors.New("usecase error"))
				} else {
					call.Return(models.Note{
						Id:         id,
						Data:       tt.noteData,
						CreateTime: currTime,
						UpdateTime: currTime,
						OwnerId:    id,
					}, nil)
				}
			}

			r := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewBufferString(tt.requestBody))
			ctx := context.WithValue(r.Context(), config.PayloadContextKey, models.JwtPayload{
				Id:       id,
				Username: "username",
			})
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			handler := CreateNotesHandler(mockUsecase)
			handler.AddNote(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
