package http

import (
	"fmt"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/profile"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
	"net/http"
)

type ProfileHandler struct {
	uc profile.ProfileUsecase
}

func CreateProfileHandler(uc profile.ProfileUsecase) *ProfileHandler {
	return &ProfileHandler{
		uc: uc,
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
	jwtPayload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	payload := models.ProfileUpdatePayload{}
	err := utils.GetRequestData(r, &payload)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	user, err := h.uc.UpdateProfile(r.Context(), jwtPayload.Id, payload)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	err = utils.WriteResponseData(w, user, http.StatusCreated)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error in SignUp handler: %s", err)
		return
	}
}
