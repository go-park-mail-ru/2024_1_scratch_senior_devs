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
				mock.EXPECT().AddCollaborator(gomock.Any(), uuid.FromStringOrNil(req.NoteId), uuid.FromStringOrNil(req.UserId), uuid.FromStringOrNil(req.GuestId)).Return(nil)
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
				mock.EXPECT().AddCollaborator(gomock.Any(), uuid.FromStringOrNil(req.NoteId), uuid.FromStringOrNil(req.UserId), uuid.FromStringOrNil(req.GuestId)).Return(errors.New("error"))
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
