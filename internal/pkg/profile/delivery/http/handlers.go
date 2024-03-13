package http

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/profile"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
	"io"
	"log/slog"
	"net/http"
)

const (
	maxFormDataSize = 5 * 1024 * 1024
)

type ProfileHandler struct {
	uc     profile.ProfileUsecase
	logger *slog.Logger
}

func CreateProfileHandler(uc profile.ProfileUsecase, logger *slog.Logger) *ProfileHandler {
	return &ProfileHandler{
		uc:     uc,
		logger: logger,
	}
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
// @Failure		400			{object}	utils.ErrorResponse				true	"error"
// @Router		/api/profile/update [post]
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", utils.GetRequestId(r.Context())), slog.String("func", utils.GFN()))

	jwtPayload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		utils.LogHandlerError(logger, http.StatusUnauthorized, utils.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	payload := models.ProfileUpdatePayload{}
	if err := utils.GetRequestData(r, &payload); err != nil {
		utils.LogHandlerError(logger, http.StatusUnauthorized, utils.ParseBodyError+err.Error())
		utils.WriteErrorMessage(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	user, err := h.uc.UpdateProfile(r.Context(), jwtPayload.Id, payload)
	if err != nil {
		utils.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := utils.WriteResponseData(w, user, http.StatusOK); err != nil {
		utils.LogHandlerError(logger, http.StatusInternalServerError, utils.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.LogHandlerInfo(logger, http.StatusOK, "success")
}

// UpdateProfileAvatar godoc
// @Summary		Update profile avatar
// @Description	Change avatar
// @Tags 		profile
// @ID			update-profile-avatar
// @Accept		multipart/form-data
// @Produce		json
// @Param		avatar 		formData	os.File							true	"avatar"
// @Success		200			{object}	models.User						true	"user"
// @Failure		400			{object}	utils.ErrorResponse				true	"error"
// @Failure		413
// @Router		/api/profile/update_avatar [post]
func (h *ProfileHandler) UpdateProfileAvatar(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", utils.GetRequestId(r.Context())), slog.String("func", utils.GFN()))

	jwtPayload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		utils.LogHandlerError(logger, http.StatusUnauthorized, utils.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxFormDataSize)
	defer r.Body.Close()

	err := r.ParseMultipartForm(maxFormDataSize)
	if err != nil {
		utils.LogHandlerError(logger, http.StatusRequestEntityTooLarge, err.Error())
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	avatar, _, err := r.FormFile("avatar")
	if err != nil {
		utils.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		utils.WriteErrorMessage(w, http.StatusBadRequest, "no file found")
		return
	}
	content, err := io.ReadAll(avatar)
	if err != nil && !errors.Is(err, io.EOF) {
		if errors.As(err, new(*http.MaxBytesError)) {
			utils.LogHandlerError(logger, http.StatusRequestEntityTooLarge, err.Error())
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}
	}

	if !utils.CheckFileFormat(content) {
		utils.LogHandlerError(logger, http.StatusBadRequest, "incorrect file format")
		utils.WriteErrorMessage(w, http.StatusBadRequest, "incorrect file format")
		return
	}

	user, err := h.uc.UpdateProfileAvatar(r.Context(), jwtPayload.Id, avatar)
	if err != nil {
		utils.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := utils.WriteResponseData(w, user, http.StatusOK); err != nil {
		utils.LogHandlerError(logger, http.StatusInternalServerError, utils.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.LogHandlerInfo(logger, http.StatusOK, "success")
}
