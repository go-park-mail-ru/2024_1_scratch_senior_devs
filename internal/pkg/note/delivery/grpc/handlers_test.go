package grpc

import (
	"bytes"
	"context"
	"errors"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc/gen"
	mock_note "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/mocks"
	"github.com/golang/mock/gomock"
	"github.com/satori/uuid"
)

func TestNoteHandler_GetAllNotes(t *testing.T) {
	const successTestName = "Test Success"

	tests := []struct {
		name         string
		wantErr      bool
		id           uuid.UUID
		username     string
		expectedData *gen.GetAllResponse
	}{
		{

			name:     successTestName,
			id:       uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
			username: "testuser",
			wantErr:  false,
			expectedData: &gen.GetAllResponse{Notes: []*gen.NoteModel{
				{
					Id:         "c80e3ea8-0813-4731-b6ee-b41604c56f95",
					OwnerId:    "a89e3ea8-0813-4731-b6ee-b41604c56f95",
					UpdateTime: "0001-01-01 00:00:00 +0000 UTC",
					CreateTime: "0001-01-01 00:00:00 +0000 UTC",
					Data:       "nil",
				},
				{
					Id:         "c80e3ea8-0813-4731-b12e-b41604c56f95",
					OwnerId:    "a89e3ea8-0813-4731-b6ee-b41604c56f95",
					UpdateTime: "0001-01-01 00:00:00 +0000 UTC",
					CreateTime: "0001-01-01 00:00:00 +0000 UTC",
					Data:       "nil",
				},
			},
			}},

		{

			name:         "Test Error",
			id:           uuid.FromStringOrNil(""),
			username:     "",
			wantErr:      true,
			expectedData: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			req := httptest.NewRequest("GET", "http://example.com/api/handler", nil)
			ctx := context.WithValue(req.Context(), config.PayloadContextKey, models.JwtPayload{Id: tt.id, Username: tt.username})

			if tt.name == successTestName {
				mockUsecase.EXPECT().GetAllNotes(ctx, tt.id, int64(10), int64(0), "").Return([]models.Note{
					{
						Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
						OwnerId:    uuid.FromStringOrNil("a89e3ea8-0813-4731-b6ee-b41604c56f95"),
						UpdateTime: time.Time{},
						CreateTime: time.Time{},
						Data:       []byte("nil"),
					},
					{
						Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b12e-b41604c56f95"),
						OwnerId:    uuid.FromStringOrNil("a89e3ea8-0813-4731-b6ee-b41604c56f95"),
						UpdateTime: time.Time{},
						CreateTime: time.Time{},
						Data:       []byte("nil"),
					},
				}, nil)
			}
			if tt.name == "Test Error" {
				mockUsecase.EXPECT().GetAllNotes(ctx, tt.id, int64(10), int64(0), "").Return([]models.Note{}, errors.New("error"))

			}
			req = req.WithContext(ctx)

			h := NewGrpcNoteHandler(mockUsecase)
			got, err := h.GetAllNotes(req.Context(), &gen.GetAllRequest{
				UserId: tt.id.String(),
				Count:  10,
				Offset: 0,
				Title:  "",
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllNotes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expectedData) {
				t.Errorf("GetAllNotes() = %v, want %v", got, tt.expectedData)
			}

		})
	}
}

func TestNoteHandler_GetNote(t *testing.T) {
	const successTestName = "Test Success"

	tests := []struct {
		name string

		wantErr      bool
		userId       uuid.UUID
		noteId       uuid.UUID
		username     string
		expectedData *gen.GetNoteResponse
	}{
		// TODO: Add test cases.
		{
			name:     successTestName,
			wantErr:  false,
			noteId:   uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
			userId:   uuid.FromStringOrNil("a89e3ea8-0813-4731-b6ee-b41604c56f95"),
			username: "test_user",
			expectedData: &gen.GetNoteResponse{
				Note: &gen.NoteModel{
					Id:         "c80e3ea8-0813-4731-b6ee-b41604c56f95",
					OwnerId:    "a89e3ea8-0813-4731-b6ee-b41604c56f95",
					Data:       "",
					CreateTime: "0001-01-01 00:00:00 +0000 UTC",
					UpdateTime: "0001-01-01 00:00:00 +0000 UTC",
				},
			},
		},
		// {
		// 	name:           "Test Unauthorized",
		// 	expectedStatus: 401,
		// 	noteId:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
		// 	userId:         uuid.FromStringOrNil(""),
		// 	username:       "test_user",
		// 	expectedData:   models.Note{},
		// },
		{
			name:         "Test Error",
			wantErr:      true,
			noteId:       uuid.FromStringOrNil(""),
			userId:       uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
			username:     "test_user",
			expectedData: nil,
		},
		// {
		// 	name:         "Test Bad Request",
		// 	wantErr:      true,
		// 	noteId:       uuid.FromStringOrNil(""),
		// 	userId:       uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
		// 	username:     "test_user",
		// 	expectedData: nil,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			req := httptest.NewRequest("GET", "http://example.com/api/note/c80e3ea8-0813-4731-b6ee-b41604c56f95", nil)

			ctx := context.WithValue(req.Context(), config.PayloadContextKey, models.JwtPayload{Id: tt.userId, Username: tt.username})

			if tt.name == successTestName {
				mockUsecase.EXPECT().GetNote(gomock.Any(), tt.noteId, tt.userId).Return(models.Note{
					Id:         uuid.FromStringOrNil(tt.expectedData.Note.Id),
					OwnerId:    uuid.FromStringOrNil(tt.expectedData.Note.OwnerId),
					CreateTime: time.Time{},
					UpdateTime: time.Time{},
					Data:       []byte(tt.expectedData.Note.Data),
				}, nil)
			}
			if tt.name == "Test Error" {
				mockUsecase.EXPECT().GetNote(gomock.Any(), tt.noteId, tt.userId).Return(models.Note{}, errors.New("error"))

			}
			req = req.WithContext(ctx)

			h := NewGrpcNoteHandler(mockUsecase)
			got, err := h.GetNote(req.Context(), &gen.GetNoteRequest{
				Id:     tt.noteId.String(),
				UserId: tt.userId.String(),
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("GetNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expectedData) {
				t.Errorf("GetNote() = %v, want %v", got, tt.expectedData)
			}
		})
	}
}
func TestNoteHandler_AddNote(t *testing.T) {
	id := uuid.NewV4()
	currTime := time.Now().UTC()

	var tests = []struct {
		name        string
		requestBody *gen.AddNoteRequest

		usecaseErr bool
		wantErr    bool
		want       *gen.AddNoteResponse
	}{
		{
			name: "Test_NoteHandler_AddNote_Success",
			requestBody: &gen.AddNoteRequest{
				Data:   `{"title": "my note"}`,
				UserId: id.String(),
			},
			wantErr:    false,
			usecaseErr: false,
			want: &gen.AddNoteResponse{
				Note: &gen.NoteModel{
					Id:         id.String(),
					CreateTime: currTime.String(),
					UpdateTime: currTime.String(),
					OwnerId:    id.String(),
					Data:       `{"title": "my note"}`,
				},
			},
		},
		// {
		// 	name: "Test_NoteHandler_AddNote_Fail_1",
		// 	requestBody: &gen.AddNoteRequest{
		// 		Data:   `{"data":{"title": "my note"}`,
		// 		UserId: id.String(),
		// 	},
		// 	usecaseErr: true,
		// 	wantErr:    true,
		// 	want:       nil,
		// },
		{
			name: "Test_NoteHandler_AddNote_Fail_2",
			requestBody: &gen.AddNoteRequest{
				Data:   `{"title": "my note"}`,
				UserId: id.String(),
			},

			usecaseErr: true,
			wantErr:    true,
			want:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()

			if tt.name != "Test_NoteHandler_AddNote_Fail_1" {
				call := mockUsecase.EXPECT().CreateNote(gomock.Any(), gomock.Any(), gomock.Any())
				if tt.usecaseErr {
					call.Return(models.Note{}, errors.New("usecase error"))
				} else {
					call.Return(models.Note{
						Id:         id,
						Data:       []byte(tt.requestBody.Data),
						CreateTime: currTime,
						UpdateTime: currTime,
						OwnerId:    id,
					}, nil)
				}
			}

			r := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewBufferString(tt.requestBody.Data))
			ctx := context.WithValue(r.Context(), config.PayloadContextKey, models.JwtPayload{
				Id:       id,
				Username: "username",
			})
			r = r.WithContext(ctx)

			handler := NewGrpcNoteHandler(mockUsecase)
			got, err := handler.AddNote(r.Context(), &gen.AddNoteRequest{
				UserId: id.String(),
				Data:   string(tt.requestBody.Data),
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("GetNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNote() = %v, want %v", got, tt.want)
			}
		})
	}
}
