package http

import (
	"encoding/base32"
	"errors"
	"fmt"
	"github.com/dgryski/dgoogauth"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/cookie"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/images"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/request"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/response"
	"github.com/skip2/go-qrcode"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

const (
	maxFormDataSize = 5 * 1024 * 1024
	qrIssuer        = "YouNoteQR"
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
// @Failure		400			{object}	response.ErrorResponse	true	"error"
// @Router		/api/auth/signup [post]
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

	userData := models.UserFormData{}
	if err := request.GetRequestData(r, &userData); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, response.ParseBodyError+err.Error())
		response.WriteErrorMessage(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	if err := userData.Validate(); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, "validation error: "+err.Error())
		response.WriteErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	newUser, token, expTime, err := h.uc.SignUp(r.Context(), userData)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		response.WriteErrorMessage(w, http.StatusBadRequest, "this username is already taken")
		return
	}

	http.SetCookie(w, cookie.GenTokenCookie(token, expTime))
	w.Header().Set("Authorization", "Bearer "+token)

	if err := response.WriteResponseData(w, newUser, http.StatusCreated); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, response.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusCreated, "success")
}

// CheckUser godoc
// @Summary		Check user
// @Description	Get user info if user is authorized
// @Tags 		auth
// @ID			check-user
// @Success		200
// @Failure		401
// @Router		/api/auth/check_user [get]
func (h *AuthHandler) CheckUser(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

	_, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, response.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
	w.WriteHeader(http.StatusOK)
}

// LogOut godoc
// @Summary		Log out
// @Description	Quit from user`s account
// @Tags 		auth
// @ID			log-out
// @Success		204
// @Router		/api/auth/logout [delete]
func (h *AuthHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

	http.SetCookie(w, cookie.DelTokenCookie())
	w.Header().Del("Authorization")
	w.WriteHeader(http.StatusNoContent)

	log.LogHandlerInfo(logger, http.StatusNoContent, "success")
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
// @Failure		400			{object}	response.ErrorResponse	true	"error"
// @Failure		401			{object}	response.ErrorResponse	true	"error"
// @Router		/api/auth/login [post]
func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

	userData := models.UserFormData{}
	if err := request.GetRequestData(r, &userData); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, response.ParseBodyError+err.Error())
		response.WriteErrorMessage(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	user, token, exp, err := h.uc.SignIn(r.Context(), userData)
	if err != nil {
		log.LogHandlerError(logger, http.StatusUnauthorized, err.Error())
		response.WriteErrorMessage(w, http.StatusUnauthorized, "incorrect username or password")
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	http.SetCookie(w, cookie.GenTokenCookie(token, exp))

	if err := response.WriteResponseData(w, user, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, response.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusNoContent, "success")
}

// GetProfile godoc
// @Summary		Get profile
// @Description	Get user info if user is authorized
// @Tags 		profile
// @ID			get-profile
// @Produce		json
// @Success		200		{object}	models.User		true	"user"
// @Failure		401
// @Router		/api/profile/get [get]
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

	jwtPayload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, response.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userId := jwtPayload.Id
	currentUser, err := h.uc.CheckUser(r.Context(), userId)
	if err != nil {
		log.LogHandlerError(logger, http.StatusUnauthorized, err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := response.WriteResponseData(w, currentUser, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, response.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

// UpdateProfile godoc
// @Summary		Update profile
// @Description	Change password and/or description
// @Tags 		profile
// @ID			update-profile
// @Accept		json
// @Produce		json
// @Param		credentials body		models.ProfileUpdatePayload		true	"update data"
// @Success		200			{object}	models.User						true	"user"
// @Failure		400			{object}	response.ErrorResponse			true	"error"
// @Router		/api/profile/update [post]
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

	jwtPayload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, response.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	payload := models.ProfileUpdatePayload{}
	if err := request.GetRequestData(r, &payload); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, response.ParseBodyError+err.Error())
		response.WriteErrorMessage(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	user, err := h.uc.UpdateProfile(r.Context(), jwtPayload.Id, payload)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := response.WriteResponseData(w, user, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, response.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

// UpdateProfileAvatar godoc
// @Summary		Update profile images
// @Description	Change images
// @Tags 		profile
// @ID			update-profile-images
// @Accept		multipart/form-data
// @Produce		json
// @Param		avatar 		formData	file						true	"avatar"
// @Success		200			{object}	models.User					true	"user"
// @Failure		400			{object}	response.ErrorResponse		true	"error"
// @Failure		413
// @Router		/api/profile/update_avatar [post]
func (h *AuthHandler) UpdateProfileAvatar(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

	jwtPayload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, response.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxFormDataSize)
	defer r.Body.Close()

	err := r.ParseMultipartForm(maxFormDataSize)
	if err != nil {
		log.LogHandlerError(logger, http.StatusRequestEntityTooLarge, err.Error())
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	defer func() {
		if err := r.MultipartForm.RemoveAll(); err != nil {
			logger.Error(err.Error())
		}
	}()

	avatar, _, err := r.FormFile("avatar")
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		response.WriteErrorMessage(w, http.StatusBadRequest, "no file found")
		return
	}
	content, err := io.ReadAll(avatar)
	if err != nil && !errors.Is(err, io.EOF) {
		if errors.As(err, new(*http.MaxBytesError)) {
			log.LogHandlerError(logger, http.StatusRequestEntityTooLarge, err.Error())
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}
	}

	fileExtension := images.CheckFileFormat(content)
	if fileExtension == "" {
		log.LogHandlerError(logger, http.StatusBadRequest, "incorrect file format")
		response.WriteErrorMessage(w, http.StatusBadRequest, "incorrect file format")
		return
	}

	user, err := h.uc.UpdateProfileAvatar(r.Context(), jwtPayload.Id, avatar, fileExtension)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := response.WriteResponseData(w, user, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, response.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

func (h *AuthHandler) GetQRCode(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

	info := models.QrPayload{}
	if err := request.GetRequestData(r, &info); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, response.ParseBodyError+err.Error())
		response.WriteErrorMessage(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	byteSecret, err := h.uc.GenerateAndUpdateSecret(r.Context(), info.Username)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		response.WriteErrorMessage(w, http.StatusBadRequest, "wrong username")
		return
	}

	secret := base32.StdEncoding.EncodeToString(byteSecret)

	URL, err := url.Parse("otpauth://totp")
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	secretParam := url.Values{}
	secretParam.Add("secret", secret)

	issuer := url.Values{}
	issuer.Add("issuer", qrIssuer)

	URL.RawQuery = secretParam.Encode() + "&" + issuer.Encode()
	URL.Path += fmt.Sprintf("/%s:%s", url.PathEscape(qrIssuer), url.PathEscape(info.Username))

	var png []byte
	png, _ = qrcode.Encode(URL.String(), qrcode.Medium, 256)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(png)

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

func (h *AuthHandler) CheckQRCode(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

	info := models.OtpPayload{}
	if err := request.GetRequestData(r, &info); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, response.ParseBodyError+err.Error())
		response.WriteErrorMessage(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	byteSecret, err := h.uc.GetSecret(r.Context(), info.Username)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		response.WriteErrorMessage(w, http.StatusBadRequest, "wrong username")
		return
	}

	otpConfig := &dgoogauth.OTPConfig{
		Secret:      base32.StdEncoding.EncodeToString(byteSecret),
		WindowSize:  30,
		HotpCounter: 0,
	}

	success, err := otpConfig.Authenticate(info.QrCode)
	if success && err == nil {
		log.LogHandlerInfo(logger, http.StatusAccepted, "qr code is good")
		w.WriteHeader(http.StatusAccepted)
	} else {
		logErrorMessage := "success is false"
		if err != nil {
			logErrorMessage = err.Error()
		}
		log.LogHandlerError(logger, http.StatusUnauthorized, response.ParseBodyError+logErrorMessage)
		w.WriteHeader(http.StatusUnauthorized)
	}
}
