package repo

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
)

var testLogger *slog.Logger

func init() {
	testLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

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
				mockPool.EXPECT().Query(gomock.Any(), getAllNotes, id, "%%", int64(1), int64(0)).Return(pgxRows, nil)
			},
			userId:      uuid.NewV4(),
			columns:     []string{"id", "data", "create_time", "update_time", "owner_id"},
			expectedErr: nil,
		},
		{
			name: "ReadAllNotes_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, pgxRows pgx.Rows, id uuid.UUID) {
				mockPool.EXPECT().Query(gomock.Any(), getAllNotes, id, "%%", int64(1), int64(0)).Return(pgxRows, pgx.ErrNoRows)
			},
			userId:      uuid.NewV4(),
			columns:     []string{"id", "data", "create_time", "update_time", "owner_id"},
			expectedErr: pgx.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			defer ctrl.Finish()

			pgxRows := pgxpoolmock.NewRows(tt.columns).AddRow(uuid.NewV4(), []byte{}, time.Now(), &time.Time{}, tt.userId).ToPgxRows()

			tt.mockRepoAction(mockPool, pgxRows, tt.userId)

			repo := CreateNoteRepo(mockPool, testLogger)
			_, err := repo.ReadAllNotes(context.Background(), tt.userId, int64(1), int64(0), "")

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
