package http

import (
	"bytes"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	mock_auth "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthHandler_SignUp(t *testing.T) {
	req := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewBufferString(`{"username":"test_user_2","password":"12345678a"}`))
	w := httptest.NewRecorder()

	ctrl := gomock.NewController(t)
	mockUsecase := mock_auth.NewMockAuthUsecase(ctrl)
	mockUsecase.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(&models.User{}, "", time.Now(), nil)
	defer ctrl.Finish()

	handler := CreateAuthHandler(mockUsecase)
	handler.SignUp(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
