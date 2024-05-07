package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNoteRepo_ReadAllNotes(t *testing.T) {
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, pgx.Rows, uuid.UUID)
		userId         uuid.UUID
		columns        []string
		expectedErr    error
	}{
		{
			name: "ReadAllNotes_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, pgxRows pgx.Rows, id uuid.UUID) {
				mockPool.EXPECT().Query(gomock.Any(), getAllNotes, id, int64(1), int64(0), []string{}).Return(pgxRows, nil)
			},
			userId:      uuid.NewV4(),
			columns:     []string{"id", "data", "create_time", "update_time", "owner_id", "parent", "children", "tags", "collaborators"},
			expectedErr: nil,
		},
		{
			name: "ReadAllNotes_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, pgxRows pgx.Rows, id uuid.UUID) {
				mockPool.EXPECT().Query(gomock.Any(), getAllNotes, id, int64(1), int64(0), []string{}).Return(pgxRows, pgx.ErrNoRows)
			},
			userId:      uuid.NewV4(),
			columns:     []string{"id", "data", "create_time", "update_time", "owner_id", "parent", "children", "tags", "collaborators"},
			expectedErr: pgx.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			defer ctrl.Finish()

			pgxRows := pgxpoolmock.NewRows(tt.columns).AddRow(uuid.NewV4(), []byte{}, time.Now(), time.Time{}, tt.userId, uuid.UUID{}, []uuid.UUID{}, []string{}, []uuid.UUID{}).ToPgxRows()

			tt.mockRepoAction(mockPool, pgxRows, tt.userId)

			repo := CreateNotePostgres(mockPool)
			_, err := repo.ReadAllNotes(context.Background(), tt.userId, int64(1), int64(0), []string{})

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestNoteRepo_ReadNote(t *testing.T) {
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, pgx.Rows, uuid.UUID)
		Id             uuid.UUID
		columns        []string
		expectedErr    error
	}{
		{
			name: "ReadNote_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, pgxRows pgx.Rows, id uuid.UUID) {
				mockPool.EXPECT().QueryRow(gomock.Any(), getNote, id).Return(pgxRows)
				pgxRows.Next()
			},
			Id:          uuid.NewV4(),
			columns:     []string{"id", "data", "create_time", "update_time", "owner_id", "parent", "children", "tags", "collaborators"},
			expectedErr: nil,
		},
		{
			name: "ReadNote_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, pgxRows pgx.Rows, id uuid.UUID) {
				mockPool.EXPECT().QueryRow(gomock.Any(), getNote, id).Return(pgxRows)
				pgxRows.Next()
			},
			Id:          uuid.NewV4(),
			columns:     []string{"id", "data", "create_time", "update_time", "owner_id", "parent", "children", "tags", "collaborators"},
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			defer ctrl.Finish()

			pgxRows := pgxpoolmock.NewRows(tt.columns).AddRow(uuid.NewV4(), []byte{}, time.Now(), time.Time{}, tt.Id, uuid.UUID{}, []uuid.UUID{}, []string{}, []uuid.UUID{}).ToPgxRows()

			tt.mockRepoAction(mockPool, pgxRows, tt.Id)

			repo := CreateNotePostgres(mockPool)
			_, err := repo.ReadNote(context.Background(), tt.Id)

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
		mockRepoAction func(*pgxpoolmock.MockPgxPool)
		err            error
	}{
		{
			name: "CreateNote_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(gomock.Any(), createNote,
					Id, []byte{}, currTime, currTime, userId, Id, []uuid.UUID{}, []string{}, []uuid.UUID{},
				).Return(nil, nil)
			},
			err: nil,
		},
		{
			name: "CreateNote_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(gomock.Any(), createNote,
					Id, []byte{}, currTime, currTime, userId, Id, []uuid.UUID{}, []string{}, []uuid.UUID{},
				).Return(nil, errors.New("err"))
			},
			err: errors.New("err"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool)

			repo := CreateNotePostgres(mockPool)
			err := repo.CreateNote(context.Background(), models.Note{
				Id:            Id,
				Data:          []byte{},
				CreateTime:    currTime,
				UpdateTime:    currTime,
				OwnerId:       userId,
				Parent:        Id,
				Children:      []uuid.UUID{},
				Tags:          []string{},
				Collaborators: []uuid.UUID{},
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
		mockRepoAction func(*pgxpoolmock.MockPgxPool)
		err            error
	}{
		{
			name: "UpdateNote_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(gomock.Any(), updateNote,
					[]byte{}, currTime, Id,
				).Return(nil, nil)
			},
			err: nil,
		},
		{
			name: "UpdateNote_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(gomock.Any(), updateNote,
					[]byte{}, currTime, Id,
				).Return(nil, pgx.ErrNoRows)
			},
			err: pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool)

			repo := CreateNotePostgres(mockPool)
			err := repo.UpdateNote(context.Background(), models.Note{
				Id:            Id,
				Data:          []byte{},
				CreateTime:    currTime,
				UpdateTime:    currTime,
				OwnerId:       userId,
				Parent:        Id,
				Children:      []uuid.UUID{},
				Tags:          []string{},
				Collaborators: []uuid.UUID{},
			})

			assert.Equal(t, tt.err, err)
		})
	}
}

func TestNoteRepo_DeleteNote(t *testing.T) {
	Id := uuid.NewV4()

	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool)
		err            error
	}{
		{
			name: "DeleteNote_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(gomock.Any(), deleteNote, Id).Return(nil, nil)
			},
			err: nil,
		},
		{
			name: "DeleteNote_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(gomock.Any(), deleteNote, Id).Return(nil, pgx.ErrNoRows)
			},
			err: pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool)

			repo := CreateNotePostgres(mockPool)
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
		mockRepoAction func(*pgxpoolmock.MockPgxPool)
		err            error
	}{
		{
			name: "AddSubNote_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(gomock.Any(), addSubNote, parentId, noteId).Return(nil, nil)
			},
			err: nil,
		},
		{
			name: "AddSubNote_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(gomock.Any(), addSubNote, parentId, noteId).Return(nil, pgx.ErrNoRows)
			},
			err: pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool)

			repo := CreateNotePostgres(mockPool)
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
		mockRepoAction func(*pgxpoolmock.MockPgxPool)
		err            error
	}{
		{
			name: "RemoveSubNote_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(gomock.Any(), removeSubNote, noteId, parentId).Return(nil, nil)
			},
			err: nil,
		},
		{
			name: "RemoveSubNote_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(gomock.Any(), removeSubNote, noteId, parentId).Return(nil, pgx.ErrNoRows)
			},
			err: pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool)

			repo := CreateNotePostgres(mockPool)
			err := repo.RemoveSubNote(context.Background(), parentId, noteId)

			assert.Equal(t, tt.err, err)
		})
	}
}

func TestNotePostgres_AddTag(t *testing.T) {
	noteId := uuid.NewV4()

	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool)
		err            error
	}{
		{
			name: "AddTag_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(gomock.Any(), addTag, "tag", noteId).Return(nil, nil)
			},
			err: nil,
		},
		{
			name: "AddTag_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(gomock.Any(), addTag, "tag", noteId).Return(nil, pgx.ErrNoRows)
			},
			err: pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool)

			repo := CreateNotePostgres(mockPool)
			err := repo.AddTag(context.Background(), "tag", noteId)

			assert.Equal(t, tt.err, err)
		})
	}
}

func TestNotePostgres_DeleteTag(t *testing.T) {
	noteId := uuid.NewV4()

	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool)
		err            error
	}{
		{
			name: "DeleteTag_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(gomock.Any(), deleteTag, "tag", noteId).Return(nil, nil)
			},
			err: nil,
		},
		{
			name: "DeleteTag_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(gomock.Any(), deleteTag, "tag", noteId).Return(nil, pgx.ErrNoRows)
			},
			err: pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool)

			repo := CreateNotePostgres(mockPool)
			err := repo.DeleteTag(context.Background(), "tag", noteId)

			assert.Equal(t, tt.err, err)
		})
	}
}

func TestNotePostgres_AddCollaborator(t *testing.T) {
	noteId := uuid.NewV4()
	guestId := uuid.NewV4()

	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool)
		err            error
	}{
		{
			name: "AddCollaborator_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(gomock.Any(), addCollaborator, guestId, noteId).Return(nil, nil)
			},
			err: nil,
		},
		{
			name: "AddCollaborator_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(gomock.Any(), addCollaborator, guestId, noteId).Return(nil, pgx.ErrNoRows)
			},
			err: pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool)

			repo := CreateNotePostgres(mockPool)
			err := repo.AddCollaborator(context.Background(), noteId, guestId)

			assert.Equal(t, tt.err, err)
		})
	}
}

func TestNotePostgres_GetTags(t *testing.T) {
	userId := uuid.NewV4()
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, pgx.Rows, uuid.UUID)
		userId         uuid.UUID
		expectedErr    error
	}{
		{
			name: "GetTags_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, pgxRows pgx.Rows, id uuid.UUID) {
				mockPool.EXPECT().Query(gomock.Any(), getTags, userId).Return(pgxRows, nil)
			},
			userId:      userId,
			expectedErr: nil,
		},
		{
			name: "GetTags_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, pgxRows pgx.Rows, id uuid.UUID) {
				mockPool.EXPECT().Query(gomock.Any(), getTags, userId).Return(pgxRows, pgx.ErrNoRows)
			},
			userId:      userId,
			expectedErr: pgx.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			defer ctrl.Finish()

			pgxRows := pgxpoolmock.NewRows([]string{"tags"}).AddRow("tag").ToPgxRows()

			tt.mockRepoAction(mockPool, pgxRows, tt.userId)

			repo := CreateNotePostgres(mockPool)
			_, err := repo.GetTags(context.Background(), tt.userId)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
