package http

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
)

type AuthHandler struct {
	uc auth.AuthUsecase
}

func CreateAuthHandler(uc auth.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		uc: uc,
	}
}

// SignUp
// @Summary		Sign up
// @Description	Add a new user to the database
// @Tags 		auth
// @ID			sign-up
// @Accept		json
// @Produce		json
// @Param		username	body		string					true	"username"
// @Param		password	body		string					true	"password"
// @Success		200			{object}	models.User				true	"user"
// @Failure		400			{object}	utils.ErrorResponse		true	"error"
// @Router		/api/auth/signup [post]
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	userData := models.UserFormData{}
	err := utils.GetRequestData(r, &userData)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	if err := userData.Validate(); err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	newUser, token, expTime, err := h.uc.SignUp(r.Context(), userData)
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

// CheckUser
// @Summary		Check user
// @Description	Get user info if user is authorized
// @Tags 		auth
// @ID			check-user
// @Produce		json
// @Success		200		{object}	models.User		true	"user"
// @Failure		401
// @Router		/api/auth/check_user [get]
func (h *AuthHandler) CheckUser(w http.ResponseWriter, r *http.Request) {
	jwtPayload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userId := jwtPayload.Id
	currentUser, err := h.uc.CheckUser(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = utils.WriteResponseData(w, currentUser, http.StatusOK)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error in CheckUser handler: %s", err)
		return
	}
}

// LogOut
// @Summary		Log out
// @Description	Quit from user`s account
// @Tags 		auth
// @ID			log-out
// @Success		204
// @Router		/api/auth/logout [delete]
func (h *AuthHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, utils.DelTokenCookie())
	w.Header().Del("Authorization")
	w.WriteHeader(http.StatusNoContent)
}

// SignIn
// @Summary		Sign in
// @Description	Login as a user
// @Tags 		auth
// @ID			sign-in
// @Accept		json
// @Produce		json
// @Param		username	body		string					true	"username"
// @Param		password	body		string					true	"password"
// @Success		200			{object}	models.User				true	"user"
// @Failure		400			{object}	utils.ErrorResponse		true	"error"
// @Router		/api/auth/login [post]
func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	userData := models.UserFormData{}
	err := utils.GetRequestData(r, &userData)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	user, token, exp, err := h.uc.SignIn(r.Context(), userData)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusUnauthorized, "incorrect username or password")
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	http.SetCookie(w, utils.GenTokenCookie(token, exp))

	err = utils.WriteResponseData(w, user, http.StatusOK)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error in SignIn handler: %s", err)
		return
	}
}
