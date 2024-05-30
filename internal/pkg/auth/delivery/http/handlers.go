package http

import (
	"encoding/base32"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/delivery/grpc/gen"
	"github.com/satori/uuid"

	noteGen "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc/gen"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/protection"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/cookie"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/filework"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/responses"
	"github.com/skip2/go-qrcode"
)

const TimeLayout = "2006-01-02 15:04:05 -0700 UTC"

type AuthHandler struct {
	client        gen.AuthClient
	noteClient    noteGen.NoteClient
	cfg           config.AuthHandlerConfig
	cfgValidation config.ValidationConfig
}

func CreateAuthHandler(client gen.AuthClient, noteClient noteGen.NoteClient, cfg config.AuthHandlerConfig, cfgValidation config.ValidationConfig) *AuthHandler {
	return &AuthHandler{
		client:        client,
		noteClient:    noteClient,
		cfg:           cfg,
		cfgValidation: cfgValidation,
	}
}

func makeHelloNoteData(username string) string {
	return fmt.Sprintf(`{"title":"YouNote❤️","content":[{"pluginName":"textBlock","content":"Привет, %s!"},{"pluginName":"div","children":[{"pluginName":"br"}]}]}`, username)
}

func getUser(user *gen.User) (models.User, error) {
	createTime, err := time.Parse(TimeLayout, user.CreateTime)
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		Id:           uuid.FromStringOrNil(user.Id),
		Description:  user.Description,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		CreateTime:   createTime,
		ImagePath:    user.ImagePath,
		SecondFactor: models.Secret(user.SecondFactor),
	}, nil
}

// SignUp godoc
// @Summary		Sign up
// @Description	Add a new user to the database
// @Tags 		auth
// @ID			sign-up
// @Accept		json
// @Produce		json
// @Param		credentials body		models.SignUpPayloadForSwagger	true	"registration data"
// @Success		200			{object}	models.UserForSwagger			true	"user"
// @Failure		400			{object}	responses.ErrorResponse			true	"error"
// @Router		/api/auth/signup [post]
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	userData := models.UserFormData{}
	if err := responses.GetRequestData(r, &userData); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, responses.ParseBodyError+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, auth.ErrIncorrectPayload)
		return
	}

	if err := userData.Validate(h.cfgValidation); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, "validation error: "+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	response, err := h.client.SignUp(r.Context(), &gen.UserFormData{
		Username: userData.Username,
		Password: userData.Password,
		Code:     userData.Code,
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, auth.ErrCreatingUser)
		return
	}

	expTime, err := time.Parse(TimeLayout, response.Expires)
	if err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie.GenJwtTokenCookie(response.Token, expTime, h.cfg.Jwt))
	w.Header().Set("Authorization", "Bearer "+response.Token)

	protection.SetCsrfToken(w, h.cfg.Csrf)

	realUser, err := getUser(response.User)
	if err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := responses.WriteResponseData(w, realUser, http.StatusCreated); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = h.noteClient.AddNote(r.Context(), &noteGen.AddNoteRequest{
		Data:   string(makeHelloNoteData(response.User.Username)),
		UserId: response.User.Id,
	})
	if err != nil {
		logger.Error(err.Error())
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
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	_, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
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
// @Failure		401
// @Router		/api/auth/logout [delete]
func (h *AuthHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	http.SetCookie(w, cookie.DelJwtTokenCookie(h.cfg.Jwt))
	w.Header().Del("Authorization")

	http.SetCookie(w, cookie.DelCsrfTokenCookie(h.cfg.Csrf))
	w.Header().Del("X-Csrf-Token")

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
// @Success		200			{object}	models.UserForSwagger	true	"user"
// @Success		202
// @Failure		400			{object}	responses.ErrorResponse	true	"error"
// @Failure		401			{object}	responses.ErrorResponse	true	"error"
// @Router		/api/auth/login [post]
func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	_, err := h.client.CheckLoginAttempts(r.Context(), &gen.CheckLoginAttemptsRequest{
		IpAddress: r.Header.Get("X-Real-IP"),
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusTooManyRequests, err.Error())
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}

	userData := models.UserFormData{}
	if err := responses.GetRequestData(r, &userData); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, responses.ParseBodyError+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, auth.ErrIncorrectPayload)
		return
	}

	response, err := h.client.SignIn(r.Context(), &gen.UserFormData{
		Username: userData.Username,
		Password: userData.Password,
		Code:     userData.Code,
	})
	if err != nil {
		if err.Error() == "rpc error: code = Unknown desc = first factor passed" {
			log.LogHandlerError(logger, http.StatusAccepted, err.Error())
			w.WriteHeader(http.StatusAccepted)
			return
		}

		log.LogHandlerError(logger, http.StatusUnauthorized, err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expTime, err := time.Parse(TimeLayout, response.Expires)
	if err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie.GenJwtTokenCookie(response.Token, expTime, h.cfg.Jwt))
	w.Header().Set("Authorization", "Bearer "+response.Token)

	protection.SetCsrfToken(w, h.cfg.Csrf)

	realUser, err := getUser(response.User)
	if err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := responses.WriteResponseData(w, realUser, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

// GetProfile godoc
// @Summary		Get profile
// @Description	Get user info if user is authorized
// @Tags 		profile
// @ID			get-profile
// @Produce		json
// @Success		200		{object}	models.UserForSwagger		true	"user"
// @Failure		401
// @Router		/api/profile/get [get]
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	jwtPayload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userId := jwtPayload.Id
	currentUser, err := h.client.CheckUser(r.Context(), &gen.CheckUserRequest{
		UserId: userId.String(),
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusUnauthorized, err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	realUser, err := getUser(currentUser)
	if err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := responses.WriteResponseData(w, realUser, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
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
// @Success		200			{object}	models.UserForSwagger			true	"user"
// @Failure		400			{object}	responses.ErrorResponse			true	"error"
// @Failure		401
// @Router		/api/profile/update [post]
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	jwtPayload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	payload := models.ProfileUpdatePayload{}
	if err := responses.GetRequestData(r, &payload); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, responses.ParseBodyError+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, auth.ErrIncorrectPayload)
		return
	}

	user, err := h.client.UpdateProfile(r.Context(), &gen.UpdateProfileRequest{
		UserId: jwtPayload.Id.String(),
		Payload: &gen.ProfileUpdatePayload{
			Description: payload.Description,
			Password: &gen.Passwords{
				Old: payload.Password.Old,
				New: payload.Password.New,
			},
		},
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	realUser, err := getUser(user)
	if err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := responses.WriteResponseData(w, realUser, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

// UpdateProfileAvatar godoc
// @Summary		Update profile filework
// @Description	Change filework
// @Tags 		profile
// @ID			update-profile-filework
// @Accept		multipart/form-data
// @Produce		json
// @Param		avatar 		formData	file						true	"avatar"
// @Success		200			{object}	models.UserForSwagger		true	"user"
// @Failure		400			{object}	responses.ErrorResponse		true	"error"
// @Failure		401
// @Failure		413
// @Router		/api/profile/update_avatar [post]
func (h *AuthHandler) UpdateProfileAvatar(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	jwtPayload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, int64(h.cfg.AvatarMaxFormDataSize))
	defer r.Body.Close()

	err := r.ParseMultipartForm(int64(h.cfg.AvatarMaxFormDataSize))
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

	files := r.MultipartForm.File["avatar"]
	if len(files) > 1 {
		log.LogHandlerError(logger, http.StatusBadRequest, auth.ErrWrongFilesNumber.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, auth.ErrWrongFilesNumber)
		return
	}

	avatar, _, err := r.FormFile("avatar")
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, auth.ErrWrongFilesNumber)
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

	fileExtension := filework.GetFormat(h.cfg.AvatarFileTypes, content)
	if fileExtension == "" {
		log.LogHandlerError(logger, http.StatusBadRequest, auth.ErrWrongFileFormat.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, auth.ErrWrongFileFormat)
		return
	}

	user, err := h.client.UpdateProfileAvatar(r.Context(), &gen.UpdateProfileAvatarRequest{
		UserId:    jwtPayload.Id.String(),
		Avatar:    content,
		Extension: fileExtension,
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	realUser, err := getUser(user)
	if err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := responses.WriteResponseData(w, realUser, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

// GetQRCode godoc
// @Summary		Get QR code
// @Description	Generate QR code for 2FA
// @Tags 		auth
// @ID			get-qr-code
// @Produce		image/png
// @Success		200
// @Failure		400
// @Failure		401
// @Router		/api/auth/get_qr [get]
func (h *AuthHandler) GetQRCode(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	jwtPayload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	byteSecret, err := h.client.GenerateAndUpdateSecret(r.Context(), &gen.SecretRequest{
		Username: jwtPayload.Username,
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	secret := base32.StdEncoding.EncodeToString(byteSecret.Secret)

	URL, err := url.Parse("otpauth://totp")
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	secretParam := url.Values{}
	secretParam.Add("secret", secret)

	issuer := url.Values{}
	issuer.Add("issuer", h.cfg.QrIssuer)

	URL.RawQuery = secretParam.Encode() + "&" + issuer.Encode()
	URL.Path += fmt.Sprintf("/%s:%s", url.PathEscape(h.cfg.QrIssuer), url.PathEscape(jwtPayload.Username))

	var png []byte
	png, _ = qrcode.Encode(URL.String(), qrcode.Medium, 256)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(png)

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

// DisableSecondFactor godoc
// @Summary		Disable second factor
// @Description	Remove secret for QR-code from database
// @Tags 		auth
// @ID			disable-second-factor
// @Success		204
// @Failure		401
// @Router		/api/auth/disable_2fa [delete]
func (h *AuthHandler) DisableSecondFactor(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	jwtPayload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if _, err := h.client.DeleteSecret(r.Context(), &gen.SecretRequest{
		Username: jwtPayload.Username,
	}); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.LogHandlerInfo(logger, http.StatusNoContent, "success")
}
