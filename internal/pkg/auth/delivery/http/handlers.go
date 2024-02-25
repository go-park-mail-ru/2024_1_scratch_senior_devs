package http

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth"
)

type AuthHandler struct {
	uc auth.AuthUsecase
}

func NewAuthHandler(uc auth.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		uc: uc,
	}
}

func (h *AuthHandler) SignIn(ctx context.Context, w http.ResponseWriter, r http.Request) {

}
