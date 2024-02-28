package usecase

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	mock_note "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/mocks"
	"github.com/golang/mock/gomock"
	"github.com/satori/uuid"
)

func TestNoteUsecase_GetAllNotes(t *testing.T) {
	uids := []uuid.UUID{
		uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
		uuid.FromStringOrNil("a839229c-4521-4ea9-9c21-90f0d2715568"),

		uuid.FromStringOrNil("8d48b193-d50a-4df2-8acf-3c3f34853a28"),
	}
	type args struct {
		ctx    context.Context
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
				mockResp := []models.Note{ //мок ответа от урованя репозитория
					{
						Id:         uids[0],
						OwnerId:    uids[2],
						UpdateTime: nil,
						CreateTime: time.Time{},
						Data:       nil,
					},
					{
						Id:         uids[1],
						OwnerId:    uids[2],
						UpdateTime: nil,
						CreateTime: time.Time{},
						Data:       nil,
					},
				}

				repo.EXPECT().ReadAllNotes(ctx, uId, int64(count), int64(offset)).Return(mockResp, nil).Times(1)
			},
			args: args{
				context.Background(),
				uuid.NewV4(),
				10,
				0,
			},
			want: []models.Note{
				{
					Id:         uids[0],
					OwnerId:    uids[2], //потом задать из args
					UpdateTime: nil,
					CreateTime: time.Time{},
					Data:       nil,
				},
				{
					Id:         uids[1],
					OwnerId:    uids[2],
					UpdateTime: nil,
					CreateTime: time.Time{},
					Data:       nil,
				},
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mock_note.NewMockNoteRepo(ctl)
			Usecase := CreateNoteUsecase(repo)

			tt.repoMocker(context.Background(), repo, tt.args.userId, tt.args.count, tt.args.offset)
			got, err := Usecase.GetAllNotes(tt.args.ctx, tt.args.userId, tt.args.count, tt.args.offset)
			//got, err := uc.GetAllNotes(tt.args.ctx, tt.args.userId, tt.args.count, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.GetAllNotes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NoteUsecase.GetAllNotes() = %v, want %v", got, tt.want)
			}
		})
	}
}
