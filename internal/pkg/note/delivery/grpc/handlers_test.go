package grpc

import (
	"bytes"
	"context"
	"errors"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc/gen"
	generatedNote "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc/gen"
	mock_note "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/mocks"
	"github.com/golang/mock/gomock"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
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
					Id:            "c80e3ea8-0813-4731-b6ee-b41604c56f95",
					OwnerId:       "a89e3ea8-0813-4731-b6ee-b41604c56f95",
					UpdateTime:    "0001-01-01 00:00:00 +0000 UTC",
					CreateTime:    "0001-01-01 00:00:00 +0000 UTC",
					Data:          "nil",
					Parent:        "c80e3ea8-0813-4731-b6ee-b41604c56f95",
					Children:      []string{},
					Tags:          []string{},
					Collaborators: []string{},
				},
				{
					Id:            "c80e3ea8-0813-4731-b12e-b41604c56f95",
					OwnerId:       "a89e3ea8-0813-4731-b6ee-b41604c56f95",
					UpdateTime:    "0001-01-01 00:00:00 +0000 UTC",
					CreateTime:    "0001-01-01 00:00:00 +0000 UTC",
					Data:          "nil",
					Parent:        "c80e3ea8-0813-4731-b6ee-b41604c56f95",
					Children:      []string{},
					Tags:          []string{},
					Collaborators: []string{},
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
				mockUsecase.EXPECT().GetAllNotes(ctx, tt.id, int64(10), int64(0), "", gomock.Any()).Return([]models.Note{
					{
						Id:            uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
						OwnerId:       uuid.FromStringOrNil("a89e3ea8-0813-4731-b6ee-b41604c56f95"),
						UpdateTime:    time.Time{},
						CreateTime:    time.Time{},
						Data:          "nil",
						Parent:        uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
						Children:      []uuid.UUID{},
						Tags:          []string{},
						Collaborators: []uuid.UUID{},
					},
					{
						Id:            uuid.FromStringOrNil("c80e3ea8-0813-4731-b12e-b41604c56f95"),
						OwnerId:       uuid.FromStringOrNil("a89e3ea8-0813-4731-b6ee-b41604c56f95"),
						UpdateTime:    time.Time{},
						CreateTime:    time.Time{},
						Data:          "nil",
						Parent:        uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
						Children:      []uuid.UUID{},
						Tags:          []string{},
						Collaborators: []uuid.UUID{},
					},
				}, nil)
			}
			if tt.name == "Test Error" {
				mockUsecase.EXPECT().GetAllNotes(ctx, tt.id, int64(10), int64(0), "", gomock.Any()).Return([]models.Note{}, errors.New("error"))

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
					Id:            "c80e3ea8-0813-4731-b6ee-b41604c56f95",
					OwnerId:       "a89e3ea8-0813-4731-b6ee-b41604c56f95",
					Data:          "",
					CreateTime:    "0001-01-01 00:00:00 +0000 UTC",
					UpdateTime:    "0001-01-01 00:00:00 +0000 UTC",
					Parent:        "c80e3ea8-0813-4731-b6ee-b41604c56f95",
					Children:      []string{},
					Tags:          []string{},
					Collaborators: []string{},
				},
			},
		},

		{
			name:         "Test Error",
			wantErr:      true,
			noteId:       uuid.FromStringOrNil(""),
			userId:       uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
			username:     "test_user",
			expectedData: nil,
		},
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
					Id:            uuid.FromStringOrNil(tt.expectedData.Note.Id),
					OwnerId:       uuid.FromStringOrNil(tt.expectedData.Note.OwnerId),
					CreateTime:    time.Time{},
					UpdateTime:    time.Time{},
					Data:          tt.expectedData.Note.Data,
					Parent:        uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
					Children:      []uuid.UUID{},
					Tags:          []string{},
					Collaborators: []uuid.UUID{},
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
					Id:            id.String(),
					CreateTime:    currTime.String(),
					UpdateTime:    currTime.String(),
					OwnerId:       id.String(),
					Data:          `{"title": "my note"}`,
					Parent:        id.String(),
					Children:      []string{},
					Tags:          []string{},
					Collaborators: []string{},
				},
			},
		},

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
						Id:            id,
						Data:          tt.requestBody.Data,
						CreateTime:    currTime,
						UpdateTime:    currTime,
						OwnerId:       id,
						Parent:        id,
						Children:      []uuid.UUID{},
						Tags:          []string{},
						Collaborators: []uuid.UUID{},
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

func TestGrpcNoteHandler_GetTags(t *testing.T) {

	tests := []struct {
		name        string
		requestBody *gen.GetTagsRequest
		wantErr     bool
		mocker      func(userId uuid.UUID, mock *mock_note.MockNoteUsecase)
		want        *gen.GetTagsResponse
	}{
		{
			name:        "Test_GetTags_Success",
			requestBody: &gen.GetTagsRequest{UserId: uuid.NewV4().String()},
			wantErr:     false,
			mocker: func(userId uuid.UUID, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().GetTags(gomock.Any(), userId).Return([]string{"first", "second"}, nil)
			},
			want: &gen.GetTagsResponse{
				Tags: []string{"first", "second"},
			},
		},
		{
			name:        "Test_GetTags_Fail",
			requestBody: &gen.GetTagsRequest{UserId: uuid.NewV4().String()},
			wantErr:     true,
			mocker: func(userId uuid.UUID, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().GetTags(gomock.Any(), userId).Return([]string{}, errors.New("error"))
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			ctx := context.Background()

			h := NewGrpcNoteHandler(mockUsecase)
			tt.mocker(uuid.FromStringOrNil(tt.requestBody.UserId), mockUsecase)
			got, err := h.GetTags(ctx, tt.requestBody)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetTags error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTags = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestGrpcNoteHandler_AddTag(t *testing.T) {
	noteId := uuid.NewV4()
	userId := uuid.NewV4()

	tests := []struct {
		name        string
		requestBody *gen.TagRequest
		mocker      func(req *gen.TagRequest, mock *mock_note.MockNoteUsecase)
		want        *gen.GetNoteResponse
		wantErr     bool
	}{
		{
			name: "Test_AddTag_Success",
			requestBody: &gen.TagRequest{
				TagName: "tag",
				NoteId:  noteId.String(),
				UserId:  userId.String(),
			},
			want: &gen.GetNoteResponse{
				Note: &gen.NoteModel{
					Id:            noteId.String(),
					OwnerId:       userId.String(),
					Tags:          []string{"tag"},
					CreateTime:    time.Time{}.String(),
					UpdateTime:    time.Time{}.String(),
					Parent:        uuid.UUID{}.String(),
					Collaborators: []string{},
					Children:      []string{},
				},
			},
			wantErr: false,
			mocker: func(req *gen.TagRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().AddTag(gomock.Any(), req.TagName,
					uuid.FromStringOrNil(req.NoteId), uuid.FromStringOrNil(req.UserId)).Return(models.Note{
					Id:      noteId,
					OwnerId: userId,
					Tags:    []string{"tag"},
				}, nil)
			},
		},
		{
			name: "Test_AddTag_Fail",
			requestBody: &gen.TagRequest{
				TagName: "tag",
				NoteId:  noteId.String(),
				UserId:  userId.String(),
			},
			want:    nil,
			wantErr: true,
			mocker: func(req *gen.TagRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().AddTag(gomock.Any(), req.TagName,
					uuid.FromStringOrNil(req.NoteId), uuid.FromStringOrNil(req.UserId)).Return(models.Note{}, errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			ctx := context.Background()

			h := NewGrpcNoteHandler(mockUsecase)
			tt.mocker(tt.requestBody, mockUsecase)

			got, err := h.AddTag(ctx, tt.requestBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcNoteHandler.AddTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGrpcNoteHandler_DeleteTag(t *testing.T) {
	noteId := uuid.NewV4()
	userId := uuid.NewV4()

	tests := []struct {
		name        string
		requestBody *gen.TagRequest
		mocker      func(req *gen.TagRequest, mock *mock_note.MockNoteUsecase)
		want        *gen.GetNoteResponse
		wantErr     bool
	}{
		{
			name: "Test_DeleteTag_Success",
			requestBody: &gen.TagRequest{
				TagName: "tag",
				NoteId:  noteId.String(),
				UserId:  userId.String(),
			},
			want: &gen.GetNoteResponse{
				Note: &gen.NoteModel{
					Id:            noteId.String(),
					OwnerId:       userId.String(),
					Tags:          []string{},
					CreateTime:    time.Time{}.String(),
					UpdateTime:    time.Time{}.String(),
					Parent:        uuid.UUID{}.String(),
					Collaborators: []string{},
					Children:      []string{},
				},
			},
			wantErr: false,
			mocker: func(req *gen.TagRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().DeleteTag(gomock.Any(), req.TagName,
					uuid.FromStringOrNil(req.NoteId), uuid.FromStringOrNil(req.UserId)).Return(models.Note{
					Id:      noteId,
					OwnerId: userId,
					Tags:    []string{},
				}, nil)
			},
		},
		{
			name: "Test_DeleteTag_Fail",
			requestBody: &gen.TagRequest{
				TagName: "tag",
				NoteId:  noteId.String(),
				UserId:  userId.String(),
			},
			want:    nil,
			wantErr: true,
			mocker: func(req *gen.TagRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().DeleteTag(gomock.Any(), req.TagName,
					uuid.FromStringOrNil(req.NoteId), uuid.FromStringOrNil(req.UserId)).Return(models.Note{}, errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			ctx := context.Background()

			h := NewGrpcNoteHandler(mockUsecase)
			tt.mocker(tt.requestBody, mockUsecase)

			got, err := h.DeleteTag(ctx, tt.requestBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcNoteHandler.AddTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGrpcNoteHandler_AddCollaborator(t *testing.T) {
	noteId := uuid.NewV4()
	userId := uuid.NewV4()
	guestId := uuid.NewV4()

	tests := []struct {
		name        string
		requestBody *gen.AddCollaboratorRequest
		mocker      func(req *gen.AddCollaboratorRequest, mock *mock_note.MockNoteUsecase)
		want        *gen.AddCollaboratorResponse
		wantErr     bool
	}{
		{
			name: "Test_AddCollaborator_Success",
			requestBody: &gen.AddCollaboratorRequest{
				NoteId:  noteId.String(),
				UserId:  userId.String(),
				GuestId: guestId.String(),
			},
			want:    &gen.AddCollaboratorResponse{},
			wantErr: false,
			mocker: func(req *gen.AddCollaboratorRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().AddCollaborator(gomock.Any(), uuid.FromStringOrNil(req.NoteId), uuid.FromStringOrNil(req.UserId), uuid.FromStringOrNil(req.GuestId)).Return("", nil)
			},
		},
		{
			name: "Test_AddCollaborator_Fail",
			requestBody: &gen.AddCollaboratorRequest{
				NoteId:  noteId.String(),
				UserId:  userId.String(),
				GuestId: guestId.String(),
			},
			want:    nil,
			wantErr: true,
			mocker: func(req *gen.AddCollaboratorRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().AddCollaborator(gomock.Any(), uuid.FromStringOrNil(req.NoteId), uuid.FromStringOrNil(req.UserId), uuid.FromStringOrNil(req.GuestId)).Return("", errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			ctx := context.Background()

			h := NewGrpcNoteHandler(mockUsecase)
			tt.mocker(tt.requestBody, mockUsecase)

			got, err := h.AddCollaborator(ctx, tt.requestBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcNoteHandler.AddTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGrpcNoteHandler_DeleteNote(t *testing.T) {
	noteId := uuid.NewV4()
	userId := uuid.NewV4()

	tests := []struct {
		name        string
		requestBody *gen.DeleteNoteRequest
		mocker      func(req *gen.DeleteNoteRequest, mock *mock_note.MockNoteUsecase)
		want        *gen.DeleteNoteResponse
		wantErr     bool
	}{
		{
			name: "Test_DeleteNote_Success",
			requestBody: &gen.DeleteNoteRequest{
				Id:     noteId.String(),
				UserId: userId.String(),
			},
			want:    &gen.DeleteNoteResponse{},
			wantErr: false,
			mocker: func(req *gen.DeleteNoteRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().DeleteNote(gomock.Any(), uuid.FromStringOrNil(req.Id), uuid.FromStringOrNil(req.UserId)).Return(nil)
			},
		},

		{
			name: "Test_DeleteNote_Fail",
			requestBody: &gen.DeleteNoteRequest{
				Id:     noteId.String(),
				UserId: userId.String(),
			},
			want:    nil,
			wantErr: true,
			mocker: func(req *gen.DeleteNoteRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().DeleteNote(gomock.Any(), uuid.FromStringOrNil(req.Id), uuid.FromStringOrNil(req.UserId)).Return(errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			ctx := context.Background()

			h := NewGrpcNoteHandler(mockUsecase)
			tt.mocker(tt.requestBody, mockUsecase)

			got, err := h.DeleteNote(ctx, tt.requestBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcNoteHandler.DeleteNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGrpcNoteHandler_CreateSubNote(t *testing.T) {
	noteId := uuid.NewV4()
	userId := uuid.NewV4()
	childId := uuid.NewV4()
	tests := []struct {
		name        string
		requestBody *gen.CreateSubNoteRequest
		mocker      func(req *gen.CreateSubNoteRequest, mock *mock_note.MockNoteUsecase)
		want        *gen.CreateSubNoteResponse
		wantErr     bool
	}{
		{
			name: "Test_CreateSubNote_Success",
			requestBody: &gen.CreateSubNoteRequest{
				UserId:   userId.String(),
				NoteData: "",
				ParentId: noteId.String(),
			},
			want: &gen.CreateSubNoteResponse{
				Note: &gen.NoteModel{
					Id:            childId.String(),
					OwnerId:       userId.String(),
					Tags:          []string{},
					CreateTime:    time.Time{}.String(),
					UpdateTime:    time.Time{}.String(),
					Parent:        noteId.String(),
					Collaborators: []string{},
					Children:      []string{},
					Data:          "",
				},
			},
			wantErr: false,
			mocker: func(req *gen.CreateSubNoteRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().CreateSubNote(gomock.Any(), uuid.FromStringOrNil(req.UserId), req.NoteData, uuid.FromStringOrNil(req.ParentId)).Return(
					models.Note{
						Id:      childId,
						OwnerId: userId,
						Parent:  noteId,
					},
					nil)
			},
		},

		{
			name: "Test_CreateSubNote_Fail",
			requestBody: &gen.CreateSubNoteRequest{
				UserId:   userId.String(),
				NoteData: "",
				ParentId: noteId.String(),
			},
			want:    nil,
			wantErr: true,
			mocker: func(req *gen.CreateSubNoteRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().CreateSubNote(gomock.Any(), uuid.FromStringOrNil(req.UserId), req.NoteData, uuid.FromStringOrNil(req.ParentId)).Return(
					models.Note{}, errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			ctx := context.Background()

			h := NewGrpcNoteHandler(mockUsecase)
			tt.mocker(tt.requestBody, mockUsecase)

			got, err := h.CreateSubNote(ctx, tt.requestBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcNoteHandler.CreateSubNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGrpcNoteHandler_UpdateNote(t *testing.T) {
	noteId := uuid.NewV4()
	userId := uuid.NewV4()

	tests := []struct {
		name        string
		requestBody *gen.UpdateNoteRequest
		mocker      func(req *gen.UpdateNoteRequest, mock *mock_note.MockNoteUsecase)
		want        *gen.UpdateNoteResponse
		wantErr     bool
	}{
		{
			name: "Test_UpdateNote_Success",
			requestBody: &gen.UpdateNoteRequest{
				Data:   "",
				Id:     noteId.String(),
				UserId: userId.String(),
			},
			want: &gen.UpdateNoteResponse{
				Note: &gen.NoteModel{
					Id:            noteId.String(),
					OwnerId:       userId.String(),
					Tags:          []string{},
					CreateTime:    time.Time{}.String(),
					UpdateTime:    time.Time{}.String(),
					Parent:        uuid.UUID{}.String(),
					Collaborators: []string{},
					Children:      []string{},
					Data:          "",
				},
			},
			wantErr: false,
			mocker: func(req *gen.UpdateNoteRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().UpdateNote(gomock.Any(), uuid.FromStringOrNil(req.Id), uuid.FromStringOrNil(req.UserId), req.Data).Return(
					models.Note{
						Id:      noteId,
						OwnerId: userId,
					}, nil)
			},
		},

		{
			name: "Test_UpdateNote_Fail",
			requestBody: &gen.UpdateNoteRequest{
				Data:   "",
				Id:     noteId.String(),
				UserId: userId.String(),
			},
			want:    nil,
			wantErr: true,
			mocker: func(req *gen.UpdateNoteRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().UpdateNote(gomock.Any(), uuid.FromStringOrNil(req.Id), uuid.FromStringOrNil(req.UserId), req.Data).Return(
					models.Note{}, errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			ctx := context.Background()

			h := NewGrpcNoteHandler(mockUsecase)
			tt.mocker(tt.requestBody, mockUsecase)

			got, err := h.UpdateNote(ctx, tt.requestBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcNoteHandler.UpdateNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGrpcNoteHandler_CheckPermissions(t *testing.T) {
	noteId := uuid.NewV4()
	userId := uuid.NewV4()

	tests := []struct {
		name        string
		requestBody *gen.CheckPermissionsRequest
		mocker      func(req *gen.CheckPermissionsRequest, mock *mock_note.MockNoteUsecase)
		want        *gen.CheckPermissionsResponse
		wantErr     bool
	}{
		{
			name: "Test_CheckPermission_Success",
			requestBody: &gen.CheckPermissionsRequest{
				NoteId: noteId.String(),
				UserId: userId.String(),
			},
			want: &gen.CheckPermissionsResponse{
				Result: true,
			},
			wantErr: false,
			mocker: func(req *gen.CheckPermissionsRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().CheckPermissions(gomock.Any(), uuid.FromStringOrNil(req.NoteId), uuid.FromStringOrNil(req.UserId)).Return(true, nil)
			},
		},

		{
			name: "Test_CheckPermission_Success",
			requestBody: &gen.CheckPermissionsRequest{
				NoteId: noteId.String(),
				UserId: userId.String(),
			},
			want:    nil,
			wantErr: true,
			mocker: func(req *gen.CheckPermissionsRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().CheckPermissions(gomock.Any(), uuid.FromStringOrNil(req.NoteId), uuid.FromStringOrNil(req.UserId)).Return(false, errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			ctx := context.Background()

			h := NewGrpcNoteHandler(mockUsecase)
			tt.mocker(tt.requestBody, mockUsecase)

			got, err := h.CheckPermissions(ctx, tt.requestBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcNoteHandler.CheckPermission() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGrpcNoteHandler_RememberTag(t *testing.T) {
	userId := uuid.NewV4()

	tests := []struct {
		name        string
		requestBody *gen.AllTagRequest
		mocker      func(req *gen.AllTagRequest, mock *mock_note.MockNoteUsecase)
		wantErr     bool
	}{
		{
			name: "Test_RememberTag_Success",
			requestBody: &gen.AllTagRequest{
				TagName: "tag",
				UserId:  userId.String(),
			},

			wantErr: false,
			mocker: func(req *gen.AllTagRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().RememberTag(gomock.Any(), req.TagName, uuid.FromStringOrNil(req.UserId)).Return(nil)
			},
		},
		{
			name: "Test_RememberTag_Fail",
			requestBody: &gen.AllTagRequest{
				TagName: "tag",
				UserId:  userId.String(),
			},
			wantErr: true,
			mocker: func(req *gen.AllTagRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().RememberTag(gomock.Any(), req.TagName,
					uuid.FromStringOrNil(req.UserId)).Return(errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			ctx := context.Background()

			h := NewGrpcNoteHandler(mockUsecase)
			tt.mocker(tt.requestBody, mockUsecase)

			_, err := h.RememberTag(ctx, tt.requestBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcNoteHandler.RememberTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGrpcNoteHandler_UpdateTag(t *testing.T) {
	userId := uuid.NewV4()

	tests := []struct {
		name        string
		requestBody *gen.UpdateTagRequest
		mocker      func(req *gen.UpdateTagRequest, mock *mock_note.MockNoteUsecase)
		wantErr     bool
	}{
		{
			name: "Test_UpdateTag_Success",
			requestBody: &gen.UpdateTagRequest{
				OldTag: "old",
				NewTag: "new",
				UserId: userId.String(),
			},

			wantErr: false,
			mocker: func(req *gen.UpdateTagRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().UpdateTag(gomock.Any(), req.OldTag, req.NewTag, uuid.FromStringOrNil(req.UserId)).Return(nil)
			},
		},
		{
			name: "Test_UpdateTag_Fail",
			requestBody: &gen.UpdateTagRequest{
				OldTag: "old",
				NewTag: "new",
				UserId: userId.String(),
			},
			wantErr: true,
			mocker: func(req *gen.UpdateTagRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().UpdateTag(gomock.Any(), req.OldTag, req.NewTag,
					uuid.FromStringOrNil(req.UserId)).Return(errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			ctx := context.Background()

			h := NewGrpcNoteHandler(mockUsecase)
			tt.mocker(tt.requestBody, mockUsecase)

			_, err := h.UpdateTag(ctx, tt.requestBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcNoteHandler.UpdateTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGrpcNoteHandler_ForgetTag(t *testing.T) {
	userId := uuid.NewV4()

	tests := []struct {
		name        string
		requestBody *gen.AllTagRequest
		mocker      func(req *gen.AllTagRequest, mock *mock_note.MockNoteUsecase)
		wantErr     bool
	}{
		{
			name: "Test_ForgetTag_Success",
			requestBody: &gen.AllTagRequest{
				TagName: "tag",
				UserId:  userId.String(),
			},

			wantErr: false,
			mocker: func(req *gen.AllTagRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().ForgetTag(gomock.Any(), req.TagName, uuid.FromStringOrNil(req.UserId)).Return(nil)
			},
		},
		{
			name: "Test_ForgetTag_Fail",
			requestBody: &gen.AllTagRequest{
				TagName: "tag",
				UserId:  userId.String(),
			},
			wantErr: true,
			mocker: func(req *gen.AllTagRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().ForgetTag(gomock.Any(), req.TagName,
					uuid.FromStringOrNil(req.UserId)).Return(errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			ctx := context.Background()

			h := NewGrpcNoteHandler(mockUsecase)
			tt.mocker(tt.requestBody, mockUsecase)

			_, err := h.ForgetTag(ctx, tt.requestBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcNoteHandler.ForgetTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGrpcNoteHandler_SetIcon(t *testing.T) {
	userId := uuid.NewV4()
	noteId := uuid.NewV4()
	tests := []struct {
		name        string
		requestBody *gen.SetIconRequest
		mocker      func(req *gen.SetIconRequest, mock *mock_note.MockNoteUsecase)
		wantErr     bool
	}{
		{
			name: "Test_SetIcon_Success",
			requestBody: &gen.SetIconRequest{
				Icon:   "icon",
				NoteId: noteId.String(),
				UserId: userId.String(),
			},

			wantErr: false,
			mocker: func(req *gen.SetIconRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().SetIcon(gomock.Any(), uuid.FromStringOrNil(req.NoteId), req.Icon, uuid.FromStringOrNil(req.UserId)).Return(models.Note{}, nil)
			},
		},
		{
			name: "Test_SetIcon_Fail",
			requestBody: &gen.SetIconRequest{
				Icon:   "icon",
				NoteId: noteId.String(),
				UserId: userId.String(),
			},
			wantErr: true,
			mocker: func(req *gen.SetIconRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().SetIcon(gomock.Any(), uuid.FromStringOrNil(req.NoteId), req.Icon, uuid.FromStringOrNil(req.UserId)).Return(models.Note{}, errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			ctx := context.Background()

			h := NewGrpcNoteHandler(mockUsecase)
			tt.mocker(tt.requestBody, mockUsecase)

			_, err := h.SetIcon(ctx, tt.requestBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcNoteHandler.SetIcon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGrpcNoteHandler_SetHeader(t *testing.T) {
	userId := uuid.NewV4()
	noteId := uuid.NewV4()
	tests := []struct {
		name        string
		requestBody *gen.SetHeaderRequest
		mocker      func(req *gen.SetHeaderRequest, mock *mock_note.MockNoteUsecase)
		wantErr     bool
	}{
		{
			name: "Test_SetHeader_Success",
			requestBody: &gen.SetHeaderRequest{
				Header: "header",
				NoteId: noteId.String(),
				UserId: userId.String(),
			},

			wantErr: false,
			mocker: func(req *gen.SetHeaderRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().SetHeader(gomock.Any(), uuid.FromStringOrNil(req.NoteId), req.Header, uuid.FromStringOrNil(req.UserId)).Return(models.Note{}, nil)
			},
		},
		{
			name: "Test_SetHeader_Fail",
			requestBody: &gen.SetHeaderRequest{
				Header: "header",
				NoteId: noteId.String(),
				UserId: userId.String(),
			},
			wantErr: true,
			mocker: func(req *gen.SetHeaderRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().SetHeader(gomock.Any(), uuid.FromStringOrNil(req.NoteId), req.Header, uuid.FromStringOrNil(req.UserId)).Return(models.Note{}, errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			ctx := context.Background()

			h := NewGrpcNoteHandler(mockUsecase)
			tt.mocker(tt.requestBody, mockUsecase)

			_, err := h.SetHeader(ctx, tt.requestBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcNoteHandler.SetHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGrpcNoteHandler_AddFav(t *testing.T) {
	userId := uuid.NewV4()
	noteId := uuid.NewV4()
	tests := []struct {
		name        string
		requestBody *gen.ChangeFlagRequest
		mocker      func(req *gen.ChangeFlagRequest, mock *mock_note.MockNoteUsecase)
		wantErr     bool
	}{
		{
			name: "Test_AddFav_Success",
			requestBody: &gen.ChangeFlagRequest{
				NoteId: noteId.String(),
				UserId: userId.String(),
			},

			wantErr: false,
			mocker: func(req *gen.ChangeFlagRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().AddFav(gomock.Any(), uuid.FromStringOrNil(req.NoteId), uuid.FromStringOrNil(req.UserId)).Return(models.Note{}, nil)
			},
		},
		{
			name: "Test_AddFav_Fail",
			requestBody: &gen.ChangeFlagRequest{
				NoteId: noteId.String(),
				UserId: userId.String(),
			},
			wantErr: true,
			mocker: func(req *gen.ChangeFlagRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().AddFav(gomock.Any(), uuid.FromStringOrNil(req.NoteId), uuid.FromStringOrNil(req.UserId)).Return(models.Note{}, errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			ctx := context.Background()

			h := NewGrpcNoteHandler(mockUsecase)
			tt.mocker(tt.requestBody, mockUsecase)

			_, err := h.AddFav(ctx, tt.requestBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcNoteHandler.AddFav() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGrpcNoteHandler_DelFav(t *testing.T) {
	userId := uuid.NewV4()
	noteId := uuid.NewV4()
	tests := []struct {
		name        string
		requestBody *gen.ChangeFlagRequest
		mocker      func(req *gen.ChangeFlagRequest, mock *mock_note.MockNoteUsecase)
		wantErr     bool
	}{
		{
			name: "Test_DelFav_Success",
			requestBody: &gen.ChangeFlagRequest{
				NoteId: noteId.String(),
				UserId: userId.String(),
			},

			wantErr: false,
			mocker: func(req *gen.ChangeFlagRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().DelFav(gomock.Any(), uuid.FromStringOrNil(req.NoteId), uuid.FromStringOrNil(req.UserId)).Return(models.Note{}, nil)
			},
		},
		{
			name: "Test_DelFav_Fail",
			requestBody: &gen.ChangeFlagRequest{
				NoteId: noteId.String(),
				UserId: userId.String(),
			},
			wantErr: true,
			mocker: func(req *gen.ChangeFlagRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().DelFav(gomock.Any(), uuid.FromStringOrNil(req.NoteId), uuid.FromStringOrNil(req.UserId)).Return(models.Note{}, errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			ctx := context.Background()

			h := NewGrpcNoteHandler(mockUsecase)
			tt.mocker(tt.requestBody, mockUsecase)

			_, err := h.DelFav(ctx, tt.requestBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcNoteHandler.DelFav() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGrpcNoteHandler_SetPublic(t *testing.T) {
	userId := uuid.NewV4()
	noteId := uuid.NewV4()
	tests := []struct {
		name        string
		requestBody *gen.AccessModeRequest
		mocker      func(req *gen.AccessModeRequest, mock *mock_note.MockNoteUsecase)
		wantErr     bool
	}{
		{
			name: "Test_SetPublic_Success",
			requestBody: &gen.AccessModeRequest{
				NoteId: noteId.String(),
				UserId: userId.String(),
			},

			wantErr: false,
			mocker: func(req *gen.AccessModeRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().SetPublic(gomock.Any(), uuid.FromStringOrNil(req.NoteId), uuid.FromStringOrNil(req.UserId)).Return(models.Note{}, nil)
			},
		},
		{
			name: "Test_SetPublic_Fail",
			requestBody: &gen.AccessModeRequest{
				NoteId: noteId.String(),
				UserId: userId.String(),
			},
			wantErr: true,
			mocker: func(req *gen.AccessModeRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().SetPublic(gomock.Any(), uuid.FromStringOrNil(req.NoteId), uuid.FromStringOrNil(req.UserId)).Return(models.Note{}, errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			ctx := context.Background()

			h := NewGrpcNoteHandler(mockUsecase)
			tt.mocker(tt.requestBody, mockUsecase)

			_, err := h.SetPublic(ctx, tt.requestBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcNoteHandler.SetPublic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGrpcNoteHandler_SetPrivate(t *testing.T) {
	userId := uuid.NewV4()
	noteId := uuid.NewV4()
	tests := []struct {
		name        string
		requestBody *gen.AccessModeRequest
		mocker      func(req *gen.AccessModeRequest, mock *mock_note.MockNoteUsecase)
		wantErr     bool
	}{
		{
			name: "Test_SetPrivate_Success",
			requestBody: &gen.AccessModeRequest{
				NoteId: noteId.String(),
				UserId: userId.String(),
			},

			wantErr: false,
			mocker: func(req *gen.AccessModeRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().SetPrivate(gomock.Any(), uuid.FromStringOrNil(req.NoteId), uuid.FromStringOrNil(req.UserId)).Return(models.Note{}, nil)
			},
		},
		{
			name: "Test_SetPrivate_Fail",
			requestBody: &gen.AccessModeRequest{
				NoteId: noteId.String(),
				UserId: userId.String(),
			},
			wantErr: true,
			mocker: func(req *gen.AccessModeRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().SetPrivate(gomock.Any(), uuid.FromStringOrNil(req.NoteId), uuid.FromStringOrNil(req.UserId)).Return(models.Note{}, errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()
			ctx := context.Background()

			h := NewGrpcNoteHandler(mockUsecase)
			tt.mocker(tt.requestBody, mockUsecase)

			_, err := h.SetPrivate(ctx, tt.requestBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcNoteHandler.SetPrivate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGrpcNoteHandler_GetAttachList(t *testing.T) {
	noteId := uuid.NewV4()
	userId := uuid.NewV4()

	type args struct {
		ctx context.Context
		in  *generatedNote.GetAttachListRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *generatedNote.GetAttachListResponse
		mocker  func(req *gen.GetAttachListRequest, mock *mock_note.MockNoteUsecase)
		wantErr bool
	}{
		{
			name: "Test_GetAttachList_Success",
			args: args{
				ctx: context.Background(),
				in: &generatedNote.GetAttachListRequest{
					NoteId: noteId.String(),
					UserId: userId.String(),
				},
			},
			want: &generatedNote.GetAttachListResponse{
				Paths: []string{"1", "2"},
			},
			wantErr: false,
			mocker: func(req *generatedNote.GetAttachListRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().GetAttachList(gomock.Any(), uuid.FromStringOrNil(req.NoteId), uuid.FromStringOrNil(req.UserId)).Return([]string{"1", "2"}, nil)
			},
		},
		{
			name: "Test_GetAttachList_Fail",
			args: args{
				ctx: context.Background(),
				in: &generatedNote.GetAttachListRequest{
					NoteId: noteId.String(),
					UserId: userId.String(),
				},
			},
			want:    nil,
			wantErr: true,
			mocker: func(req *generatedNote.GetAttachListRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().GetAttachList(gomock.Any(), uuid.FromStringOrNil(req.NoteId), uuid.FromStringOrNil(req.UserId)).Return([]string{}, errors.New("err"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()

			h := NewGrpcNoteHandler(mockUsecase)
			tt.mocker(tt.args.in, mockUsecase)

			got, err := h.GetAttachList(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcNoteHandler.GetAttachList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GrpcNoteHandler.GetAttachList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrpcNoteHandler_GetSharedAttachList(t *testing.T) {
	noteId := uuid.NewV4()

	type args struct {
		ctx context.Context
		in  *generatedNote.GetSharedAttachListRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *generatedNote.GetAttachListResponse
		mocker  func(req *gen.GetSharedAttachListRequest, mock *mock_note.MockNoteUsecase)
		wantErr bool
	}{
		{
			name: "Test_GetSharedAttachList_Success",
			args: args{
				ctx: context.Background(),
				in: &generatedNote.GetSharedAttachListRequest{
					NoteId: noteId.String(),
				},
			},
			want: &generatedNote.GetAttachListResponse{
				Paths: []string{"1", "2"},
			},
			wantErr: false,
			mocker: func(req *generatedNote.GetSharedAttachListRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().GetSharedAttachList(gomock.Any(), uuid.FromStringOrNil(req.NoteId)).Return([]string{"1", "2"}, nil)
			},
		},
		{
			name: "Test_GetSharedAttachList_Fail",
			args: args{
				ctx: context.Background(),
				in: &generatedNote.GetSharedAttachListRequest{
					NoteId: noteId.String(),
				},
			},
			want:    nil,
			wantErr: true,
			mocker: func(req *generatedNote.GetSharedAttachListRequest, mock *mock_note.MockNoteUsecase) {
				mock.EXPECT().GetSharedAttachList(gomock.Any(), uuid.FromStringOrNil(req.NoteId)).Return([]string{}, errors.New("err"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
			defer ctrl.Finish()

			h := NewGrpcNoteHandler(mockUsecase)
			tt.mocker(tt.args.in, mockUsecase)

			got, err := h.GetSharedAttachList(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcNoteHandler.GetSharedAttachList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GrpcNoteHandler.GetSharedAttachList() = %v, want %v", got, tt.want)
			}
		})
	}
}
