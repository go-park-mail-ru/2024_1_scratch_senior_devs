package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	mock_note "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/mocks"
	"github.com/golang/mock/gomock"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/require"
)

func TestNoteUsecase_GetAllNotes(t *testing.T) {
	/*type fields struct {
		repo note.NoteRepo
	}
	type args struct {
		ctx    context.Context
		userId uuid.UUID
		count  int64
		offset int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Note
		wantErr bool
	}{
		// TODO: Add test cases.
	}*/
	ctx := context.Background()
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	repo := mock_note.NewMockNoteRepo(ctl)
	uId := uuid.NewV4()
	expected := []models.Note{ //ожидаемый от юзкейса результат.
		{
			Id:         uuid.NewV4(),
			OwnerId:    uId,
			UpdateTime: nil,
			CreateTime: time.Time{},
			Data:       nil,
		},
		{
			Id:         uuid.NewV4(),
			OwnerId:    uId,
			UpdateTime: nil,
			CreateTime: time.Time{},
			Data:       nil,
		},
	}
	mockResp := []models.Note{ //мок ответа от урованя репозитория
		{
			Id:         uuid.NewV4(),
			OwnerId:    uId,
			UpdateTime: nil,
			CreateTime: time.Time{},
			Data:       nil,
		},
		{
			Id:         uuid.NewV4(),
			OwnerId:    uId,
			UpdateTime: nil,
			CreateTime: time.Time{},
			Data:       nil,
		},
	}
	repo.EXPECT().ReadAllNotes(ctx, uId, int64(10), int64(0)).Return(mockResp, nil).Times(1)
	Usecase := CreateNoteUsecase(repo)
	notes, err := Usecase.GetAllNotes(ctx, uId, int64(10), int64(0)) //Get(ctx, log, in)
	require.NoError(t, err)
	require.ElementsMatch(t, expected, notes)
	/*for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &NoteUsecase{
				repo: tt.fields.repo,
			}
			got, err := uc.GetAllNotes(tt.args.ctx, tt.args.userId, tt.args.count, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteUsecase.GetAllNotes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NoteUsecase.GetAllNotes() = %v, want %v", got, tt.want)
			}
		})
	}*/
}
