package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	mock_metrics "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/metrics/mocks"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNoteRepo_ReadAllNotes(t *testing.T) {
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics, pgx.Rows, uuid.UUID)
		userId         uuid.UUID
		columns        []string
		expectedErr    error
	}{
		{
			name: "ReadAllNotes_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, id uuid.UUID) {
				mockPool.EXPECT().Query(gomock.Any(), getAllNotes, id, int64(1), int64(0), []string{}).Return(pgxRows, nil)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			userId:      uuid.NewV4(),
			columns:     []string{"id", "data", "create_time", "update_time", "owner_id", "parent", "children", "tags", "collaborators", "icon", "header", "favorite", "is_public"},
			expectedErr: nil,
		},
		{
			name: "ReadAllNotes_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, id uuid.UUID) {
				mockPool.EXPECT().Query(gomock.Any(), getAllNotes, id, int64(1), int64(0), []string{}).Return(pgxRows, pgx.ErrNoRows)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},
			userId:      uuid.NewV4(),
			columns:     []string{"id", "data", "create_time", "update_time", "owner_id", "parent", "children", "tags", "collaborators", "icon", "header", "favorite", "is_public"},
			expectedErr: pgx.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			pgxRows := pgxpoolmock.NewRows(tt.columns).AddRow(uuid.NewV4(), "", time.Now(), time.Time{}, tt.userId, uuid.UUID{}, []uuid.UUID{}, []string{}, []uuid.UUID{}, "", "", false, false).ToPgxRows()

			tt.mockRepoAction(mockPool, mockMetrics, pgxRows, tt.userId)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			_, err := repo.ReadAllNotes(context.Background(), tt.userId, int64(1), int64(0), []string{})

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestNoteRepo_ReadNote(t *testing.T) {
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics, pgx.Rows, uuid.UUID)
		Id             uuid.UUID
		columns        []string
		expectedErr    error
	}{
		{
			name: "ReadNote_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, id uuid.UUID) {
				mockPool.EXPECT().QueryRow(gomock.Any(), getNote, id).Return(pgxRows)
				pgxRows.Next()
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			Id:          uuid.NewV4(),
			columns:     []string{"id", "data", "create_time", "update_time", "owner_id", "parent", "children", "tags", "collaborators", "icon", "header", "favorite", "is_public"},
			expectedErr: nil,
		},
		{
			name: "ReadNote_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, id uuid.UUID) {
				mockPool.EXPECT().QueryRow(gomock.Any(), getNote, id).Return(pgxRows)
				pgxRows.Next()
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			Id:          uuid.NewV4(),
			columns:     []string{"id", "data", "create_time", "update_time", "owner_id", "parent", "children", "tags", "collaborators", "icon", "header", "favorite", "is_public"},
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			pgxRows := pgxpoolmock.NewRows(tt.columns).AddRow(uuid.NewV4(), "", time.Now(), time.Time{}, tt.Id, uuid.UUID{}, []uuid.UUID{}, []string{}, []uuid.UUID{}, "", "", false, false).ToPgxRows()

			tt.mockRepoAction(mockPool, mockMetrics, pgxRows, tt.Id)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			_, err := repo.ReadNote(context.Background(), tt.Id, uuid.UUID{})

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestNoteRepo_CreateNote(t *testing.T) {
	Id := uuid.NewV4()
	userId := uuid.NewV4()
	currTime := time.Now().UTC()

	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics)
		err            error
	}{
		{
			name: "CreateNote_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
				mockPool.EXPECT().Exec(gomock.Any(), createNote,
					Id, "", currTime, currTime, userId, Id, []uuid.UUID{}, []string{}, []uuid.UUID{}, "", "", false,
				).Return(nil, nil)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			err: nil,
		},
		{
			name: "CreateNote_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
				mockPool.EXPECT().Exec(gomock.Any(), createNote,
					Id, "", currTime, currTime, userId, Id, []uuid.UUID{}, []string{}, []uuid.UUID{}, "", "", false,
				).Return(nil, errors.New("err"))
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},
			err: errors.New("err"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool, mockMetrics)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			err := repo.CreateNote(context.Background(), models.Note{
				Id:            Id,
				Data:          "",
				CreateTime:    currTime,
				UpdateTime:    currTime,
				OwnerId:       userId,
				Parent:        Id,
				Children:      []uuid.UUID{},
				Tags:          []string{},
				Collaborators: []uuid.UUID{},
				Icon:          "",
				Header:        "",
				Favorite:      false,
				Public:        false,
			})

			assert.Equal(t, tt.err, err)
		})
	}
}

func TestNoteRepo_UpdateNote(t *testing.T) {
	Id := uuid.NewV4()
	userId := uuid.NewV4()
	currTime := time.Now().UTC()

	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics)
		err            error
	}{
		{
			name: "UpdateNote_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
				mockPool.EXPECT().Exec(gomock.Any(), updateNote,
					"", currTime, Id,
				).Return(nil, nil)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			err: nil,
		},
		{
			name: "UpdateNote_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
				mockPool.EXPECT().Exec(gomock.Any(), updateNote,
					"", currTime, Id,
				).Return(nil, pgx.ErrNoRows)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},
			err: pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool, mockMetrics)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			err := repo.UpdateNote(context.Background(), models.Note{
				Id:            Id,
				Data:          "",
				CreateTime:    currTime,
				UpdateTime:    currTime,
				OwnerId:       userId,
				Parent:        Id,
				Children:      []uuid.UUID{},
				Tags:          []string{},
				Collaborators: []uuid.UUID{},
				Icon:          "",
				Header:        "",
				Favorite:      false,
				Public:        false,
			})

			assert.Equal(t, tt.err, err)
		})
	}
}

func TestNoteRepo_DeleteNote(t *testing.T) {
	Id := uuid.NewV4()

	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics)
		err            error
	}{
		{
			name: "DeleteNote_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
				mockPool.EXPECT().Exec(gomock.Any(), deleteNote, Id).Return(nil, nil)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			err: nil,
		},
		{
			name: "DeleteNote_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
				mockPool.EXPECT().Exec(gomock.Any(), deleteNote, Id).Return(nil, pgx.ErrNoRows)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},
			err: pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool, mockMetrics)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			err := repo.DeleteNote(context.Background(), Id)

			assert.Equal(t, tt.err, err)
		})
	}
}

func TestNotePostgres_AddSubNote(t *testing.T) {
	noteId := uuid.NewV4()
	parentId := uuid.NewV4()

	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics)
		err            error
	}{
		{
			name: "AddSubNote_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
				mockPool.EXPECT().Exec(gomock.Any(), addSubNote, parentId, noteId).Return(nil, nil)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			err: nil,
		},
		{
			name: "AddSubNote_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
				mockPool.EXPECT().Exec(gomock.Any(), addSubNote, parentId, noteId).Return(nil, pgx.ErrNoRows)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},
			err: pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool, mockMetrics)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			err := repo.AddSubNote(context.Background(), noteId, parentId)

			assert.Equal(t, tt.err, err)
		})
	}
}

func TestNotePostgres_RemoveSubNote(t *testing.T) {
	noteId := uuid.NewV4()
	parentId := uuid.NewV4()

	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics)
		err            error
	}{
		{
			name: "RemoveSubNote_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
				mockPool.EXPECT().Exec(gomock.Any(), removeSubNote, noteId, parentId).Return(nil, nil)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			err: nil,
		},
		{
			name: "RemoveSubNote_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
				mockPool.EXPECT().Exec(gomock.Any(), removeSubNote, noteId, parentId).Return(nil, pgx.ErrNoRows)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},
			err: pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool, mockMetrics)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			err := repo.RemoveSubNote(context.Background(), parentId, noteId)

			assert.Equal(t, tt.err, err)
		})
	}
}

func TestNotePostgres_AddTag(t *testing.T) {
	noteId := uuid.NewV4()

	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics)
		err            error
	}{
		{
			name: "AddTag_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
				mockPool.EXPECT().Exec(gomock.Any(), addTag, "tag", noteId).Return(nil, nil)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			err: nil,
		},
		{
			name: "AddTag_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
				mockPool.EXPECT().Exec(gomock.Any(), addTag, "tag", noteId).Return(nil, pgx.ErrNoRows)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},
			err: pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool, mockMetrics)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			err := repo.AddTag(context.Background(), "tag", noteId)

			assert.Equal(t, tt.err, err)
		})
	}
}

func TestNotePostgres_DeleteTag(t *testing.T) {
	noteId := uuid.NewV4()

	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics)
		err            error
	}{
		{
			name: "DeleteTag_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
				mockPool.EXPECT().Exec(gomock.Any(), deleteTag, "tag", noteId).Return(nil, nil)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			err: nil,
		},
		{
			name: "DeleteTag_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
				mockPool.EXPECT().Exec(gomock.Any(), deleteTag, "tag", noteId).Return(nil, pgx.ErrNoRows)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},
			err: pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool, mockMetrics)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			err := repo.DeleteTag(context.Background(), "tag", noteId)

			assert.Equal(t, tt.err, err)
		})
	}
}

//func TestNotePostgres_AddCollaborator(t *testing.T) {
//	noteId := uuid.NewV4()
//	guestId := uuid.NewV4()
//
//	tests := []struct {
//		name           string
//		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics)
//		err            error
//	}{
//		{
//			name: "AddCollaborator_Success",
//			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
//				mockPool.EXPECT().QueryRow(gomock.Any(), addCollaborator, guestId, noteId).Return(pgxpoolmock.NewRows([]string{"title"}).AddRow("title").ToPgxRows())
//				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
//			},
//			err: nil,
//		},
//		{
//			name: "AddCollaborator_Fail",
//			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
//				mockPool.EXPECT().QueryRow(gomock.Any(), addCollaborator, guestId, noteId).Return(pgxpoolmock.NewRows([]string{"title"}).AddRow("title").ToPgxRows())
//				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
//				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
//			},
//			err: pgx.ErrNoRows,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
//			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
//			defer ctrl.Finish()
//
//			tt.mockRepoAction(mockPool, mockMetrics)
//
//			repo := CreateNotePostgres(mockPool, mockMetrics)
//			_, err := repo.AddCollaborator(context.Background(), noteId, guestId)
//
//			assert.Equal(t, tt.err, err)
//		})
//	}
//}

func TestNotePostgres_GetTags(t *testing.T) {
	userId := uuid.NewV4()
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics, pgx.Rows, uuid.UUID)
		userId         uuid.UUID
		expectedErr    error
	}{
		{
			name: "GetTags_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, id uuid.UUID) {
				mockPool.EXPECT().Query(gomock.Any(), getTags, userId).Return(pgxRows, nil)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			userId:      userId,
			expectedErr: nil,
		},
		{
			name: "GetTags_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, id uuid.UUID) {
				mockPool.EXPECT().Query(gomock.Any(), getTags, userId).Return(pgxRows, pgx.ErrNoRows)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},
			userId:      userId,
			expectedErr: pgx.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			pgxRows := pgxpoolmock.NewRows([]string{"tags"}).AddRow("tag").ToPgxRows()

			tt.mockRepoAction(mockPool, mockMetrics, pgxRows, tt.userId)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			_, err := repo.GetTags(context.Background(), tt.userId)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestNotePostgres_RememberTag(t *testing.T) {
	userId := uuid.NewV4()
	type args struct {
		userId  uuid.UUID
		tagName string
	}
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics, pgx.Rows, args)
		args           args
		expectedErr    error
	}{
		{
			name: "RememberTag_Success",
			args: args{
				userId:  userId,
				tagName: "tag",
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, args args) {
				commandTag := pgconn.CommandTag("INSERT 0 1")
				mockPool.EXPECT().Exec(gomock.Any(), rememberTag, args.tagName, args.userId).Return(commandTag, nil)

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},

			expectedErr: nil,
		},
		{
			name: "RememberTag_Fail_No_rows_affected",
			args: args{
				userId:  userId,
				tagName: "tag",
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), rememberTag, args.tagName, args.userId).Return(nil, nil)

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},

			expectedErr: errors.New("tag already exists"),
		},
		{
			name: "RememberTag_Fail_Error",
			args: args{
				userId:  userId,
				tagName: "tag",
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), rememberTag, args.tagName, args.userId).Return(nil, errors.New("err"))

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},

			expectedErr: errors.New("err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			pgxRows := pgxpoolmock.NewRows([]string{"tag_name", "user_id"}).AddRow("t", "0").ToPgxRows()

			tt.mockRepoAction(mockPool, mockMetrics, pgxRows, tt.args)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			err := repo.RememberTag(context.Background(), tt.args.tagName, tt.args.userId)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestNotePostgres_ForgetTag(t *testing.T) {
	userId := uuid.NewV4()
	type args struct {
		userId  uuid.UUID
		tagName string
	}
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics, pgx.Rows, args)
		args           args
		expectedErr    error
	}{
		{
			name: "ForgetTag_Success",
			args: args{
				userId:  userId,
				tagName: "tag",
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), forgetTag, args.tagName, args.userId).Return(nil, nil)

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},

			expectedErr: nil,
		},

		{
			name: "ForgetTag_Fail_Error",
			args: args{
				userId:  userId,
				tagName: "tag",
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), forgetTag, args.tagName, args.userId).Return(nil, errors.New("err"))

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},

			expectedErr: errors.New("err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			pgxRows := pgxpoolmock.NewRows([]string{"tag_name", "user_id"}).AddRow("t", "0").ToPgxRows()

			tt.mockRepoAction(mockPool, mockMetrics, pgxRows, tt.args)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			err := repo.ForgetTag(context.Background(), tt.args.tagName, tt.args.userId)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestNotePostgres_DeleteTagFromAllNotes(t *testing.T) {
	userId := uuid.NewV4()
	type args struct {
		userId  uuid.UUID
		tagName string
	}
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics, pgx.Rows, args)
		args           args
		expectedErr    error
	}{
		{
			name: "Test_DeleteTagFromAllNotes_Success",
			args: args{
				userId:  userId,
				tagName: "tag",
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), deleteTagFromAllNotes, args.tagName, args.userId).Return(nil, nil)

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},

			expectedErr: nil,
		},

		{
			name: "Test_DeleteTagFromAllNotes_Fail_Error",
			args: args{
				userId:  userId,
				tagName: "tag",
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), deleteTagFromAllNotes, args.tagName, args.userId).Return(nil, errors.New("err"))

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},

			expectedErr: errors.New("err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			pgxRows := pgxpoolmock.NewRows([]string{"tag_name", "user_id"}).AddRow("t", "0").ToPgxRows()

			tt.mockRepoAction(mockPool, mockMetrics, pgxRows, tt.args)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			err := repo.DeleteTagFromAllNotes(context.Background(), tt.args.tagName, tt.args.userId)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestNotePostgres_UpdateTag(t *testing.T) {
	userId := uuid.NewV4()
	type args struct {
		userId  uuid.UUID
		oldName string
		newName string
	}
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics, pgx.Rows, args)
		args           args
		expectedErr    error
	}{
		{
			name: "Test_UpdateTag_Success",
			args: args{
				userId:  userId,
				oldName: "tag",
				newName: "tag1",
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), updateTag, args.newName, args.userId, args.oldName).Return(nil, nil)

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},

			expectedErr: nil,
		},

		{
			name: "Test_UpdateTag_Fail_Error",
			args: args{
				userId:  userId,
				oldName: "tag",
				newName: "tag1",
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), updateTag, args.newName, args.userId, args.oldName).Return(nil, errors.New("err"))

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},

			expectedErr: errors.New("err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			pgxRows := pgxpoolmock.NewRows([]string{"tag_name", "user_id"}).AddRow("t", "0").ToPgxRows()

			tt.mockRepoAction(mockPool, mockMetrics, pgxRows, tt.args)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			err := repo.UpdateTag(context.Background(), tt.args.oldName, tt.args.newName, tt.args.userId)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestNotePostgres_SetIcon(t *testing.T) {
	noteId := uuid.NewV4()
	type args struct {
		noteId uuid.UUID
		icon   string
	}
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics, pgx.Rows, args)
		args           args
		expectedErr    error
	}{
		{
			name: "Test_SetIcon_Success",
			args: args{
				noteId: noteId,
				icon:   "icon",
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), setIcon, args.icon, args.noteId).Return(nil, nil)

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},

			expectedErr: nil,
		},

		{
			name: "Test_SetIcons_Fail_Error",
			args: args{
				noteId: noteId,
				icon:   "icon",
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), setIcon, args.icon, args.noteId).Return(nil, errors.New("err"))

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},

			expectedErr: errors.New("err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			pgxRows := pgxpoolmock.NewRows([]string{"tag_name", "user_id"}).AddRow("t", "0").ToPgxRows()

			tt.mockRepoAction(mockPool, mockMetrics, pgxRows, tt.args)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			err := repo.SetIcon(context.Background(), tt.args.noteId, tt.args.icon)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestNotePostgres_SetHeader(t *testing.T) {
	noteId := uuid.NewV4()
	type args struct {
		noteId uuid.UUID
		header string
	}
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics, args)
		args           args
		expectedErr    error
	}{
		{
			name: "Test_SetHeader_Success",
			args: args{
				noteId: noteId,
				header: "header",
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), setHeader, args.header, args.noteId).Return(nil, nil)

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},

			expectedErr: nil,
		},

		{
			name: "Test_SetHeader_Fail_Error",
			args: args{
				noteId: noteId,
				header: "header",
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), setHeader, args.header, args.noteId).Return(nil, errors.New("err"))

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},

			expectedErr: errors.New("err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool, mockMetrics, tt.args)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			err := repo.SetHeader(context.Background(), tt.args.noteId, tt.args.header)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestNotePostgres_AddFav(t *testing.T) {
	noteId := uuid.NewV4()
	userId := uuid.NewV4()
	type args struct {
		noteId uuid.UUID
		userId uuid.UUID
	}
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics, args)
		args           args
		expectedErr    error
	}{
		{
			name: "Test_AddFav_Success",
			args: args{
				noteId: noteId,
				userId: userId,
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), addFav, args.noteId, args.userId).Return(nil, nil)

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},

			expectedErr: nil,
		},

		{
			name: "Test_AddFav_Fail_Error",
			args: args{
				noteId: noteId,
				userId: userId,
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), addFav, args.noteId, args.userId).Return(nil, errors.New("err"))

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},

			expectedErr: errors.New("err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool, mockMetrics, tt.args)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			err := repo.AddFav(context.Background(), tt.args.noteId, tt.args.userId)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestNotePostgres_DelFav(t *testing.T) {
	noteId := uuid.NewV4()
	userId := uuid.NewV4()
	type args struct {
		noteId uuid.UUID
		userId uuid.UUID
	}
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics, args)
		args           args
		expectedErr    error
	}{
		{
			name: "Test_DelFav_Success",
			args: args{
				noteId: noteId,
				userId: userId,
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), delFav, args.noteId, args.userId).Return(nil, nil)

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},

			expectedErr: nil,
		},

		{
			name: "Test_DelFav_Fail_Error",
			args: args{
				noteId: noteId,
				userId: userId,
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), delFav, args.noteId, args.userId).Return(nil, errors.New("err"))

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},

			expectedErr: errors.New("err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool, mockMetrics, tt.args)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			err := repo.DelFav(context.Background(), tt.args.noteId, tt.args.userId)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestNotePostgres_SetPublic(t *testing.T) {
	noteId := uuid.NewV4()

	type args struct {
		noteId uuid.UUID
	}
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics, args)
		args           args
		expectedErr    error
	}{
		{
			name: "Test_SetPublic_Success",
			args: args{
				noteId: noteId,
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), setPublic, args.noteId).Return(nil, nil)

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},

			expectedErr: nil,
		},

		{
			name: "Test_SetPublic_Fail_Error",
			args: args{
				noteId: noteId,
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), setPublic, args.noteId).Return(nil, errors.New("err"))

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},

			expectedErr: errors.New("err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool, mockMetrics, tt.args)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			err := repo.SetPublic(context.Background(), tt.args.noteId)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestNotePostgres_SetPrivate(t *testing.T) {
	noteId := uuid.NewV4()

	type args struct {
		noteId uuid.UUID
	}
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics, args)
		args           args
		expectedErr    error
	}{
		{
			name: "Test_SetPrivate_Success",
			args: args{
				noteId: noteId,
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), setPrivate, args.noteId).Return(nil, nil)

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},

			expectedErr: nil,
		},

		{
			name: "Test_SetPrivate_Fail_Error",
			args: args{
				noteId: noteId,
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), setPrivate, args.noteId).Return(nil, errors.New("err"))

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},

			expectedErr: errors.New("err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool, mockMetrics, tt.args)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			err := repo.SetPrivate(context.Background(), tt.args.noteId)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestNotePostgres_GetAttachList(t *testing.T) {
	noteId := uuid.NewV4()
	type args struct {
		noteId uuid.UUID
	}
	tests := []struct {
		args           args
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics, args)
		userId         uuid.UUID
		expectedErr    error
	}{
		{
			name: "GetTags_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, args args) {
				pgxRows := pgxpoolmock.NewRows([]string{"path"}).AddRow("path").ToPgxRows()

				mockPool.EXPECT().Query(gomock.Any(), getAttachList, args.noteId).Return(pgxRows, nil)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			args: args{
				noteId: noteId,
			},

			expectedErr: nil,
		},
		{
			name: "GetTags_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, args args) {
				pgxRows := pgxpoolmock.NewRows([]string{"path"}).ToPgxRows()

				mockPool.EXPECT().Query(gomock.Any(), getAttachList, args.noteId).Return(pgxRows, pgx.ErrNoRows)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},
			args: args{
				noteId: noteId,
			},
			expectedErr: pgx.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool, mockMetrics, tt.args)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			_, err := repo.GetAttachList(context.Background(), tt.args.noteId)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestNotePostgres_UpdateTagOnAllNotes(t *testing.T) {
	userId := uuid.NewV4()
	type args struct {
		userId  uuid.UUID
		oldName string
		newName string
	}
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics, pgx.Rows, args)
		args           args
		expectedErr    error
	}{
		{
			name: "Test_UpdateTag_Success",
			args: args{
				userId:  userId,
				oldName: "tag",
				newName: "tag1",
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), updateTagInAllNotes, args.oldName, args.newName, args.userId).Return(nil, nil)

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},

			expectedErr: nil,
		},

		{
			name: "Test_UpdateTag_Fail_Error",
			args: args{
				userId:  userId,
				oldName: "tag",
				newName: "tag1",
			},
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, args args) {
				mockPool.EXPECT().Exec(gomock.Any(), updateTagInAllNotes, args.oldName, args.newName, args.userId).Return(nil, errors.New("err"))

				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},

			expectedErr: errors.New("err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			pgxRows := pgxpoolmock.NewRows([]string{"tag_name", "user_id"}).AddRow("t", "0").ToPgxRows()

			tt.mockRepoAction(mockPool, mockMetrics, pgxRows, tt.args)

			repo := CreateNotePostgres(mockPool, mockMetrics)
			err := repo.UpdateTagOnAllNotes(context.Background(), tt.args.oldName, tt.args.newName, tt.args.userId)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
