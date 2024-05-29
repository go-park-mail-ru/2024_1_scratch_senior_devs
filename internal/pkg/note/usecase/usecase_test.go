package usecase

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
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

	constraintsConfig := config.ConstraintsConfig{
		MaxDepth:         3,
		MaxCollaborators: 10,
		MaxTags:          10,
		MaxSubnotes:      10,
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
		want       []models.NoteResponse
		wantErr    bool
	}{
		{
			name: "TestSuccess",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, uId uuid.UUID, count int64, offset int64) {
				mockResp := []models.NoteResponse{ //мок ответа от уровня репозитория
					{
						Note: models.Note{
							Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
							OwnerId:    uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
							UpdateTime: time.Time{},
							CreateTime: time.Time{},
							Data:       "",
							Parent:     uuid.UUID{},
							Children:   []uuid.UUID{},
						},
					},
					{
						Note: models.Note{
							Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b12e-b41604c56f95"),
							OwnerId:    uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
							UpdateTime: time.Time{},
							CreateTime: time.Time{},
							Data:       "",
							Parent:     uuid.UUID{},
							Children:   []uuid.UUID{},
						},
					},
				}

				baseRepo.EXPECT().ReadAllNotes(ctx, uId, int64(count), int64(offset), []string{"first"}).Return(mockResp, nil).Times(1)
				baseRepo.EXPECT().GetOwnerInfo(gomock.Any(), gomock.Any()).Return(models.OwnerInfo{}, nil).Times(2)
			},
			args: args{

				uuid.NewV4(),
				10,
				0,
			},
			want: []models.NoteResponse{
				{
					Note: models.Note{
						Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
						OwnerId:    uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"), //потом задать из args
						UpdateTime: time.Time{},
						CreateTime: time.Time{},
						Data:       "",
						Parent:     uuid.UUID{},
						Children:   []uuid.UUID{},
					},
				},
				{
					Note: models.Note{
						Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b12e-b41604c56f95"),
						OwnerId:    uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
						UpdateTime: time.Time{},
						CreateTime: time.Time{},
						Data:       "",
						Parent:     uuid.UUID{},
						Children:   []uuid.UUID{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "TestFail",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, uId uuid.UUID, count int64, offset int64) {
				mockResp := []models.NoteResponse{ //мок ответа от уровня репозитория

				}

				baseRepo.EXPECT().ReadAllNotes(ctx, uId, int64(count), int64(offset), []string{"first"}).Return(mockResp, errors.New("repo error")).Times(1)
			},
			args: args{

				uuid.NewV4(),
				10,
				0,
			},
			want:    make([]models.NoteResponse, 0),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctl := gomock.NewController(t)
			defer ctl.Finish()
			baseRepo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			Usecase := CreateNoteUsecase(baseRepo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), baseRepo, searchRepo, tt.args.userId, tt.args.count, tt.args.offset)
			got, err := Usecase.GetAllNotes(context.Background(), tt.args.userId, tt.args.count, tt.args.offset, "", []string{"first"})

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

	constraintsConfig := config.ConstraintsConfig{
		MaxDepth:         3,
		MaxCollaborators: 10,
		MaxTags:          10,
		MaxSubnotes:      10,
	}

	type args struct {
		ctx    context.Context
		noteId uuid.UUID
		userId uuid.UUID
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, repo *mock_note.MockNoteBaseRepo, args args)
		args       args
		want       models.NoteResponse
		wantErr    bool
	}{
		{
			name: "TestSuccess",
			repoMocker: func(ctx context.Context, repo *mock_note.MockNoteBaseRepo, args args) {
				mockResp := models.NoteResponse{
					Note: models.Note{ //мок ответа от уровня репозитория
						Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
						OwnerId:    uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
						UpdateTime: time.Time{},
						CreateTime: time.Time{},
						Data:       "",
						Parent:     uuid.UUID{},
						Children:   []uuid.UUID{},
					},
				}

				repo.EXPECT().ReadNote(ctx, args.noteId, args.userId).Return(mockResp, nil).Times(1)
				repo.EXPECT().GetOwnerInfo(gomock.Any(), gomock.Any()).Return(models.OwnerInfo{}, nil).Times(1)
			},
			args: args{
				context.Background(),
				uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
				uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
			},
			want: models.NoteResponse{
				Note: models.Note{
					Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
					OwnerId:    uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
					UpdateTime: time.Time{},
					CreateTime: time.Time{},
					Data:       "",
					Parent:     uuid.UUID{},
					Children:   []uuid.UUID{},
				},
			},

			wantErr: false,
		},
		{
			name: "TestFail",
			repoMocker: func(ctx context.Context, repo *mock_note.MockNoteBaseRepo, args args) {
				mockResp := models.NoteResponse{ //мок ответа от уровня репозитория

				}

				repo.EXPECT().ReadNote(ctx, args.noteId, args.userId).Return(mockResp, errors.New("error")).Times(1)
			},
			args: args{
				context.Background(),
				uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
				uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
			},
			want: models.NoteResponse{},

			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, tt.args)

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

func TestNoteUsecase_GetPublicNote(t *testing.T) {
	elasticConfig := config.ElasticConfig{
		ElasticIndexName:            "notes",
		ElasticSearchValueMinLength: 2,
	}

	constraintsConfig := config.ConstraintsConfig{
		MaxDepth:         3,
		MaxCollaborators: 10,
		MaxTags:          10,
		MaxSubnotes:      10,
	}

	type args struct {
		ctx    context.Context
		noteId uuid.UUID
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, repo *mock_note.MockNoteBaseRepo, args args)
		args       args
		want       models.NoteResponse
		wantErr    bool
	}{
		{
			name: "TestSuccess",
			repoMocker: func(ctx context.Context, repo *mock_note.MockNoteBaseRepo, args args) {
				mockResp := models.NoteResponse{
					Note: models.Note{ //мок ответа от уровня репозитория
						Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
						OwnerId:    uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
						UpdateTime: time.Time{},
						CreateTime: time.Time{},
						Data:       "",
						Parent:     uuid.UUID{},
						Children:   []uuid.UUID{},
						Public:     true,
					},
				}

				repo.EXPECT().ReadPublicNote(ctx, args.noteId).Return(mockResp, nil).Times(1)
				repo.EXPECT().GetOwnerInfo(gomock.Any(), gomock.Any()).Return(models.OwnerInfo{}, nil).Times(1)
			},
			args: args{
				context.Background(),
				uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
			},
			want: models.NoteResponse{
				Note: models.Note{
					Id:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
					OwnerId:    uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
					UpdateTime: time.Time{},
					CreateTime: time.Time{},
					Data:       "",
					Parent:     uuid.UUID{},
					Children:   []uuid.UUID{},
					Public:     true,
				},
			},

			wantErr: false,
		},
		{
			name: "TestFail",
			repoMocker: func(ctx context.Context, repo *mock_note.MockNoteBaseRepo, args args) {
				mockResp := models.NoteResponse{ //мок ответа от уровня репозитория

				}

				repo.EXPECT().ReadPublicNote(ctx, args.noteId).Return(mockResp, errors.New("error")).Times(1)
			},
			args: args{
				context.Background(),
				uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
			},
			want: models.NoteResponse{},

			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, tt.args)

			got, err := uc.GetPublicNote(tt.args.ctx, tt.args.noteId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.GetPublicNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NoteUsecase.GetPublicNote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoteUsecase_CreateNote(t *testing.T) {
	elasticConfig := config.ElasticConfig{
		ElasticIndexName:            "notes",
		ElasticSearchValueMinLength: 2,
	}

	constraintsConfig := config.ConstraintsConfig{
		MaxDepth:         3,
		MaxCollaborators: 10,
		MaxTags:          10,
		MaxSubnotes:      10,
	}

	type args struct {
		ctx      context.Context
		userId   uuid.UUID
		noteData string
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
				noteData: "{\"title\":\"title\"}",
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
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

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

	constraintsConfig := config.ConstraintsConfig{
		MaxDepth:         3,
		MaxCollaborators: 10,
		MaxTags:          10,
		MaxSubnotes:      10,
	}

	id := uuid.NewV4()

	type args struct {
		ctx      context.Context
		userId   uuid.UUID
		noteData string
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
				baseRepo.EXPECT().ReadNote(ctx, gomock.Any(), gomock.Any()).Return(models.NoteResponse{}, nil).Times(1)
				searchRepo.EXPECT().UpdateNote(ctx, gomock.Any()).Return(nil).Times(1)
			},
			args: args{
				ctx:      context.Background(),
				userId:   uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
				noteData: "{\"title\":\"title\"}",
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
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

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

	constraintsConfig := config.ConstraintsConfig{
		MaxDepth:         3,
		MaxCollaborators: 10,
		MaxTags:          10,
		MaxSubnotes:      10,
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
				baseRepo.EXPECT().ReadNote(ctx, gomock.Any(), gomock.Any()).Return(models.NoteResponse{}, nil).Times(1)
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
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo)

			err := uc.DeleteNote(tt.args.ctx, id, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.DeleteNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNoteUsecase_GetTags(t *testing.T) {

	elasticConfig := config.ElasticConfig{
		ElasticIndexName:            "notes",
		ElasticSearchValueMinLength: 2,
	}

	constraintsConfig := config.ConstraintsConfig{
		MaxDepth:         3,
		MaxCollaborators: 10,
		MaxTags:          10,
		MaxSubnotes:      10,
	}

	userId := uuid.NewV4()

	type args struct {
		ctx    context.Context
		userId uuid.UUID
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo)
		args       args
		wantErr    bool
		want       []string
	}{
		{
			name: "TestSuccess",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().GetTags(ctx, userId).Return([]string{"tag1", "tag2"}, nil)
			},
			want: []string{"tag1", "tag2"},
			args: args{
				ctx:    context.Background(),
				userId: userId,
			},
		},
		{
			name: "TestFail",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().GetTags(ctx, userId).Return([]string{}, errors.New("error"))
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
			},
			wantErr: true,
			want:    []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo)

			got, err := uc.GetTags(tt.args.ctx, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.CheckPermissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNoteUsecase_DeleteTag(t *testing.T) {
	elasticConfig := config.ElasticConfig{
		ElasticIndexName:            "notes",
		ElasticSearchValueMinLength: 2,
	}

	constraintsConfig := config.ConstraintsConfig{
		MaxDepth:         3,
		MaxCollaborators: 10,
		MaxTags:          10,
		MaxSubnotes:      10,
	}

	noteId := uuid.NewV4()
	userId := uuid.NewV4()

	type args struct {
		ctx     context.Context
		userId  uuid.UUID
		noteId  uuid.UUID
		tagName string
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo)
		args       args
		wantErr    bool
		want       models.Note
	}{
		{
			name: "Test_DeleteTag_Success",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().DeleteTag(ctx, "tag1", noteId).Return(nil).Times(1)
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: userId,
						Id:      noteId,
						Tags:    []string{"tag1", "tag2"},
					},
				}, nil).Times(1)
				searchRepo.EXPECT().DeleteTag(ctx, "tag1", noteId).Return(nil)

			},
			args: args{
				ctx:     context.Background(),
				userId:  userId,
				noteId:  noteId,
				tagName: "tag1",
			},
			wantErr: false,
			want: models.Note{
				OwnerId: userId,
				Id:      noteId,
				Tags:    []string{"tag2"},
			},
		},
		{
			name: "Test_DeleteTag_FailOnDelete",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().DeleteTag(ctx, "tag1", noteId).Return(errors.New("error")).Times(1)
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: userId,
						Id:      noteId,
						Tags:    []string{"tag1", "tag2"},
					},
				}, nil).Times(1)

			},
			args: args{
				ctx:     context.Background(),
				userId:  userId,
				noteId:  noteId,
				tagName: "tag1",
			},
			wantErr: true,
			want:    models.Note{},
		},
		{
			name: "Test_DeleteTag_FailOnRead",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{}, errors.New("error")).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    models.Note{},
		},
		{
			name: "Test_DeleteTag_NotOwner",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: uuid.NewV4(),
						Id:      noteId,
					},
				}, nil).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    models.Note{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo)

			got, err := uc.DeleteTag(tt.args.ctx, tt.args.tagName, tt.args.noteId, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.DeleteTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNoteUsecase_AddTag(t *testing.T) {
	elasticConfig := config.ElasticConfig{
		ElasticIndexName:            "notes",
		ElasticSearchValueMinLength: 2,
	}

	constraintsConfig := config.ConstraintsConfig{
		MaxDepth:         3,
		MaxCollaborators: 10,
		MaxTags:          4,
		MaxSubnotes:      10,
	}

	noteId := uuid.NewV4()
	userId := uuid.NewV4()

	type args struct {
		ctx     context.Context
		userId  uuid.UUID
		noteId  uuid.UUID
		tagName string
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo)
		args       args
		wantErr    bool
		want       models.Note
	}{
		{
			name: "Test_AddTag_Success",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().AddTag(ctx, "tag1", noteId).Return(nil).Times(1)
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: userId,
						Id:      noteId,
						Tags:    []string{"tag2"},
					},
				}, nil).Times(1)
				searchRepo.EXPECT().AddTag(ctx, "tag1", noteId).Return(nil)
				baseRepo.EXPECT().RememberTag(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			args: args{
				ctx:     context.Background(),
				userId:  userId,
				noteId:  noteId,
				tagName: "tag1",
			},
			wantErr: false,
			want: models.Note{
				OwnerId: userId,
				Id:      noteId,
				Tags:    []string{"tag2", "tag1"},
			},
		},
		{
			name: "Test_AddTag_FailOnAdd",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().AddTag(ctx, "tag1", noteId).Return(errors.New("error")).Times(1)
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: userId,
						Id:      noteId,
						Tags:    []string{"tag2"},
					},
				}, nil).Times(1)

			},
			args: args{
				ctx:     context.Background(),
				userId:  userId,
				noteId:  noteId,
				tagName: "tag1",
			},
			wantErr: true,
			want:    models.Note{},
		},
		{
			name: "Test_AddTag_FailTooMany",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: userId,
						Id:      noteId,
						Tags:    []string{"1", "2", "3", "4"},
					},
				}, nil).Times(1)

			},
			args: args{
				ctx:     context.Background(),
				userId:  userId,
				noteId:  noteId,
				tagName: "tag1",
			},
			wantErr: true,
			want:    models.Note{},
		},
		{
			name: "Test_AddTag_FailOnRead",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{}, errors.New("error")).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    models.Note{},
		},
		{
			name: "Test_AddTag_NotOwner",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: uuid.NewV4(),
						Id:      noteId,
					},
				}, nil).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    models.Note{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo)

			got, err := uc.AddTag(tt.args.ctx, tt.args.tagName, tt.args.noteId, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.AddTag error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNoteUsecase_addCollaboratorRecursive(t *testing.T) {
	elasticConfig := config.ElasticConfig{
		ElasticIndexName:            "notes",
		ElasticSearchValueMinLength: 2,
	}

	constraintsConfig := config.ConstraintsConfig{
		MaxDepth:         3,
		MaxCollaborators: 10,
		MaxTags:          4,
		MaxSubnotes:      10,
	}

	noteId := uuid.NewV4()
	guestId := uuid.NewV4()

	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo)
		wantErr    bool
	}{
		{
			name:    "Test_addCollaboratorReqursive_ReadError",
			wantErr: true,
			repoMocker: func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(gomock.Any(), noteId, gomock.Any()).Return(models.NoteResponse{}, errors.New("error"))
			},
		},
		{
			name:    "Test_addCollaboratorReqursive_AddError",
			wantErr: true,
			repoMocker: func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(gomock.Any(), noteId, gomock.Any()).Return(models.NoteResponse{
					Note: models.Note{
						Id: noteId,
					},
				}, nil)
				baseRepo.EXPECT().AddCollaborator(gomock.Any(), noteId, guestId).Return("", errors.New("error"))

			},
		},
		{
			name:    "Test_addCollaboratorReqursive_NoChildren",
			wantErr: false,
			repoMocker: func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(gomock.Any(), noteId, gomock.Any()).Return(models.NoteResponse{
					Note: models.Note{
						Id: noteId,
					},
				}, nil)
				baseRepo.EXPECT().AddCollaborator(gomock.Any(), noteId, guestId).Return("", nil)

			},
		},
		{
			name:    "Test_addCollaboratorReqursive_Success",
			wantErr: false,
			repoMocker: func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(gomock.Any(), noteId, gomock.Any()).Return(models.NoteResponse{
					Note: models.Note{
						Id:       noteId,
						Children: []uuid.UUID{noteId},
					},
				}, nil)
				baseRepo.EXPECT().ReadNote(gomock.Any(), noteId, gomock.Any()).Return(models.NoteResponse{
					Note: models.Note{
						Id: noteId,
					},
				}, nil)
				baseRepo.EXPECT().AddCollaborator(gomock.Any(), noteId, guestId).Return("", nil).Times(2)

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo)

			err := uc.addCollaboratorRecursive(context.Background(), noteId, guestId, uuid.NewV4())
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.addCollaboratorReqursive error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestNoteUsecase_AddCollaborator(t *testing.T) {
	elasticConfig := config.ElasticConfig{
		ElasticIndexName:            "notes",
		ElasticSearchValueMinLength: 2,
	}

	constraintsConfig := config.ConstraintsConfig{
		MaxDepth:         3,
		MaxCollaborators: 10,
		MaxTags:          4,
		MaxSubnotes:      10,
	}

	noteId := uuid.NewV4()
	userId := uuid.NewV4()
	guestId := uuid.NewV4()

	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo)
		wantErr    bool
	}{
		{
			name: "Test_AddCollaborator_FailOnRead",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{}, errors.New("error")).Times(1)

			},

			wantErr: true,
		},
		{
			name: "Test_AddCollaborator_FailOnNotOwner",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{}, nil).Times(1)

			},

			wantErr: true,
		},
		{
			name: "Test_AddCollaborator_FailParentNotEmpty",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						Id:      noteId,
						OwnerId: userId,
						Parent:  uuid.NewV4(),
					},
				}, nil).Times(1)

			},

			wantErr: true,
		},
		{
			name: "Test_AddCollaborator_FailAlreadyCollaborator",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						Id:            noteId,
						OwnerId:       userId,
						Parent:        uuid.UUID{},
						Collaborators: []uuid.UUID{guestId},
					},
				}, nil).Times(1)

			},

			wantErr: true,
		},
		{
			name: "Test_AddCollaborator_TooMany",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						Id:            noteId,
						OwnerId:       userId,
						Parent:        uuid.UUID{},
						Collaborators: make([]uuid.UUID, 20),
					},
				}, nil).Times(1)

			},

			wantErr: true,
		},
		{
			name: "Test_AddCollaborator_AddErr",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						Id:      noteId,
						OwnerId: userId,
						Parent:  uuid.UUID{},
					},
				}, nil).Times(1)
				baseRepo.EXPECT().AddCollaborator(gomock.Any(), noteId, guestId).Return("", errors.New("error"))

			},

			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo)

			_, err := uc.AddCollaborator(context.Background(), noteId, userId, guestId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.AddCollaborator error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNoteUsecase_getDepth(t *testing.T) {
	elasticConfig := config.ElasticConfig{
		ElasticIndexName:            "notes",
		ElasticSearchValueMinLength: 2,
	}

	constraintsConfig := config.ConstraintsConfig{
		MaxDepth:         3,
		MaxCollaborators: 10,
		MaxTags:          4,
		MaxSubnotes:      10,
	}

	parentId := uuid.NewV4()
	childId := uuid.NewV4()
	tests := []struct {
		parentId   uuid.UUID
		childId    uuid.UUID
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo)
		wantErr    bool
		currDepth  int
	}{
		{
			name:      "Test_getDepth",
			parentId:  uuid.UUID{},
			childId:   uuid.UUID{},
			currDepth: 1,
			wantErr:   false,
			repoMocker: func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
			},
		},
		{
			name:      "Test_ReadErr",
			parentId:  parentId,
			childId:   childId,
			currDepth: 1,
			wantErr:   true,
			repoMocker: func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(gomock.Any(), childId, gomock.Any()).Return(models.NoteResponse{}, errors.New("err"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo)

			_, err := uc.getDepth(context.Background(), tt.childId, tt.currDepth, uuid.UUID{})
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.getDepth error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestNoteUsecase_CreateSubNote(t *testing.T) {
	elasticConfig := config.ElasticConfig{
		ElasticIndexName:            "notes",
		ElasticSearchValueMinLength: 2,
	}

	constraintsConfig := config.ConstraintsConfig{
		MaxDepth:         3,
		MaxCollaborators: 10,
		MaxTags:          4,
		MaxSubnotes:      10,
	}

	parentId := uuid.NewV4()
	userId := uuid.NewV4()
	childId := uuid.NewV4()
	tests := []struct {
		parentId   uuid.UUID
		childId    uuid.UUID
		userId     uuid.UUID
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo)
		wantErr    bool
		currDepth  int
	}{
		{
			name:      "Test_CreateSubnote_readError",
			parentId:  parentId,
			childId:   childId,
			userId:    userId,
			currDepth: 1,
			wantErr:   true,
			repoMocker: func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(gomock.Any(), parentId, userId).Return(models.NoteResponse{}, errors.New("error"))
			},
		},
		{
			name:      "Test_CreateSubnote_TooManyChildren",
			parentId:  parentId,
			childId:   childId,
			userId:    userId,
			currDepth: 1,
			wantErr:   true,
			repoMocker: func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(gomock.Any(), parentId, userId).Return(models.NoteResponse{
					Note: models.Note{
						Id:       parentId,
						OwnerId:  userId,
						Children: make([]uuid.UUID, 20),
					},
				}, nil)

			},
		},
		{
			name:      "Test_CreateSubnote_getDepthErr",
			parentId:  parentId,
			childId:   childId,
			userId:    userId,
			currDepth: 1,
			wantErr:   true,
			repoMocker: func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(gomock.Any(), parentId, userId).Return(models.NoteResponse{
					Note: models.Note{
						Id:      parentId,
						OwnerId: userId,
						Parent:  uuid.NewV4(),
					},
				}, nil)
				baseRepo.EXPECT().ReadNote(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.NoteResponse{}, errors.New("error"))

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo)

			_, err := uc.CreateSubNote(context.Background(), tt.userId, "", tt.parentId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.CreateSubnote error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestNoteUsecase_CheckPermissions(t *testing.T) {
	elasticConfig := config.ElasticConfig{
		ElasticIndexName:            "notes",
		ElasticSearchValueMinLength: 2,
	}

	constraintsConfig := config.ConstraintsConfig{
		MaxDepth:         3,
		MaxCollaborators: 10,
		MaxTags:          10,
		MaxSubnotes:      10,
	}

	noteId := uuid.NewV4()
	userId := uuid.NewV4()

	type args struct {
		ctx    context.Context
		userId uuid.UUID
		noteId uuid.UUID
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo)
		args       args
		wantErr    bool
		want       bool
	}{
		{
			name: "TestSuccess",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: userId,
					},
				}, nil).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "TestFail",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{}, errors.New("error")).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo)

			got, err := uc.CheckPermissions(tt.args.ctx, tt.args.noteId, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.CheckPermissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNoteUsecase_SetIcon(t *testing.T) {
	elasticConfig := config.ElasticConfig{}

	constraintsConfig := config.ConstraintsConfig{}

	noteId := uuid.NewV4()
	userId := uuid.NewV4()

	type args struct {
		ctx    context.Context
		userId uuid.UUID
		noteId uuid.UUID
		icon   string
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args)
		args       args
		wantErr    bool
		want       models.Note
	}{
		{
			name: "Test_SetIcon_Success",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {

				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: args.userId,
						Id:      args.noteId,
						Icon:    args.icon,
					},
				}, nil).Times(1)
				baseRepo.EXPECT().SetIcon(ctx, args.noteId, args.icon).Return(nil)
				searchRepo.EXPECT().SetIcon(ctx, args.noteId, args.icon).Return(nil)

			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
				icon:   "icon",
			},
			wantErr: false,
			want: models.Note{
				OwnerId: userId,
				Id:      noteId,
				Icon:    "icon",
			},
		},
		{
			name: "Test_SetIcon_FailOnSet",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: args.userId,
						Id:      args.noteId,
						Icon:    args.icon,
					},
				}, nil).Times(1)
				baseRepo.EXPECT().SetIcon(ctx, args.noteId, args.icon).Return(errors.New("err"))

			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
				icon:   "icon",
			},
			wantErr: true,
			want:    models.Note{},
		},

		{
			name: "Test_SetIcon_FailOnRead",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: args.userId,
						Id:      args.noteId,
						Icon:    args.icon,
					},
				}, errors.New("err")).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    models.Note{},
		},
		{
			name: "Test_SetIcon_NotOwner",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: uuid.NewV4(),
						Id:      args.noteId,
					},
				}, nil).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    models.Note{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo, tt.args)

			got, err := uc.SetIcon(tt.args.ctx, tt.args.noteId, tt.args.icon, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.SetIcon error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNoteUsecase_SetHeader(t *testing.T) {
	elasticConfig := config.ElasticConfig{}

	constraintsConfig := config.ConstraintsConfig{}

	noteId := uuid.NewV4()
	userId := uuid.NewV4()

	type args struct {
		ctx    context.Context
		userId uuid.UUID
		noteId uuid.UUID
		header string
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args)
		args       args
		wantErr    bool
		want       models.Note
	}{
		{
			name: "Test_SetHeader_Success",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {

				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: args.userId,
						Id:      args.noteId,
						Header:  args.header,
					},
				}, nil).Times(1)
				baseRepo.EXPECT().SetHeader(ctx, args.noteId, args.header).Return(nil)
				searchRepo.EXPECT().SetHeader(ctx, args.noteId, args.header).Return(nil)

			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
				header: "header",
			},
			wantErr: false,
			want: models.Note{
				OwnerId: userId,
				Id:      noteId,
				Header:  "header",
			},
		},
		{
			name: "Test_SetHeader_FailOnSet",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: args.userId,
						Id:      args.noteId,
						Header:  args.header,
					},
				}, nil).Times(1)
				baseRepo.EXPECT().SetHeader(ctx, args.noteId, args.header).Return(errors.New("err"))

			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
				header: "header",
			},
			wantErr: true,
			want:    models.Note{},
		},

		{
			name: "Test_SetHeader_FailOnRead",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: args.userId,
						Id:      args.noteId,
						Header:  args.header,
					},
				}, errors.New("err")).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    models.Note{},
		},
		{
			name: "Test_SetHeader_NotOwner",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: uuid.NewV4(),
						Id:      args.noteId,
					},
				}, nil).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    models.Note{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo, tt.args)

			got, err := uc.SetHeader(tt.args.ctx, tt.args.noteId, tt.args.header, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.SetHeader error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNoteUsecase_AddFav(t *testing.T) {
	elasticConfig := config.ElasticConfig{}

	constraintsConfig := config.ConstraintsConfig{}

	noteId := uuid.NewV4()
	userId := uuid.NewV4()

	type args struct {
		ctx    context.Context
		userId uuid.UUID
		noteId uuid.UUID
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args)
		args       args
		wantErr    bool
		want       models.Note
	}{
		{
			name: "Test_AddFav_Success",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {

				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId:  args.userId,
						Id:       args.noteId,
						Favorite: false,
					},
				}, nil).Times(1)
				baseRepo.EXPECT().AddFav(ctx, args.noteId, args.userId).Return(nil)
				searchRepo.EXPECT().ChangeFlag(ctx, args.noteId, true).Return(nil)

			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: false,
			want: models.Note{
				OwnerId:  userId,
				Id:       noteId,
				Favorite: true,
			},
		},
		{
			name: "Test_AddFav_FailOnSet",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: args.userId,
						Id:      args.noteId,
					},
				}, nil).Times(1)
				baseRepo.EXPECT().AddFav(args.ctx, args.noteId, args.userId).Return(errors.New("err"))

			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    models.Note{},
		},

		{
			name: "Test_AddFav_FailOnRead",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: args.userId,
						Id:      args.noteId,
					},
				}, errors.New("err")).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    models.Note{},
		},
		{
			name: "Test_AddFav_NotOwner",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: uuid.NewV4(),
						Id:      args.noteId,
					},
				}, nil).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    models.Note{},
		},
		{
			name: "Test_AddFav_AlreadyTrue",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId:  args.userId,
						Id:       args.noteId,
						Favorite: true,
					},
				}, nil).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: false,
			want: models.Note{
				OwnerId:  userId,
				Id:       noteId,
				Favorite: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo, tt.args)

			got, err := uc.AddFav(context.Background(), tt.args.noteId, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.AddFav error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNoteUsecase_DelFav(t *testing.T) {
	elasticConfig := config.ElasticConfig{}

	constraintsConfig := config.ConstraintsConfig{}

	noteId := uuid.NewV4()
	userId := uuid.NewV4()

	type args struct {
		ctx    context.Context
		userId uuid.UUID
		noteId uuid.UUID
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args)
		args       args
		wantErr    bool
		want       models.Note
	}{
		{
			name: "Test_DelFav_Success",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {

				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId:  args.userId,
						Id:       args.noteId,
						Favorite: true,
					},
				}, nil).Times(1)
				baseRepo.EXPECT().DelFav(ctx, args.noteId, args.userId).Return(nil)
				searchRepo.EXPECT().ChangeFlag(ctx, args.noteId, false).Return(nil)

			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: false,
			want: models.Note{
				OwnerId:  userId,
				Id:       noteId,
				Favorite: false,
			},
		},
		{
			name: "Test_DelFav_FailOnDel",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId:  args.userId,
						Id:       args.noteId,
						Favorite: true,
					},
				}, nil).Times(1)
				baseRepo.EXPECT().DelFav(args.ctx, args.noteId, args.userId).Return(errors.New("err"))

			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    models.Note{},
		},

		{
			name: "Test_DelFav_FailOnRead",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: args.userId,
						Id:      args.noteId,
					},
				}, errors.New("err")).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    models.Note{},
		},
		{
			name: "Test_DelFav_NotOwner",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: uuid.NewV4(),
						Id:      args.noteId,
					},
				}, nil).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    models.Note{},
		},
		{
			name: "Test_DelFav_AlreadyTrue",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId:  args.userId,
						Id:       args.noteId,
						Favorite: false,
					},
				}, nil).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: false,
			want: models.Note{
				OwnerId:  userId,
				Id:       noteId,
				Favorite: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo, tt.args)

			got, err := uc.DelFav(context.Background(), tt.args.noteId, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.DelFav error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNoteUsecase_SetPublic(t *testing.T) {
	elasticConfig := config.ElasticConfig{}

	constraintsConfig := config.ConstraintsConfig{}

	noteId := uuid.NewV4()
	userId := uuid.NewV4()

	type args struct {
		ctx    context.Context
		userId uuid.UUID
		noteId uuid.UUID
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args)
		args       args
		wantErr    bool
		want       models.Note
	}{
		{
			name: "Test_SetPublic_Success",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {

				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: args.userId,
						Id:      args.noteId,
						Public:  false,
					},
				}, nil).Times(1)
				baseRepo.EXPECT().SetPublic(ctx, args.noteId).Return(nil)
				searchRepo.EXPECT().SetPublic(ctx, args.noteId).Return(nil)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: false,
			want: models.Note{
				OwnerId: userId,
				Id:      noteId,
				Public:  true,
			},
		},
		{
			name: "Test_SetPublic_FailOnSet",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId:  args.userId,
						Id:       args.noteId,
						Favorite: true,
					},
				}, nil).Times(1)
				baseRepo.EXPECT().SetPublic(ctx, args.noteId).Return(errors.New("err"))

			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    models.Note{},
		},

		{
			name: "Test_SetPublic_FailOnRead",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: args.userId,
						Id:      args.noteId,
					},
				}, errors.New("err")).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    models.Note{},
		},
		{
			name: "Test_SetPublic_NotOwner",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: uuid.NewV4(),
						Id:      args.noteId,
					},
				}, nil).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    models.Note{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo, tt.args)

			got, err := uc.SetPublic(context.Background(), tt.args.noteId, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.SetPublic error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNoteUsecase_SetPrivate(t *testing.T) {
	elasticConfig := config.ElasticConfig{}

	constraintsConfig := config.ConstraintsConfig{}

	noteId := uuid.NewV4()
	userId := uuid.NewV4()

	type args struct {
		ctx    context.Context
		userId uuid.UUID
		noteId uuid.UUID
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args)
		args       args
		wantErr    bool
		want       models.Note
	}{
		{
			name: "Test_SetPrivate_Success",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {

				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: args.userId,
						Id:      args.noteId,
						Public:  true,
					},
				}, nil).Times(1)
				baseRepo.EXPECT().SetPrivate(ctx, args.noteId).Return(nil)
				searchRepo.EXPECT().SetPrivate(ctx, args.noteId).Return(nil)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: false,
			want: models.Note{
				OwnerId: userId,
				Id:      noteId,
				Public:  false,
			},
		},
		{
			name: "Test_SetPrivate_FailOnSet",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId:  args.userId,
						Id:       args.noteId,
						Favorite: true,
					},
				}, nil).Times(1)
				baseRepo.EXPECT().SetPrivate(ctx, args.noteId).Return(errors.New("err"))

			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    models.Note{},
		},

		{
			name: "Test_SetPublic_FailOnRead",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: args.userId,
						Id:      args.noteId,
					},
				}, errors.New("err")).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    models.Note{},
		},
		{
			name: "Test_SetPrivate_NotOwner",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: uuid.NewV4(),
						Id:      args.noteId,
					},
				}, nil).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    models.Note{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo, tt.args)

			got, err := uc.SetPrivate(context.Background(), tt.args.noteId, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.SetPrivate error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNoteUsecase_GetAttachList(t *testing.T) {
	elasticConfig := config.ElasticConfig{}

	constraintsConfig := config.ConstraintsConfig{}

	noteId := uuid.NewV4()
	userId := uuid.NewV4()

	type args struct {
		ctx    context.Context
		userId uuid.UUID
		noteId uuid.UUID
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args)
		args       args
		wantErr    bool
		want       []string
	}{
		{
			name: "Test_GetAttachList_Success",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {

				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: args.userId,
						Id:      args.noteId,
						Public:  true,
					},
				}, nil).Times(1)
				baseRepo.EXPECT().GetAttachList(ctx, args.noteId).Return([]string{"1", "2"}, nil)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: false,
			want:    []string{"1", "2"},
		},
		{
			name: "Test_GetAttachList_FailOnGet",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId:  args.userId,
						Id:       args.noteId,
						Favorite: true,
					},
				}, nil).Times(1)
				baseRepo.EXPECT().GetAttachList(ctx, args.noteId).Return([]string{}, errors.New("err"))

			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    []string{},
		},

		{
			name: "Test_GetAttachList_FailOnRead",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: args.userId,
						Id:      args.noteId,
					},
				}, errors.New("err")).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    []string{},
		},
		{
			name: "Test_GetAttachList_NotOwner",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadNote(ctx, noteId, userId).Return(models.NoteResponse{
					Note: models.Note{
						OwnerId: uuid.NewV4(),
						Id:      args.noteId,
					},
				}, nil).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,
				noteId: noteId,
			},
			wantErr: true,
			want:    []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo, tt.args)

			got, err := uc.GetAttachList(context.Background(), tt.args.noteId, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.GetAttachList error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNoteUsecase_RememberTag(t *testing.T) {
	elasticConfig := config.ElasticConfig{}

	constraintsConfig := config.ConstraintsConfig{}

	userId := uuid.NewV4()

	type args struct {
		ctx    context.Context
		userId uuid.UUID

		tagName string
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args)
		args       args
		wantErr    bool
	}{
		{
			name: "Test_RememberTag_Success",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().RememberTag(ctx, args.tagName, args.userId).Return(nil).Times(1)

			},
			args: args{
				ctx:    context.Background(),
				userId: userId,

				tagName: "tag1",
			},
			wantErr: false,
		},
		{
			name: "Test_RememberTag_FailOnAdd",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().RememberTag(ctx, args.tagName, args.userId).Return(errors.New("error")).Times(1)

			},
			args: args{
				ctx:    context.Background(),
				userId: userId,

				tagName: "tag1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo, tt.args)

			err := uc.RememberTag(tt.args.ctx, tt.args.tagName, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.RememberTag error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNoteUsecase_ForgetTag(t *testing.T) {
	elasticConfig := config.ElasticConfig{}

	constraintsConfig := config.ConstraintsConfig{}

	userId := uuid.NewV4()

	type args struct {
		ctx    context.Context
		userId uuid.UUID

		tagName string
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args)
		args       args
		wantErr    bool
	}{
		{
			name: "Test_ForgetTag_Success",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ForgetTag(ctx, args.tagName, args.userId).Return(nil).Times(1)
				baseRepo.EXPECT().DeleteTagFromAllNotes(ctx, args.tagName, args.userId)
				searchRepo.EXPECT().DeleteTagFromAllNotes(ctx, args.tagName, args.userId)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,

				tagName: "tag1",
			},
			wantErr: false,
		},
		{
			name: "Test_ForgetTag_FailOnForget",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ForgetTag(ctx, args.tagName, args.userId).Return(errors.New("error")).Times(1)

			},
			args: args{
				ctx:    context.Background(),
				userId: userId,

				tagName: "tag1",
			},
			wantErr: true,
		},
		{
			name: "Test_ForgetTag_FailOnForget",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ForgetTag(ctx, args.tagName, args.userId).Return(nil).Times(1)
				baseRepo.EXPECT().DeleteTagFromAllNotes(ctx, args.tagName, args.userId).Return(errors.New("err"))

			},
			args: args{
				ctx:    context.Background(),
				userId: userId,

				tagName: "tag1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo, tt.args)

			err := uc.ForgetTag(tt.args.ctx, tt.args.tagName, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.RememberTag error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNoteUsecase_UpdateTag(t *testing.T) {
	elasticConfig := config.ElasticConfig{}

	constraintsConfig := config.ConstraintsConfig{}

	userId := uuid.NewV4()

	type args struct {
		ctx    context.Context
		userId uuid.UUID

		tagName string
	}
	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args)
		args       args
		wantErr    bool
	}{
		{
			name: "Test_UpdateTag_Success",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().UpdateTag(ctx, "", args.tagName, args.userId).Return(nil)
				searchRepo.EXPECT().UpdateTagOnAllNotes(ctx, "", args.tagName, args.userId).Return(nil)
				baseRepo.EXPECT().UpdateTagOnAllNotes(ctx, "", args.tagName, args.userId).Return(nil)
			},
			args: args{
				ctx:    context.Background(),
				userId: userId,

				tagName: "tag1",
			},
			wantErr: false,
		},
		{
			name: "Test_UpdatetTag_FailOnUpdate",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().UpdateTag(ctx, "", args.tagName, args.userId).Return(errors.New("error")).Times(1)

			},
			args: args{
				ctx:    context.Background(),
				userId: userId,

				tagName: "tag1",
			},
			wantErr: true,
		},
		{
			name: "Test_UpdatetTag_FailOnUpdateAll",
			repoMocker: func(ctx context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().UpdateTag(ctx, "", args.tagName, args.userId).Return(nil).Times(1)
				baseRepo.EXPECT().UpdateTagOnAllNotes(ctx, "", args.tagName, args.userId).Return(errors.New("err"))

			},
			args: args{
				ctx:    context.Background(),
				userId: userId,

				tagName: "tag1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo, tt.args)

			err := uc.UpdateTag(tt.args.ctx, "", tt.args.tagName, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.UpdateTag error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNoteUsecase_GetSharedAttachList(t *testing.T) {
	elasticConfig := config.ElasticConfig{}

	constraintsConfig := config.ConstraintsConfig{}

	noteId := uuid.NewV4()
	type args struct {
		noteID uuid.UUID
	}
	tests := []struct {
		name       string
		args       args
		repoMocker func(baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args)
		want       []string
		wantErr    bool
	}{
		{
			name: "Test_GetSharedAttachList_Success",
			args: args{
				noteID: noteId,
			},
			repoMocker: func(baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadPublicNote(gomock.Any(), args.noteID).Return(models.NoteResponse{
					Note: models.Note{
						Public: true,
						Id:     args.noteID,
					},
				}, nil)
				baseRepo.EXPECT().GetAttachList(gomock.Any(), args.noteID).Return([]string{"1", "2"}, nil)
			},
			want:    []string{"1", "2"},
			wantErr: false,
		},
		{
			name: "Test_GetSharedAttachList_NotPublic",
			args: args{
				noteID: noteId,
			},
			repoMocker: func(baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadPublicNote(gomock.Any(), args.noteID).Return(models.NoteResponse{
					Note: models.Note{
						Public: false,
					},
				}, nil)

			},
			want:    []string{},
			wantErr: true,
		},
		{
			name: "Test_GetSharedAttachList_FailRead",
			args: args{
				noteID: noteId,
			},
			repoMocker: func(baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadPublicNote(gomock.Any(), args.noteID).Return(models.NoteResponse{}, errors.New("err"))

			},
			want:    []string{},
			wantErr: true,
		},
		{
			name: "Test_GetSharedAttachList_FailGetList",
			args: args{
				noteID: noteId,
			},
			repoMocker: func(baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo, args args) {
				baseRepo.EXPECT().ReadPublicNote(gomock.Any(), args.noteID).Return(models.NoteResponse{
					Note: models.Note{
						Public: true,
					},
				}, nil)
				baseRepo.EXPECT().GetAttachList(gomock.Any(), args.noteID).Return([]string{}, errors.New("err"))

			},
			want:    []string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(repo, searchRepo, tt.args)
			got, err := uc.GetSharedAttachList(context.Background(), tt.args.noteID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.GetSharedAttachList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNoteUsecase_changeModeRecursive(t *testing.T) {
	elasticConfig := config.ElasticConfig{
		ElasticIndexName:            "notes",
		ElasticSearchValueMinLength: 2,
	}

	constraintsConfig := config.ConstraintsConfig{
		MaxDepth:         3,
		MaxCollaborators: 10,
		MaxTags:          4,
		MaxSubnotes:      10,
	}

	noteId := uuid.NewV4()
	guestId := uuid.NewV4()

	tests := []struct {
		name       string
		repoMocker func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo)
		isPublic   bool
		wantErr    bool
	}{
		{
			name:    "Test_changeModeRecursive_ReadError",
			wantErr: true,
			repoMocker: func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(gomock.Any(), noteId, gomock.Any()).Return(models.NoteResponse{}, errors.New("error"))
			},
			isPublic: true,
		},
		{
			name:    "Test_changeModeRecursive_Public_Error",
			wantErr: true,
			repoMocker: func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(gomock.Any(), noteId, gomock.Any()).Return(models.NoteResponse{
					Note: models.Note{
						Id: noteId,
					},
				}, nil)
				baseRepo.EXPECT().SetPublic(gomock.Any(), noteId).Return(errors.New("error"))
			},
			isPublic: true,
		},
		{
			name:    "Test_changeModeRecursive_Private_Error",
			wantErr: true,
			repoMocker: func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(gomock.Any(), noteId, gomock.Any()).Return(models.NoteResponse{
					Note: models.Note{
						Id: noteId,
					},
				}, nil)
				baseRepo.EXPECT().SetPrivate(gomock.Any(), noteId).Return(errors.New("error"))
			},
			isPublic: false,
		},
		{
			name:    "Test_changeModeRecursive_NoChildren",
			wantErr: false,
			repoMocker: func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(gomock.Any(), noteId, gomock.Any()).Return(models.NoteResponse{
					Note: models.Note{
						Id: noteId,
					},
				}, nil)
				baseRepo.EXPECT().SetPublic(gomock.Any(), noteId).Return(nil)
			},
			isPublic: true,
		},
		{
			name:    "Test_changeModeRecursive_Public_Success",
			wantErr: false,
			repoMocker: func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(gomock.Any(), noteId, gomock.Any()).Return(models.NoteResponse{
					Note: models.Note{
						Id:       noteId,
						Children: []uuid.UUID{noteId},
					},
				}, nil)
				baseRepo.EXPECT().ReadNote(gomock.Any(), noteId, gomock.Any()).Return(models.NoteResponse{
					Note: models.Note{
						Id: noteId,
					},
				}, nil)
				baseRepo.EXPECT().SetPublic(gomock.Any(), noteId).Return(nil).Times(2)
			},
			isPublic: true,
		},
		{
			name:    "Test_changeModeRecursive_Private_Success",
			wantErr: false,
			repoMocker: func(context context.Context, baseRepo *mock_note.MockNoteBaseRepo, searchRepo *mock_note.MockNoteSearchRepo) {
				baseRepo.EXPECT().ReadNote(gomock.Any(), noteId, gomock.Any()).Return(models.NoteResponse{
					Note: models.Note{
						Id:       noteId,
						Children: []uuid.UUID{noteId},
					},
				}, nil)
				baseRepo.EXPECT().ReadNote(gomock.Any(), noteId, gomock.Any()).Return(models.NoteResponse{
					Note: models.Note{
						Id: noteId,
					},
				}, nil)
				baseRepo.EXPECT().SetPrivate(gomock.Any(), noteId).Return(nil).Times(2)
			},
			isPublic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteBaseRepo(ctl)
			searchRepo := mock_note.NewMockNoteSearchRepo(ctl)
			uc := CreateNoteUsecase(repo, searchRepo, elasticConfig, constraintsConfig, &sync.WaitGroup{})

			tt.repoMocker(context.Background(), repo, searchRepo)

			err := uc.changeModeRecursive(context.Background(), noteId, guestId, tt.isPublic)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.changeModeRecursive error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
