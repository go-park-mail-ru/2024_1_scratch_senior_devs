package http

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	mock_attach "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/attach/mocks"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
)

var testLogger *slog.Logger
var testConfig *config.Config

func init() {
	testLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	testConfig = config.LoadConfig("../../../config/config.yaml", testLogger)
}

const (
	testNameUnauthorized = "Test_Unauthorized"
	testNameBadRequest   = "Test_Bad_Request"
)

func TestAttachHandler_DeleteAttach(t *testing.T) {

	tests := []struct {
		name           string
		ucMocker       func(ctx context.Context, uc *mock_attach.MockAttachUsecase, attachID uuid.UUID, userID uuid.UUID)
		expectedStatus int
		username       string
		attachID       uuid.UUID
		userID         uuid.UUID
	}{{
		name: "Test_Success",
		ucMocker: func(ctx context.Context, uc *mock_attach.MockAttachUsecase, attachID uuid.UUID, userID uuid.UUID) {
			uc.EXPECT().DeleteAttach(ctx, attachID, userID).Return(nil)
		},
		expectedStatus: http.StatusNoContent,
		username:       "alla",
		attachID:       uuid.FromStringOrNil("ac6966bc-3c26-45a0-963e-b168fc34fd79"),
		userID:         uuid.FromStringOrNil("ac5566bc-3c26-45a0-963e-b168fc34fd79"),
	},
		{
			name: "Test_Fail_NotFound",
			ucMocker: func(ctx context.Context, uc *mock_attach.MockAttachUsecase, attachID uuid.UUID, userID uuid.UUID) {
				uc.EXPECT().DeleteAttach(ctx, attachID, userID).Return(errors.New("uc error"))
			},
			expectedStatus: http.StatusNotFound,
			username:       "alla",
			attachID:       uuid.FromStringOrNil("ac6966bc-3c26-45a0-963e-b168fc34fd79"),
			userID:         uuid.FromStringOrNil("ac5566bc-3c26-45a0-963e-b168fc34fd79"),
		},
		{
			name: testNameUnauthorized,
			ucMocker: func(ctx context.Context, uc *mock_attach.MockAttachUsecase, attachID uuid.UUID, userID uuid.UUID) {
			},
			expectedStatus: http.StatusUnauthorized,
			username:       "alla",
			attachID:       uuid.FromStringOrNil("ac6966bc-3c26-45a0-963e-b168fc34fd79"),
			userID:         uuid.FromStringOrNil("ac5566bc-3c26-45a0-963e-b168fc34fd79"),
		},
		{
			name: testNameBadRequest,
			ucMocker: func(ctx context.Context, uc *mock_attach.MockAttachUsecase, attachID uuid.UUID, userID uuid.UUID) {
			},
			expectedStatus: http.StatusBadRequest,
			username:       "alla",
			attachID:       uuid.FromStringOrNil(""),
			userID:         uuid.FromStringOrNil("ac5566bc-3c26-45a0-963e-b168fc34fd79"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			uc := mock_attach.NewMockAttachUsecase(ctrl)
			defer ctrl.Finish()
			req := httptest.NewRequest("POST", "http://example.com/api/handler/", bytes.NewBufferString(""))
			w := httptest.NewRecorder()
			ctx := context.WithValue(req.Context(), config.PayloadContextKey, models.JwtPayload{Id: tt.userID, Username: tt.username})
			req = req.WithContext(ctx)
			if tt.name == testNameUnauthorized {
				req = req.WithContext(context.Background())
			}
			if tt.name != testNameBadRequest {
				req = mux.SetURLVars(req, map[string]string{"id": tt.attachID.String()})
			}

			tt.ucMocker(req.Context(), uc, tt.attachID, tt.userID)

			h := CreateAttachHandler(uc, testLogger, testConfig.Attach)
			h.DeleteAttach(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAttachHandler_GetAttach(t *testing.T) {

	tests := []struct {
		name           string
		ucMocker       func(ctx context.Context, uc *mock_attach.MockAttachUsecase, attachID uuid.UUID, userID uuid.UUID)
		expectedStatus int
		username       string
		attachID       uuid.UUID
		userID         uuid.UUID
	}{{
		name: "Test_Fail_NotFound_File",
		ucMocker: func(ctx context.Context, uc *mock_attach.MockAttachUsecase, attachID uuid.UUID, userID uuid.UUID) {
			uc.EXPECT().GetAttach(ctx, attachID, userID).Return(models.Attach{}, nil)
		},
		expectedStatus: http.StatusNotFound,
		username:       "alla",
		attachID:       uuid.FromStringOrNil("ac6966bc-3c26-45a0-963e-b168fc34fd79"),
		userID:         uuid.FromStringOrNil("ac5566bc-3c26-45a0-963e-b168fc34fd79"),
	},
		{
			name: "Test_Fail_NotFound_GetAttach",
			ucMocker: func(ctx context.Context, uc *mock_attach.MockAttachUsecase, attachID uuid.UUID, userID uuid.UUID) {
				uc.EXPECT().GetAttach(ctx, attachID, userID).Return(models.Attach{}, errors.New("error"))
			},
			expectedStatus: http.StatusNotFound,
			username:       "alla",
			attachID:       uuid.FromStringOrNil("ac6966bc-3c26-45a0-963e-b168fc34fd79"),
			userID:         uuid.FromStringOrNil("ac5566bc-3c26-45a0-963e-b168fc34fd79"),
		},

		{
			name: testNameUnauthorized,
			ucMocker: func(ctx context.Context, uc *mock_attach.MockAttachUsecase, attachID uuid.UUID, userID uuid.UUID) {
			},
			expectedStatus: http.StatusUnauthorized,
			username:       "alla",
			attachID:       uuid.FromStringOrNil("ac6966bc-3c26-45a0-963e-b168fc34fd79"),
			userID:         uuid.FromStringOrNil("ac5566bc-3c26-45a0-963e-b168fc34fd79"),
		},
		{
			name: testNameBadRequest,
			ucMocker: func(ctx context.Context, uc *mock_attach.MockAttachUsecase, attachID uuid.UUID, userID uuid.UUID) {
			},
			expectedStatus: http.StatusBadRequest,
			username:       "alla",
			attachID:       uuid.FromStringOrNil(""),
			userID:         uuid.FromStringOrNil("ac5566bc-3c26-45a0-963e-b168fc34fd79"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			uc := mock_attach.NewMockAttachUsecase(ctrl)
			defer ctrl.Finish()
			req := httptest.NewRequest("POST", "http://example.com/api/handler/", bytes.NewBufferString(""))
			w := httptest.NewRecorder()
			ctx := context.WithValue(req.Context(), config.PayloadContextKey, models.JwtPayload{Id: tt.userID, Username: tt.username})
			req = req.WithContext(ctx)
			if tt.name == testNameUnauthorized {
				req = req.WithContext(context.Background())
			}
			if tt.name != testNameBadRequest {
				req = mux.SetURLVars(req, map[string]string{"id": tt.attachID.String()})
			}

			tt.ucMocker(req.Context(), uc, tt.attachID, tt.userID)

			h := CreateAttachHandler(uc, testLogger, testConfig.Attach)
			h.GetAttach(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAttachHandler_AddAttach(t *testing.T) {

	tests := []struct {
		name           string
		ucMocker       func(ctx context.Context, uc *mock_attach.MockAttachUsecase, attachID uuid.UUID, userID uuid.UUID)
		expectedStatus int
		username       string
		attachID       uuid.UUID
		userID         uuid.UUID
	}{{
		name: "Test_Fail_MultipartProblem",
		ucMocker: func(ctx context.Context, uc *mock_attach.MockAttachUsecase, attachID uuid.UUID, userID uuid.UUID) {

		},
		expectedStatus: http.StatusRequestEntityTooLarge,
		username:       "alla",
		attachID:       uuid.FromStringOrNil("ac6966bc-3c26-45a0-963e-b168fc34fd79"),
		userID:         uuid.FromStringOrNil("ac5566bc-3c26-45a0-963e-b168fc34fd79"),
	},

		{
			name: testNameUnauthorized,
			ucMocker: func(ctx context.Context, uc *mock_attach.MockAttachUsecase, attachID uuid.UUID, userID uuid.UUID) {
			},
			expectedStatus: http.StatusUnauthorized,
			username:       "alla",
			attachID:       uuid.FromStringOrNil("ac6966bc-3c26-45a0-963e-b168fc34fd79"),
			userID:         uuid.FromStringOrNil("ac5566bc-3c26-45a0-963e-b168fc34fd79"),
		},
		{
			name: testNameBadRequest,
			ucMocker: func(ctx context.Context, uc *mock_attach.MockAttachUsecase, attachID uuid.UUID, userID uuid.UUID) {
			},
			expectedStatus: http.StatusBadRequest,
			username:       "alla",
			attachID:       uuid.FromStringOrNil(""),
			userID:         uuid.FromStringOrNil("ac5566bc-3c26-45a0-963e-b168fc34fd79"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			uc := mock_attach.NewMockAttachUsecase(ctrl)
			defer ctrl.Finish()
			body := new(bytes.Buffer)
			req := httptest.NewRequest("POST", "http://example.com/api/handler/", body)
			ctx := context.WithValue(req.Context(), config.PayloadContextKey, models.JwtPayload{Id: tt.userID, Username: tt.username})

			req = req.WithContext(ctx)
			req.Header.Set("Content-Type", "multipart/form-data; boundary=---")
			if tt.name == testNameUnauthorized {
				req = req.WithContext(context.Background())
			}
			if tt.name != testNameBadRequest {
				req = mux.SetURLVars(req, map[string]string{"id": tt.attachID.String()})
			}

			w := httptest.NewRecorder()

			tt.ucMocker(req.Context(), uc, tt.attachID, tt.userID)

			h := CreateAttachHandler(uc, testLogger, testConfig.Attach)
			h.AddAttach(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
