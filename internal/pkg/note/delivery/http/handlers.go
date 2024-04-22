package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc/gen"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/paging"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/responses"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

const TimeLayout = "2006-01-02 15:04:05 -0700 UTC"

type NoteHandler struct {
	client gen.NoteClient
}

func CreateNotesHandler(client gen.NoteClient) *NoteHandler {
	return &NoteHandler{
		client: client,
	}
}

const (
	incorrectIdErr = "incorrect id parameter"
)

// GetAllNotes godoc
// @Summary		Get all notes
// @Description	Get a list of notes of current user
// @Tags 		note
// @ID			get-all-notes
// @Produce		json
// @Param		count	query		int							false	"notes count"
// @Param		offset	query		int							false	"notes offset"
// @Param		title	query		string						false	"notes title substring"
// @Success		200		{object}	[]models.NoteForSwagger		true	"notes"
// @Failure		400		{object}	responses.ErrorResponse		true	"error"
// @Failure		401
// @Router		/api/note/get_all [get]
func (h *NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	count, offset, err := paging.GetParams(r)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("invalid parameters"))
		return
	}

	titleSubstr := r.URL.Query().Get("title")

	_, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	protoData, err := h.client.GetAllNotes(r.Context(), &gen.GetAllRequest{
		Count:  int64(count),
		Offset: int64(offset),
		Title:  titleSubstr,
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, err)
		return
	}
	data := make([]models.Note, len(protoData.Notes))

	for i, note := range protoData.Notes {
		createTime, err := time.Parse(TimeLayout, note.CreateTime)
		if err != nil {
			log.LogHandlerError(logger, http.StatusInternalServerError, err.Error())
		}
		updateTime, err := time.Parse(TimeLayout, note.UpdateTime)
		if err != nil {
			log.LogHandlerError(logger, http.StatusInternalServerError, err.Error())
		}
		data[i] = models.Note{
			Id:         uuid.FromStringOrNil(note.Id),
			Data:       []byte(note.Data),
			CreateTime: createTime,
			UpdateTime: updateTime,
			OwnerId:    uuid.FromStringOrNil(note.OwnerId),
		}
	}
	if err := responses.WriteResponseData(w, data, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

// GetNote godoc
// @Summary		Get one note
// @Description	Get one of notes of current user
// @Tags 		note
// @ID			get-note
// @Param		id		path		string					true	"note id"
// @Success		200		{object}	models.NoteForSwagger	true	"note"
// @Failure		400		{object}	responses.ErrorResponse	true	"incorrect id"
// @Failure		401
// @Failure		404		{object}	responses.ErrorResponse	true	"note not found"
// @Router		/api/note/{id} [get]
func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	noteIdString := mux.Vars(r)["id"]
	_, err := uuid.FromString(noteIdString)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, incorrectIdErr+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("note id must be a type of uuid"))
		return
	}

	_, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	protoNote, err := h.client.GetNote(r.Context(), &gen.GetNoteRequest{Id: noteIdString})
	if err != nil {
		log.LogHandlerError(logger, http.StatusNotFound, err.Error())
		responses.WriteErrorMessage(w, http.StatusNotFound, err)
		return
	}
	createTime, err := time.Parse(TimeLayout, protoNote.Note.CreateTime)
	if err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, err.Error())

	}
	updateTime, err := time.Parse(TimeLayout, protoNote.Note.UpdateTime)
	if err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, err.Error())

	}
	resultNote := models.Note{
		Id:         uuid.FromStringOrNil(protoNote.Note.Id),
		Data:       []byte(protoNote.Note.Data),
		CreateTime: createTime,
		UpdateTime: updateTime,
		OwnerId:    uuid.FromStringOrNil(protoNote.Note.OwnerId),
	}
	if err := responses.WriteResponseData(w, resultNote, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

// AddNote godoc
// @Summary		Add note
// @Description	Create new note to current user
// @Tags 		note
// @ID			add-note
// @Accept		json
// @Produce		json
// @Param    	payload		body    	models.UpsertNoteRequestForSwagger  	true  	"note data"
// @Success		200			{object}	models.NoteForSwagger					true	"note"
// @Failure		400			{object}	responses.ErrorResponse					true	"error"
// @Failure		401
// @Router		/api/note/add [post]
func (h *NoteHandler) AddNote(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	_, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	payload := models.UpsertNoteRequest{}
	if err := responses.GetRequestData(r, &payload); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, responses.ParseBodyError+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("incorrect data format"))
		return
	}

	noteData, err := json.Marshal(payload.Data)
	if err != nil {

		log.LogHandlerError(logger, http.StatusBadRequest, responses.ParseBodyError+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("incorrect data format"))
		return
	}

	protoNote, err := h.client.AddNote(r.Context(), &gen.AddNoteRequest{Data: string(noteData)})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("invalid query"))
		return
	}
	createTime, err := time.Parse(TimeLayout, protoNote.Note.CreateTime)
	if err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, err.Error())

	}
	updateTime, err := time.Parse(TimeLayout, protoNote.Note.UpdateTime)
	if err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, err.Error())

	}
	resultNote := models.Note{
		Id:         uuid.FromStringOrNil(protoNote.Note.Id),
		Data:       []byte(protoNote.Note.Data),
		CreateTime: createTime,
		UpdateTime: updateTime,
		OwnerId:    uuid.FromStringOrNil(protoNote.Note.OwnerId),
	}
	if err := responses.WriteResponseData(w, resultNote, http.StatusCreated); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

// UpdateNote godoc
// @Summary		Update note
// @Description	Create new note to current user
// @Tags 		note
// @ID			update-note
// @Accept		json
// @Produce		json
// @Param		id			path		string								true	"note id"
// @Param    	payload		body    	models.UpsertNoteRequestForSwagger	true  	"note data"
// @Success		200			{object}	models.NoteForSwagger				true	"note"
// @Failure		400			{object}	responses.ErrorResponse				true	"error"
// @Failure		401
// @Router		/api/note/{id}/edit [post]
func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	noteIdString := mux.Vars(r)["id"]
	_, err := uuid.FromString(noteIdString)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, incorrectIdErr+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("note id must be a type of uuid"))
		return
	}

	_, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	payload := models.UpsertNoteRequest{}
	if err := responses.GetRequestData(r, &payload); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, responses.ParseBodyError+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("incorrect data format"))
		return
	}

	noteData, err := json.Marshal(payload.Data)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, responses.ParseBodyError+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("incorrect data format"))
		return
	}

	protoNote, err := h.client.UpdateNote(r.Context(), &gen.UpdateNoteRequest{Data: string(noteData), Id: noteIdString})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("note not found"))
		return
	}
	createTime, err := time.Parse(TimeLayout, protoNote.Note.CreateTime)
	if err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, err.Error())

	}
	updateTime, err := time.Parse(TimeLayout, protoNote.Note.UpdateTime)
	if err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, err.Error())

	}
	resultNote := models.Note{
		Id:         uuid.FromStringOrNil(protoNote.Note.Id),
		Data:       []byte(protoNote.Note.Data),
		CreateTime: createTime,
		UpdateTime: updateTime,
		OwnerId:    uuid.FromStringOrNil(protoNote.Note.OwnerId),
	}
	if err := responses.WriteResponseData(w, resultNote, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

// DeleteNote godoc
// @Summary		Delete note
// @Description	Delete selected note of current user
// @Tags 		note
// @ID			delete-note
// @Param		id		path		string					true	"note id"
// @Success		204
// @Failure		400		{object}	responses.ErrorResponse	true	"incorrect id"
// @Failure		401
// @Failure		404		{object}	responses.ErrorResponse	true	"note not found"
// @Router		/api/note/{id}/delete [delete]
func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	noteIdString := mux.Vars(r)["id"]
	_, err := uuid.FromString(noteIdString)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, incorrectIdErr+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("note id must be a type of uuid"))
		return
	}

	_, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err = h.client.DeleteNote(r.Context(), &gen.DeleteNoteRequest{Id: noteIdString})
	if err != nil {
		log.LogHandlerError(logger, http.StatusNotFound, err.Error())
		responses.WriteErrorMessage(w, http.StatusNotFound, errors.New("note not found"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.LogHandlerInfo(logger, http.StatusNoContent, "success")
}
