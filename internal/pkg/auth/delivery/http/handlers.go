package http

import (
	"log/slog"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
)

type AuthHandler struct {
	uc     auth.AuthUsecase
	logger *slog.Logger
}

func CreateAuthHandler(uc auth.AuthUsecase, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		uc:     uc,
		logger: logger,
	}
}

// SignUp godoc
// @Summary		Sign up
// @Description	Add a new user to the database
// @Tags 		auth
// @ID			sign-up
// @Accept		json
// @Produce		json
// @Param		credentials body		models.UserFormData		true	"registration data"
// @Success		200			{object}	models.User				true	"user"
// @Failure		400			{object}	utils.ErrorResponse		true	"error"
// @Router		/api/auth/signup [post]
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", utils.GetRequestId(r.Context())), slog.String("func", utils.GFN()))

	userData := models.UserFormData{}
	if err := utils.GetRequestData(r, &userData); err != nil {
		utils.LogHandlerError(logger, http.StatusBadRequest, utils.ParseBodyError+err.Error())
		utils.WriteErrorMessage(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	if err := userData.Validate(); err != nil {
		utils.LogHandlerError(logger, http.StatusBadRequest, "validation error: "+err.Error())
		utils.WriteErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	newUser, token, expTime, err := h.uc.SignUp(r.Context(), userData)
	if err != nil {
		utils.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		utils.WriteErrorMessage(w, http.StatusBadRequest, "this username is already taken")
		return
	}

	http.SetCookie(w, utils.GenTokenCookie(token, expTime))
	w.Header().Set("Authorization", "Bearer "+token)

	if err := utils.WriteResponseData(w, newUser, http.StatusCreated); err != nil {
		utils.LogHandlerError(logger, http.StatusInternalServerError, utils.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.LogHandlerInfo(logger, http.StatusCreated, "success")
}

// CheckUser godoc
// @Summary		Check user
// @Description	Get user info if user is authorized
// @Tags 		auth
// @ID			check-user
// @Produce		json
// @Success		200		{object}	models.User		true	"user"
// @Failure		401
// @Router		/api/auth/check_user [get]
func (h *AuthHandler) CheckUser(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", utils.GetRequestId(r.Context())), slog.String("func", utils.GFN()))

	jwtPayload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	h.logger.Info(utils.GFN())
	if !ok {
		utils.LogHandlerError(logger, http.StatusUnauthorized, utils.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userId := jwtPayload.Id
	currentUser, err := h.uc.CheckUser(r.Context(), userId)
	if err != nil {
		utils.LogHandlerError(logger, http.StatusUnauthorized, err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := utils.WriteResponseData(w, currentUser, http.StatusOK); err != nil {
		utils.LogHandlerError(logger, http.StatusInternalServerError, utils.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.LogHandlerInfo(logger, http.StatusOK, "success")
}

// LogOut godoc
// @Summary		Log out
// @Description	Quit from user`s account
// @Tags 		auth
// @ID			log-out
// @Success		204
// @Router		/api/auth/logout [delete]
func (h *AuthHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", utils.GetRequestId(r.Context())), slog.String("func", utils.GFN()))

	http.SetCookie(w, utils.DelTokenCookie())
	w.Header().Del("Authorization")
	w.WriteHeader(http.StatusNoContent)

	utils.LogHandlerInfo(logger, http.StatusNoContent, "success")
}

// SignIn godoc
// @Summary		Sign in
// @Description	Login as a user
// @Tags 		auth
// @ID			sign-in
// @Accept		json
// @Produce		json
// @Param		credentials body		models.UserFormData		true	"login data"
// @Success		200			{object}	models.User				true	"user"
// @Failure		400			{object}	utils.ErrorResponse		true	"error"
// @Failure		401			{object}	utils.ErrorResponse		true	"error"
// @Router		/api/auth/login [post]
func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", utils.GetRequestId(r.Context())), slog.String("func", utils.GFN()))

	userData := models.UserFormData{}
	if err := utils.GetRequestData(r, &userData); err != nil {
		utils.LogHandlerError(logger, http.StatusBadRequest, utils.ParseBodyError+err.Error())
		utils.WriteErrorMessage(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	user, token, exp, err := h.uc.SignIn(r.Context(), userData)
	if err != nil {
		utils.LogHandlerError(logger, http.StatusUnauthorized, err.Error())
		utils.WriteErrorMessage(w, http.StatusUnauthorized, "incorrect username or password")
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	http.SetCookie(w, utils.GenTokenCookie(token, exp))

	if err := utils.WriteResponseData(w, user, http.StatusOK); err != nil {
		utils.LogHandlerError(logger, http.StatusInternalServerError, utils.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.LogHandlerInfo(logger, http.StatusNoContent, "success")
}
