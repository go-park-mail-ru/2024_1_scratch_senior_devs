package http

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/attach"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/filework"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/responses"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"
)

const (
	incorrectIdErr = "incorrect id parameter"
)

type AttachHandler struct {
	uc  attach.AttachUsecase
	cfg config.AttachConfig
}

func CreateAttachHandler(uc attach.AttachUsecase, cfg config.AttachConfig) *AttachHandler {
	return &AttachHandler{
		uc:  uc,
		cfg: cfg,
	}
}

// AddAttach godoc
// @Summary		Add attachment
// @Description	Attach new file to note
// @Tags 		attach
// @ID			add-attach
// @Accept		multipart/form-data
// @Produce		json
// @Param		id			path		string						true	"note id"
// @Param		attach 		formData	file						true	"attach file"
// @Success		200			{object}	models.Attach				true	"attach model"
// @Failure		400			{object}	responses.ErrorResponse		true	"error"
// @Failure		401
// @Failure		413
// @Router		/api/note/{id}/add_attach [post]
func (h *AttachHandler) AddAttach(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	jwtPayload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	noteIdString := mux.Vars(r)["id"]
	noteId, err := uuid.FromString(noteIdString)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, incorrectIdErr+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("note id must be a type of uuid"))
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, h.cfg.AttachMaxFormDataSize)
	defer r.Body.Close()

	err = r.ParseMultipartForm(h.cfg.AttachMaxFormDataSize)
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
		responses.WriteErrorMessage(w, http.StatusBadRequest, auth.ErrWrongFilesNumber)
		return
	}

	attachFile, _, err := r.FormFile("attach")
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, auth.ErrWrongFilesNumber)
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

	fileExtension := filework.GetFormat(h.cfg.AttachFileTypes, content)
	if fileExtension == "" {
		log.LogHandlerError(logger, http.StatusBadRequest, auth.ErrWrongFileFormat.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, auth.ErrWrongFileFormat)
		return
	}

	attachModel, err := h.uc.AddAttach(r.Context(), noteId, jwtPayload.Id, attachFile, fileExtension)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := responses.WriteResponseData(w, attachModel, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

// DeleteAttach godoc
// @Summary		Delete attach
// @Description	Remove attach from note
// @Tags 		attach
// @ID			delete-attach
// @Param		id			path		string		true	"attach id"
// @Success		204
// @Failure		400			{object}	responses.ErrorResponse			true	"incorrect id"
// @Failure		401
// @Failure		404			{object}	responses.ErrorResponse			true	"not found"
// @Router		/api/attach/{id}/delete [delete]
func (h *AttachHandler) DeleteAttach(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	jwtPayload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	attachIdString := mux.Vars(r)["id"]
	attachId, err := uuid.FromString(attachIdString)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, incorrectIdErr+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("attach id must be a type of uuid"))
		return
	}

	err = h.uc.DeleteAttach(r.Context(), attachId, jwtPayload.Id)
	if err != nil {
		log.LogHandlerError(logger, http.StatusNotFound, err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.LogHandlerInfo(logger, http.StatusNoContent, "success")
}

// GetAttach godoc
// @Summary		Get attach
// @Description	Get attach if it belongs to current user
// @Tags 		attach
// @ID			get-attach
// @Param		id		path		string					true	"attach id"
// @Produce		image/webp
// @Success		200
// @Failure		401
// @Router		/api/attaches/{id} [get]
func (h *AttachHandler) GetAttach(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	jwtPayload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	attachIdString := mux.Vars(r)["id"]
	attachId, err := uuid.FromString(attachIdString)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, incorrectIdErr+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("attach id must be a type of uuid"))
		return
	}
	attach, err := h.uc.GetAttach(r.Context(), attachId, jwtPayload.Id)
	if err != nil {
		log.LogHandlerError(logger, http.StatusNotFound, err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	targetPath := path.Join(os.Getenv("ATTACHES_BASE_PATH"), attach.Path)
	log.LogHandlerInfo(logger, http.StatusOK, "success")
	fileInfo, err := os.Stat(targetPath)
	if err != nil {
		responses.WriteErrorMessage(w, http.StatusNotFound, errors.New("file not found"))
		return
	}
	w.Header().Add("etag", fileInfo.ModTime().UTC().String())
	header := fmt.Sprintf("attachment; filename=\"%s\"", attach.Path)
	w.Header().Add("Content-Disposition", header)
	http.ServeFile(w, r, targetPath)

}
