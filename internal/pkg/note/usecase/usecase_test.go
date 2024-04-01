package usecase

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	mock_note "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/mocks"
	"github.com/golang/mock/gomock"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
)

var testLogger *slog.Logger

func init() {
	testLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func TestNoteUsecase_GetAllNotes(t *testing.T) {

	type args struct {
		userId uuid.UUID
		count  int64
		offset int64
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, repo *mock_note.MockNoteRepo, uId uuid.UUID, count int64, offset int64)
		args       args
		want       []models.Note
		wantErr    bool
	}{
		{
			name: "TestSuccess",
			repoMocker: func(ctx context.Context, repo *mock_note.MockNoteRepo, uId uuid.UUID, count int64, offset int64) {
				mockResp := []models.Note{ //мок ответа от уровня репозитория
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
				}

				repo.EXPECT().ReadAllNotes(ctx, uId, int64(count), int64(offset), "").Return(mockResp, nil).Times(1)
			},
			args: args{

				uuid.NewV4(),
				10,
				0,
			},
			want: []models.Note{
				{
					Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
					OwnerId:    uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"), //потом задать из args
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
			wantErr: false,
		},
		{
			name: "TestFail",
			repoMocker: func(ctx context.Context, repo *mock_note.MockNoteRepo, uId uuid.UUID, count int64, offset int64) {
				mockResp := []models.Note{ //мок ответа от уровня репозитория

				}

				repo.EXPECT().ReadAllNotes(ctx, uId, int64(count), int64(offset), "").Return(mockResp, errors.New("repo error")).Times(1)
			},
			args: args{

				uuid.NewV4(),
				10,
				0,
			},
			want:    make([]models.Note, 0),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteRepo(ctl)
			Usecase := CreateNoteUsecase(repo, testLogger)

			tt.repoMocker(context.Background(), repo, tt.args.userId, tt.args.count, tt.args.offset)
			got, err := Usecase.GetAllNotes(context.Background(), tt.args.userId, tt.args.count, tt.args.offset, "")

			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.GetAllNotes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("NoteUsecase.GetAllNotes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoteUsecase_GetNote(t *testing.T) {

	type args struct {
		ctx    context.Context
		noteId uuid.UUID
		userId uuid.UUID
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, repo *mock_note.MockNoteRepo, nId uuid.UUID)
		args       args
		want       models.Note
		wantErr    bool
	}{
		{
			name: "TestSuccess",
			repoMocker: func(ctx context.Context, repo *mock_note.MockNoteRepo, nId uuid.UUID) {
				mockResp := models.Note{ //мок ответа от уровня репозитория
					Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
					OwnerId:    uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
					UpdateTime: time.Time{},
					CreateTime: time.Time{},
					Data:       nil,
				}

				repo.EXPECT().ReadNote(ctx, nId).Return(mockResp, nil).Times(1)
			},
			args: args{
				context.Background(),
				uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
				uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
			},
			want: models.Note{
				Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
				OwnerId:    uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
				UpdateTime: time.Time{},
				CreateTime: time.Time{},
				Data:       nil,
			},

			wantErr: false,
		},
		{
			name: "TestFail",
			repoMocker: func(ctx context.Context, repo *mock_note.MockNoteRepo, nId uuid.UUID) {
				mockResp := models.Note{ //мок ответа от уровня репозитория

				}

				repo.EXPECT().ReadNote(ctx, nId).Return(mockResp, errors.New("error")).Times(1)
			},
			args: args{
				context.Background(),
				uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
				uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
			},
			want: models.Note{},

			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteRepo(ctl)
			uc := CreateNoteUsecase(repo, testLogger)

			tt.repoMocker(context.Background(), repo, tt.args.noteId)

			got, err := uc.GetNote(tt.args.ctx, tt.args.noteId, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.GetNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NoteUsecase.GetNote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoteUsecase_CreateNote(t *testing.T) {

	type args struct {
		ctx      context.Context
		userId   uuid.UUID
		noteData []byte
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, repo *mock_note.MockNoteRepo)
		args       args
		wantErr    bool
		want       models.Note
	}{
		{
			name: "TestSuccess",
			repoMocker: func(ctx context.Context, repo *mock_note.MockNoteRepo) {

				repo.EXPECT().CreateNote(ctx, gomock.Any()).Return(nil).Times(1)

			},
			args: args{
				ctx:      context.Background(),
				userId:   uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
				noteData: []byte("{\"title\":\"title\"}"),
			},
			wantErr: false,
			want:    models.Note{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteRepo(ctl)
			uc := CreateNoteUsecase(repo, testLogger)

			tt.repoMocker(context.Background(), repo)

			got, err := uc.CreateNote(tt.args.ctx, tt.args.userId, tt.args.noteData)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.CreateNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("NoteUsecase.CreateNote() returned %v, want %v", err, tt.wantErr)
				return
			}

		})
	}
}
