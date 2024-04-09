package repo

import (
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

var testLogger *slog.Logger
var testConfig *config.Config

func init() {
	testLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	testConfig = config.LoadConfig("../../config/config.yaml", testLogger)
}

func TestNoteRepo_MakeHelloNoteData(t *testing.T) {
	username := "testuser"
	expected := []byte(`
		{
			"title": "You-note ❤️",
			"content": [
				{
					"id": "1",
					"type": "div",
					"content": [
						{
							"id": "2",
							"content": "Привет, testuser!"
						}
					]
				}
			]
		}
	`)

	t.Run("Test_MakeHelloNoteData", func(t *testing.T) {
		repo := CreateNoteRepo(&elastic.Client{}, testLogger, testConfig.Elastic)
		result := repo.MakeHelloNoteData(username)

		assert.Equal(t, expected, result)
	})
}

//func TestNoteRepo_ReadAllNotes(t *testing.T) {
//	tests := []struct {
//		name           string
//		mockRepoAction func(*pgxpoolmock.MockPgxPool, pgx.Rows, uuid.UUID)
//		userId         uuid.UUID
//		columns        []string
//		expectedErr    error
//	}{
//		{
//			name: "ReadAllNotes_Success",
//			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, pgxRows pgx.Rows, id uuid.UUID) {
//				mockPool.EXPECT().Query(gomock.Any(), getAllNotes, id, "%%", int64(1), int64(0)).Return(pgxRows, nil)
//			},
//			userId:      uuid.NewV4(),
//			columns:     []string{"id", "data", "create_time", "update_time", "owner_id"},
//			expectedErr: nil,
//		},
//		{
//			name: "ReadAllNotes_Fail",
//			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, pgxRows pgx.Rows, id uuid.UUID) {
//				mockPool.EXPECT().Query(gomock.Any(), getAllNotes, id, "%%", int64(1), int64(0)).Return(pgxRows, pgx.ErrNoRows)
//			},
//			userId:      uuid.NewV4(),
//			columns:     []string{"id", "data", "create_time", "update_time", "owner_id"},
//			expectedErr: pgx.ErrNoRows,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
//			defer ctrl.Finish()
//
//			pgxRows := pgxpoolmock.NewRows(tt.columns).AddRow(uuid.NewV4(), []byte{}, time.Now(), time.Time{}, tt.userId).ToPgxRows()
//
//			tt.mockRepoAction(mockPool, pgxRows, tt.userId)
//
//			repo := CreateNoteRepo(mockPool, testLogger)
//			_, err := repo.ReadAllNotes(context.Background(), tt.userId, int64(1), int64(0), "")
//
//			assert.Equal(t, tt.expectedErr, err)
//		})
//	}
//}
//
//func TestNoteRepo_ReadNote(t *testing.T) {
//	tests := []struct {
//		name           string
//		mockRepoAction func(*pgxpoolmock.MockPgxPool, pgx.Rows, uuid.UUID)
//		Id             uuid.UUID
//		columns        []string
//		expectedErr    error
//	}{
//		{
//			name: "ReadNote_Success",
//			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, pgxRows pgx.Rows, id uuid.UUID) {
//				mockPool.EXPECT().QueryRow(gomock.Any(), getNote, id).Return(pgxRows)
//				pgxRows.Next()
//			},
//			Id:          uuid.NewV4(),
//			columns:     []string{"id", "data", "create_time", "update_time", "owner_id"},
//			expectedErr: nil,
//		},
//		{
//			name: "ReadNote_Fail",
//			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, pgxRows pgx.Rows, id uuid.UUID) {
//				mockPool.EXPECT().QueryRow(gomock.Any(), getNote, id).Return(pgxRows)
//				pgxRows.Next()
//			},
//			Id:          uuid.NewV4(),
//			columns:     []string{"id", "data", "create_time", "update_time", "owner_id"},
//			expectedErr: nil,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
//			defer ctrl.Finish()
//
//			pgxRows := pgxpoolmock.NewRows(tt.columns).AddRow(uuid.NewV4(), []byte{}, time.Now(), time.Time{}, tt.Id).ToPgxRows()
//
//			tt.mockRepoAction(mockPool, pgxRows, tt.Id)
//
//			repo := CreateNoteRepo(mockPool, testLogger)
//			_, err := repo.ReadNote(context.Background(), tt.Id)
//
//			assert.Equal(t, tt.expectedErr, err)
//		})
//	}
//}
//
//func TestNoteRepo_CreateNote(t *testing.T) {
//	Id := uuid.NewV4()
//	userId := uuid.NewV4()
//	currTime := time.Now().UTC()
//
//	tests := []struct {
//		name           string
//		mockRepoAction func(*pgxpoolmock.MockPgxPool)
//		err            error
//	}{
//		{
//			name: "CreateNote_Success",
//			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
//				mockPool.EXPECT().Exec(gomock.Any(), createNote,
//					Id, []byte{}, currTime, currTime, userId,
//				).Return(nil, nil)
//			},
//			err: nil,
//		},
//		{
//			name: "CreateNote_Fail",
//			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
//				mockPool.EXPECT().Exec(gomock.Any(), createNote,
//					Id, []byte{}, currTime, currTime, userId,
//				).Return(nil, errors.New("err"))
//			},
//			err: errors.New("err"),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
//			defer ctrl.Finish()
//
//			tt.mockRepoAction(mockPool)
//
//			repo := CreateNoteRepo(mockPool, testLogger)
//			err := repo.CreateNote(context.Background(), models.Note{
//				Id:         Id,
//				Data:       []byte{},
//				CreateTime: currTime,
//				UpdateTime: currTime,
//				OwnerId:    userId,
//			})
//
//			assert.Equal(t, tt.err, err)
//		})
//	}
//}
//
//func TestNoteRepo_UpdateNote(t *testing.T) {
//	Id := uuid.NewV4()
//	userId := uuid.NewV4()
//	currTime := time.Now().UTC()
//
//	tests := []struct {
//		name           string
//		mockRepoAction func(*pgxpoolmock.MockPgxPool)
//		err            error
//	}{
//		{
//			name: "UpdateNote_Success",
//			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
//				mockPool.EXPECT().Exec(gomock.Any(), updateNote,
//					[]byte{}, currTime, Id,
//				).Return(nil, nil)
//			},
//			err: nil,
//		},
//		{
//			name: "UpdateNote_Fail",
//			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
//				mockPool.EXPECT().Exec(gomock.Any(), updateNote,
//					[]byte{}, currTime, Id,
//				).Return(nil, pgx.ErrNoRows)
//			},
//			err: pgx.ErrNoRows,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
//			defer ctrl.Finish()
//
//			tt.mockRepoAction(mockPool)
//
//			repo := CreateNoteRepo(mockPool, testLogger)
//			err := repo.UpdateNote(context.Background(), models.Note{
//				Id:         Id,
//				Data:       []byte{},
//				CreateTime: currTime,
//				UpdateTime: currTime,
//				OwnerId:    userId,
//			})
//
//			assert.Equal(t, tt.err, err)
//		})
//	}
//}
//
//func TestNoteRepo_DeleteNote(t *testing.T) {
//	Id := uuid.NewV4()
//
//	tests := []struct {
//		name           string
//		mockRepoAction func(*pgxpoolmock.MockPgxPool)
//		err            error
//	}{
//		{
//			name: "DeleteNote_Success",
//			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
//				mockPool.EXPECT().Exec(gomock.Any(), deleteNote, Id).Return(nil, nil)
//			},
//			err: nil,
//		},
//		{
//			name: "DeleteNote_Fail",
//			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
//				mockPool.EXPECT().Exec(gomock.Any(), deleteNote, Id).Return(nil, pgx.ErrNoRows)
//			},
//			err: pgx.ErrNoRows,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
//			defer ctrl.Finish()
//
//			tt.mockRepoAction(mockPool)
//
//			repo := CreateNoteRepo(mockPool, testLogger)
//			err := repo.DeleteNote(context.Background(), Id)
//
//			assert.Equal(t, tt.err, err)
//		})
//	}
//}
