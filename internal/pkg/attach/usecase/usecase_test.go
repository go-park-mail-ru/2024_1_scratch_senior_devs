package usecase

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	mock_attach "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/attach/mocks"
	mock_note "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/mocks"
	"github.com/golang/mock/gomock"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
)

var testLogger *slog.Logger

func init() {
	testLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func TestAttachUsecase_DeleteAttach(t *testing.T) {
	attachId := uuid.NewV4()
	userId := uuid.NewV4()
	noteId := uuid.NewV4()
	type args struct {
		ctx      context.Context
		attachID uuid.UUID
		userID   uuid.UUID
		noteID   uuid.UUID
	}
	tests := []struct {
		name        string
		repoMocker  func(context context.Context, repo *mock_attach.MockAttachRepo, noteRepo *mock_note.MockNoteRepo, data args)
		args        args
		expectedErr error
	}{
		{
			name: "TestDeleteAtatch_Success",
			repoMocker: func(ctx context.Context, repo *mock_attach.MockAttachRepo, noteRepo *mock_note.MockNoteRepo, data args) {
				repo.EXPECT().DeleteAttach(ctx, attachId).Return(nil).Times(1)
				noteRepo.EXPECT().ReadNote(ctx, noteId).Return(models.Note{
					Id:         data.noteID,
					Data:       []byte{},
					CreateTime: time.Time{},
					UpdateTime: time.Time{},
					OwnerId:    data.userID,
				}, nil).Times(1)
				repo.EXPECT().GetAttach(ctx, attachId).Return(models.Attach{
					Id:     data.attachID,
					NoteId: data.noteID,
					Path:   "",
				}, nil).Times(1)

			},
			args: args{
				ctx:      context.Background(),
				attachID: attachId,
				userID:   userId,
				noteID:   noteId,
			},
			expectedErr: nil,
		},
		{
			name: "TestDeleteAtatch_Fail_On_DeleteAttach",
			repoMocker: func(ctx context.Context, repo *mock_attach.MockAttachRepo, noteRepo *mock_note.MockNoteRepo, data args) {
				repo.EXPECT().DeleteAttach(ctx, attachId).Return(errors.New("delete error")).Times(1)
				noteRepo.EXPECT().ReadNote(ctx, noteId).Return(models.Note{
					Id:         data.noteID,
					Data:       []byte{},
					CreateTime: time.Time{},
					UpdateTime: time.Time{},
					OwnerId:    data.userID,
				}, nil).Times(1)
				repo.EXPECT().GetAttach(ctx, attachId).Return(models.Attach{
					Id:     data.attachID,
					NoteId: data.noteID,
					Path:   "",
				}, nil).Times(1)

			},
			args: args{
				ctx:      context.Background(),
				attachID: attachId,
				userID:   userId,
				noteID:   noteId,
			},
			expectedErr: errors.New("delete error"),
		},
		{
			name: "TestDeleteAtatch_Fail_On_ReadNote",
			repoMocker: func(ctx context.Context, repo *mock_attach.MockAttachRepo, noteRepo *mock_note.MockNoteRepo, data args) {
				noteRepo.EXPECT().ReadNote(ctx, noteId).Return(models.Note{}, errors.New("read note error")).Times(1)
				repo.EXPECT().GetAttach(ctx, attachId).Return(models.Attach{
					Id:     data.attachID,
					NoteId: data.noteID,
					Path:   "",
				}, nil).Times(1)

			},
			args: args{
				ctx:      context.Background(),
				attachID: attachId,
				userID:   userId,
				noteID:   noteId,
			},
			expectedErr: errors.New("read note error"),
		},
		{
			name: "TestDeleteAtatch_Fail_On_GetAttach",
			repoMocker: func(ctx context.Context, repo *mock_attach.MockAttachRepo, noteRepo *mock_note.MockNoteRepo, data args) {
				repo.EXPECT().GetAttach(ctx, attachId).Return(models.Attach{}, errors.New("get attach error")).Times(1)

			},
			args: args{
				ctx:      context.Background(),
				attachID: attachId,
				userID:   userId,
				noteID:   noteId,
			},
			expectedErr: errors.New("get attach error"),
		},
		{
			name: "TestDeleteAtatch_Fail_NotFound",
			repoMocker: func(ctx context.Context, repo *mock_attach.MockAttachRepo, noteRepo *mock_note.MockNoteRepo, data args) {
				noteRepo.EXPECT().ReadNote(ctx, noteId).Return(models.Note{
					Id:         data.noteID,
					Data:       []byte{},
					CreateTime: time.Time{},
					UpdateTime: time.Time{},
					OwnerId:    uuid.NewV4(),
				}, nil).Times(1)
				repo.EXPECT().GetAttach(ctx, attachId).Return(models.Attach{
					Id:     data.attachID,
					NoteId: data.noteID,
					Path:   "",
				}, nil).Times(1)

			},
			args: args{
				ctx:      context.Background(),
				attachID: attachId,
				userID:   userId,
				noteID:   noteId,
			},
			expectedErr: errors.New("note not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			noteRepo := mock_note.NewMockNoteRepo(ctl)
			repo := mock_attach.NewMockAttachRepo(ctl)
			uc := CreateAttachUsecase(repo, noteRepo, testLogger)

			tt.repoMocker(context.Background(), repo, noteRepo, tt.args)

			err := uc.DeleteAttach(tt.args.ctx, tt.args.attachID, tt.args.userID)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestAttachUsecase_GetAttach(t *testing.T) {
	attachId := uuid.NewV4()
	userId := uuid.NewV4()
	noteId := uuid.NewV4()
	type args struct {
		ctx      context.Context
		attachID uuid.UUID
		userID   uuid.UUID
	}
	tests := []struct {
		name        string
		repoMocker  func(context context.Context, repo *mock_attach.MockAttachRepo, noteRepo *mock_note.MockNoteRepo)
		args        args
		want        models.Attach
		expectedErr error
	}{
		{
			name: "Test_GetAttach_Success",
			repoMocker: func(ctx context.Context, repo *mock_attach.MockAttachRepo, noteRepo *mock_note.MockNoteRepo) {
				noteRepo.EXPECT().ReadNote(ctx, noteId).Return(models.Note{
					Id:         noteId,
					Data:       []byte{},
					CreateTime: time.Time{},
					UpdateTime: time.Time{},
					OwnerId:    userId,
				}, nil).Times(1)
				repo.EXPECT().GetAttach(ctx, attachId).Return(models.Attach{
					Id:     attachId,
					NoteId: noteId,
					Path:   "",
				}, nil).Times(1)
			},
			args: args{
				attachID: attachId,
				userID:   userId,
				ctx:      context.Background(),
			},
			want: models.Attach{
				Id:     attachId,
				Path:   "",
				NoteId: noteId,
			},
			expectedErr: nil,
		},
		{
			name: "Test_GetAttach_Fail_On_ReadNote",
			repoMocker: func(ctx context.Context, repo *mock_attach.MockAttachRepo, noteRepo *mock_note.MockNoteRepo) {
				noteRepo.EXPECT().ReadNote(ctx, noteId).Return(models.Note{}, errors.New("read note error")).Times(1)
				repo.EXPECT().GetAttach(ctx, attachId).Return(models.Attach{
					Id:     attachId,
					NoteId: noteId,
					Path:   "",
				}, nil).Times(1)
			},
			args: args{
				attachID: attachId,
				userID:   userId,
				ctx:      context.Background(),
			},
			want:        models.Attach{},
			expectedErr: errors.New("read note error"),
		},
		{
			name: "Test_GetAttach_Fail_On_GetAttach",
			repoMocker: func(ctx context.Context, repo *mock_attach.MockAttachRepo, noteRepo *mock_note.MockNoteRepo) {

				repo.EXPECT().GetAttach(ctx, attachId).Return(models.Attach{}, errors.New("get attach error")).Times(1)
			},
			args: args{
				attachID: attachId,
				userID:   userId,
				ctx:      context.Background(),
			},
			want:        models.Attach{},
			expectedErr: errors.New("get attach error"),
		},
		{
			name: "Test_GetAttach_NotFound",
			repoMocker: func(ctx context.Context, repo *mock_attach.MockAttachRepo, noteRepo *mock_note.MockNoteRepo) {
				noteRepo.EXPECT().ReadNote(ctx, noteId).Return(models.Note{
					Id:         noteId,
					Data:       []byte{},
					CreateTime: time.Time{},
					UpdateTime: time.Time{},
					OwnerId:    uuid.NewV4(),
				}, nil).Times(1)
				repo.EXPECT().GetAttach(ctx, attachId).Return(models.Attach{
					Id:     attachId,
					NoteId: noteId,
					Path:   "",
				}, nil).Times(1)
			},
			args: args{
				attachID: attachId,
				userID:   userId,
				ctx:      context.Background(),
			},
			want:        models.Attach{},
			expectedErr: errors.New("note not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			noteRepo := mock_note.NewMockNoteRepo(ctl)
			repo := mock_attach.NewMockAttachRepo(ctl)
			uc := CreateAttachUsecase(repo, noteRepo, testLogger)

			tt.repoMocker(context.Background(), repo, noteRepo)

			attachData, err := uc.GetAttach(tt.args.ctx, tt.args.attachID, tt.args.userID)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.want, attachData)
		})
	}
}
