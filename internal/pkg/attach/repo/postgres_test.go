package repo

import (
	"context"
	"errors"
	mock_metrics "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/metrics/mocks"
	"testing"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAttachRepo_GetAttach(t *testing.T) {
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics, pgx.Rows, uuid.UUID)
		Id             uuid.UUID
		columns        []string
		expectedErr    error
	}{
		{
			name: "GetAttach_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, id uuid.UUID) {
				mockPool.EXPECT().QueryRow(gomock.Any(), getAttach, id).Return(pgxRows)
				pgxRows.Next()
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			Id:          uuid.NewV4(),
			columns:     []string{"id", "path", "note_id"},
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			pgxRows := pgxpoolmock.NewRows(tt.columns).AddRow(uuid.NewV4(), "", uuid.NewV4()).ToPgxRows()

			tt.mockRepoAction(mockPool, mockMetrics, pgxRows, tt.Id)

			repo := CreateAttachRepo(mockPool, mockMetrics)
			_, err := repo.GetAttach(context.Background(), tt.Id)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestAttachRepo_AddAttach(t *testing.T) {
	attachId := uuid.NewV4()
	noteId := uuid.NewV4()
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics)
		err            error
	}{
		{
			name: "AddAttach_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
				mockPool.EXPECT().Exec(gomock.Any(), createAttach,
					attachId, "", noteId,
				).Return(nil, nil)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			err: nil,
		},
		{
			name: "AddAttach_",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
				mockPool.EXPECT().Exec(gomock.Any(), createAttach,
					attachId, "", noteId,
				).Return(nil, errors.New("db error"))
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},
			err: errors.New("db error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()
			tt.mockRepoAction(mockPool, mockMetrics)

			repo := CreateAttachRepo(mockPool, mockMetrics)
			err := repo.AddAttach(context.Background(), models.Attach{
				Id:     attachId,
				Path:   "",
				NoteId: noteId,
			})
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestAttachRepo_DeleteAttach(t *testing.T) {
	attachId := uuid.NewV4()

	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics)
		err            error
	}{
		{
			name: "DeleteAttach_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
				mockPool.EXPECT().Exec(gomock.Any(), deleteAttach,
					attachId,
				).Return(nil, nil)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			err: nil,
		},
		{
			name: "DeleteAttach_Fail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
				mockPool.EXPECT().Exec(gomock.Any(), deleteAttach,
					attachId,
				).Return(nil, errors.New("db error"))
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},
			err: errors.New("db error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()
			tt.mockRepoAction(mockPool, mockMetrics)

			repo := CreateAttachRepo(mockPool, mockMetrics)
			err := repo.DeleteAttach(context.Background(), attachId)
			assert.Equal(t, tt.err, err)
		})
	}
}
