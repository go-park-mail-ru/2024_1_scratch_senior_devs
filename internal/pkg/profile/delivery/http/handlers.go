package http

import (
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/profile"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
	"log/slog"
	"net/http"
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
		utils.LogHandlerError(logger, http.StatusUnauthorized, err.Error())
		utils.WriteErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := utils.WriteResponseData(w, user, http.StatusCreated); err != nil {
		utils.LogHandlerError(logger, http.StatusInternalServerError, utils.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.LogHandlerInfo(logger, http.StatusCreated, "success")
}
