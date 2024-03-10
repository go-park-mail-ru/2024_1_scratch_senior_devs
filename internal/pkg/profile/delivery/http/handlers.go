package http

import (
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

func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	payload := models.ProfileUpdatePayload{}
	err := utils.GetRequestData(r, &payload)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "incorrect data format")
		return
	}
}
