package usecase

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	mock_note "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/mocks"
	"github.com/golang/mock/gomock"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNoteUsecase_GetAllNotes(t *testing.T) {
	elasticConfig := config.ElasticConfig{
		ElasticIndexName:            "notes",
		ElasticSearchValueMinLength: 2,
	}

	type args struct {
		userId uuid.UUID
		count  int64
		offset int64
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, uId uuid.UUID, count int64, offset int64)
		args       args
		want       []models.Note
		wantErr    bool
	}{
		{
			name: "TestSuccess",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, uId uuid.UUID, count int64, offset int64) {
				mockResp := []models.Note{ //мок ответа от уровня репозитория
					{
						Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
						OwnerId:    uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
						UpdateTime: time.Time{},
						CreateTime: time.Time{},
						Data:       nil,
						Parent:     uuid.UUID{},
						Children:   []uuid.UUID{},
					},
					{
						Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b12e-b41604c56f95"),
						OwnerId:    uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
						UpdateTime: time.Time{},
						CreateTime: time.Time{},
						Data:       nil,
						Parent:     uuid.UUID{},
						Children:   []uuid.UUID{},
					},
				}

				baseRepo.EXPECT().ReadAllNotes(ctx, uId, int64(count), int64(offset)).Return(mockResp, nil).Times(1)
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
					Parent:     uuid.UUID{},
					Children:   []uuid.UUID{},
				},
				{
					Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b12e-b41604c56f95"),
					OwnerId:    uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
					UpdateTime: time.Time{},
					CreateTime: time.Time{},
					Data:       nil,
					Parent:     uuid.UUID{},
					Children:   []uuid.UUID{},
				},
			},
			wantErr: false,
		},
		{
			name: "TestFail",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, uId uuid.UUID, count int64, offset int64) {
				mockResp := []models.Note{ //мок ответа от уровня репозитория

				}

				baseRepo.EXPECT().ReadAllNotes(ctx, uId, int64(count), int64(offset)).Return(mockResp, errors.New("repo error")).Times(1)
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
			baseRepo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			Usecase := CreateNoteUsecase(baseRepo, searchRepo, elasticConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), baseRepo, searchRepo, tt.args.userId, tt.args.count, tt.args.offset)
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
	elasticConfig := config.ElasticConfig{
		ElasticIndexName:            "notes",
		ElasticSearchValueMinLength: 2,
	}

	type args struct {
		ctx    context.Context
		noteId uuid.UUID
		userId uuid.UUID
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, repo *mock_note.MockNoteBaseRepo, nId uuid.UUID)
		args       args
		want       models.Note
		wantErr    bool
	}{
		{
			name: "TestSuccess",
			repoMocker: func(ctx context.Context, repo *mock_note.MockNoteBaseRepo, nId uuid.UUID) {
				mockResp := models.Note{ //мок ответа от уровня репозитория
					Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
					OwnerId:    uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
					UpdateTime: time.Time{},
					CreateTime: time.Time{},
					Data:       nil,
					Parent:     uuid.UUID{},
					Children:   []uuid.UUID{},
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
				Parent:     uuid.UUID{},
				Children:   []uuid.UUID{},
			},

			wantErr: false,
		},
		{
			name: "TestFail",
			repoMocker: func(ctx context.Context, repo *mock_note.MockNoteBaseRepo, nId uuid.UUID) {
				mockResp := models.Note{ //мок ответа от уровня репозитория

				}

				repo.EXPECT().ReadNote(ctx, nId).Return(mockResp, errors.New("error")).Times(1)
				repo.EXPECT().CheckCollaborator(ctx, nId, uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95")).Return(false, nil)
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
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, &sync.WaitGroup{})

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
	elasticConfig := config.ElasticConfig{
		ElasticIndexName:            "notes",
		ElasticSearchValueMinLength: 2,
	}

	type args struct {
		ctx      context.Context
		userId   uuid.UUID
		noteData []byte
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo)
		args       args
		wantErr    bool
		want       models.Note
	}{
		{
			name: "TestSuccess",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().CreateNote(ctx, gomock.Any()).Return(nil).Times(1)
				searchRepo.EXPECT().CreateNote(ctx, gomock.Any()).Return(nil).Times(1)
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
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo)

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

func TestNoteUsecase_UpdateNote(t *testing.T) {
	elasticConfig := config.ElasticConfig{
		ElasticIndexName:            "notes",
		ElasticSearchValueMinLength: 2,
	}

	id := uuid.NewV4()

	type args struct {
		ctx      context.Context
		userId   uuid.UUID
		noteData []byte
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo)
		args       args
		wantErr    bool
		want       models.Note
	}{
		{
			name: "TestSuccess",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().UpdateNote(ctx, gomock.Any()).Return(nil).Times(1)
				baseRepo.EXPECT().ReadNote(ctx, gomock.Any()).Return(models.Note{}, nil).Times(1)
				searchRepo.EXPECT().UpdateNote(ctx, gomock.Any()).Return(nil).Times(1)
				searchRepo.EXPECT().ReadNote(ctx, gomock.Any()).Return(models.ElasticNote{}, nil).Times(1)
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
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo)

			got, err := uc.UpdateNote(tt.args.ctx, id, tt.args.userId, tt.args.noteData)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.UpdateNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("NoteUsecase.UpdateNote() returned %v, want %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestNoteUsecase_DeleteNote(t *testing.T) {
	elasticConfig := config.ElasticConfig{
		ElasticIndexName:            "notes",
		ElasticSearchValueMinLength: 2,
	}

	id := uuid.NewV4()

	type args struct {
		ctx      context.Context
		userId   uuid.UUID
		noteData []byte
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo)
		args       args
		wantErr    bool
		want       models.Note
	}{
		{
			name: "TestSuccess",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().DeleteNote(ctx, gomock.Any()).Return(nil).Times(1)
				baseRepo.EXPECT().ReadNote(ctx, gomock.Any()).Return(models.Note{}, nil).Times(1)
				searchRepo.EXPECT().DeleteNote(ctx, gomock.Any()).Return(nil).Times(1)
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
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo)

			err := uc.DeleteNote(tt.args.ctx, id, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.DeleteNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
