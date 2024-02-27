package http

import (
	"fmt"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
	"net/http"
)

type AuthHandler struct {
	uc auth.AuthUsecase
}

func CreateAuthHandler(uc auth.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		uc: uc,
	}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	userData := models.UserFormData{}
	err := utils.GetRequestData(r, &userData)
	if err != nil {
		utils.WriteIncorrectFormatError(w)
		return
	}

	newUser, token, expTime, err := h.uc.SignUp(r.Context(), &userData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"this username is already taken"}`))
		return
	}

	http.SetCookie(w, utils.GenTokenCookie(token, expTime))
	w.Header().Set("Authorization", "Bearer "+token)

	err = utils.WriteResponseData(w, newUser, http.StatusCreated)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error in SignUp handler: %s", err)
		return
	}
}

func (h *AuthHandler) CheckUser(w http.ResponseWriter, r *http.Request) {
	jwtPayload, ok := r.Context().Value("payload").(models.JwtPayload)
	if !ok {
		fmt.Println("here")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userId := jwtPayload.Id
	currentUser, err := h.uc.CheckUser(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = utils.WriteResponseData(w, currentUser, http.StatusAccepted)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error in CheckUser handler: %s", err)
		return
	}
}

func (h *AuthHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, utils.DelTokenCookie())
	w.Header().Del("Authorization")
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	userData := models.UserFormData{}
	err := utils.GetRequestData(r, &userData)
	if err != nil {
		utils.WriteIncorrectFormatError(w)
		return
	}

	user, token, exp, err := h.uc.SignIn(r.Context(), &userData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"incorrect username or password"}`))
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	http.SetCookie(w, utils.GenTokenCookie(token, exp))

	err = utils.WriteResponseData(w, user, http.StatusAccepted)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error in SignIn handler: %s", err)
		return
	}
}
