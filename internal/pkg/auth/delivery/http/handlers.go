package http

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
)

type AuthHandler struct {
	uc auth.AuthUsecase
}

func NewAuthHandler(uc auth.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		uc: uc,
	}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	userData := models.UserFormData{}
	err := utils.GetRequestData(r, &userData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect user data format"))
		return
	}

	newUser, token, expTime, err := h.uc.SignUp(r.Context(), &userData)

	http.SetCookie(w, utils.GenTokenCookie(token, expTime))

	err = utils.WriteResponseData(w, newUser)
	w.Header().Set("Authorization", token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r http.Request) {
	user := models.UserFormData{}
	utils.GetRequestData(w, &r, models.UserFormData)
}

func (h *AuthHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, utils.GenTokenCookie("", time.Now()))
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) CheckUser(w http.ResponseWriter, r http.Request) {
	userId := r.Context().Value("payload").(models.JwtPayload).Id
	currentUser, err := h.uc.CheckUser(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	utils.WriteResponseData(w, currentUser)
	w.WriteHeader(http.StatusOK)
}
