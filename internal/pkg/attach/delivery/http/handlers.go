package http

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/attach"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/delivery"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/sources"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"
	"io"
	"log/slog"
	"net/http"
)

const (
	incorrectIdErr = "incorrect id parameter"
)

type AttachHandler struct {
	uc     attach.AttachUsecase
	logger *slog.Logger
}

func CreateAttachHandler(uc attach.AttachUsecase, logger *slog.Logger) *AttachHandler {
	return &AttachHandler{
		uc:     uc,
		logger: logger,
	}
}

// AddAttach godoc
// @Summary		Add attachment
// @Description	Attach new file to note
// @Tags 		note
// @ID			add-attach
// @Accept		multipart/form-data
// @Produce		json
// @Param		id			path		string						true	"note id"
// @Param		attach 		formData	file						true	"attach file"
// @Success		200			{object}	models.Attach				true	"attach model"
// @Failure		400			{object}	delivery.ErrorResponse		true	"error"
// @Failure		401
// @Failure		413
// @Router		/api/note/{id}/add_attach [post]
func (h *AttachHandler) AddAttach(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

	_, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, delivery.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	noteIdString := mux.Vars(r)["id"]
	noteId, err := uuid.FromString(noteIdString)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, incorrectIdErr+err.Error())
		delivery.WriteErrorMessage(w, http.StatusBadRequest, "note id must be a type of uuid")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, config.AttachMaxFormDataSize)
	defer r.Body.Close()

	err = r.ParseMultipartForm(config.AttachMaxFormDataSize)
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

	files := r.MultipartForm.File["attach"]
	if len(files) > 1 {
		log.LogHandlerError(logger, http.StatusBadRequest, auth.ErrWrongFilesNumber.Error())
		delivery.WriteErrorMessage(w, http.StatusBadRequest, auth.ErrWrongFilesNumber.Error())
		return
	}

	attachFile, _, err := r.FormFile("attach")
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		delivery.WriteErrorMessage(w, http.StatusBadRequest, auth.ErrWrongFilesNumber.Error())
		return
	}
	content, err := io.ReadAll(attachFile)
	if err != nil && !errors.Is(err, io.EOF) {
		if errors.As(err, new(*http.MaxBytesError)) {
			log.LogHandlerError(logger, http.StatusRequestEntityTooLarge, err.Error())
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}
	}

	fileExtension := sources.CheckFormat(config.AttachFileTypes, content)
	if fileExtension == "" {
		log.LogHandlerError(logger, http.StatusBadRequest, auth.ErrWrongFileFormat.Error())
		delivery.WriteErrorMessage(w, http.StatusBadRequest, auth.ErrWrongFileFormat.Error())
		return
	}

	attachModel, err := h.uc.AddAttach(r.Context(), noteId, attachFile, fileExtension)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := delivery.WriteResponseData(w, attachModel, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, delivery.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}
