package http

import (
	"encoding/json"
	"errors"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/hub"
	"github.com/gorilla/websocket"
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
	hub    hub.HubInterface
}

func CreateNotesHandler(client gen.NoteClient, hub hub.HubInterface) *NoteHandler {
	return &NoteHandler{
		client: client,
		hub:    hub,
	}
}

const (
	incorrectIdErr = "incorrect id parameter"
)

var (
	upgrader = websocket.Upgrader{}
)

func getNote(note *gen.NoteModel) (models.Note, error) {
	createTime, err := time.Parse("2006-01-02 15:04:05 -0700 UTC", note.CreateTime)
	if err != nil {
		return models.Note{}, err
	}

	updateTime, err := time.Parse("2006-01-02 15:04:05 -0700 UTC", note.UpdateTime)
	if err != nil {
		return models.Note{}, err
	}

	children := make([]uuid.UUID, len(note.Children))
	for i, child := range note.Children {
		children[i] = uuid.FromStringOrNil(child)
	}

	return models.Note{
		Id:         uuid.FromStringOrNil(note.Id),
		OwnerId:    uuid.FromStringOrNil(note.OwnerId),
		Data:       []byte(note.Data),
		CreateTime: createTime,
		UpdateTime: updateTime,
		Parent:     uuid.FromStringOrNil(note.Parent),
		Children:   children,
	}, nil
}

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

	payload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	protoData, err := h.client.GetAllNotes(r.Context(), &gen.GetAllRequest{
		Count:  int64(count),
		Offset: int64(offset),
		Title:  titleSubstr,
		UserId: payload.Id.String(),
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, err)
		return
	}
	data := make([]models.Note, len(protoData.Notes))

	for i, note := range protoData.Notes {

		data[i], err = getNote(note)
		if err != nil {
			log.LogHandlerError(logger, http.StatusInternalServerError, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
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

	payload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	protoNote, err := h.client.GetNote(r.Context(), &gen.GetNoteRequest{
		Id:     noteIdString,
		UserId: payload.Id.String(),
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusNotFound, err.Error())
		responses.WriteErrorMessage(w, http.StatusNotFound, err)
		return
	}

	resultNote, err := getNote(protoNote.Note)
	if err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
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

	jwtPayload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
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

	protoNote, err := h.client.AddNote(r.Context(), &gen.AddNoteRequest{
		Data:   string(noteData),
		UserId: jwtPayload.Id.String(),
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("invalid query"))
		return
	}

	resultNote, err := getNote(protoNote.Note)
	if err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
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

	jwtPayload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
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

	protoNote, err := h.client.UpdateNote(r.Context(), &gen.UpdateNoteRequest{
		Data:   string(noteData),
		Id:     noteIdString,
		UserId: jwtPayload.Id.String(),
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("note not found"))
		return
	}

	resultNote, err := getNote(protoNote.Note)
	if err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
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

	payload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err = h.client.DeleteNote(r.Context(), &gen.DeleteNoteRequest{
		Id:     noteIdString,
		UserId: payload.Id.String(),
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusNotFound, err.Error())
		responses.WriteErrorMessage(w, http.StatusNotFound, errors.New("note not found"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.LogHandlerInfo(logger, http.StatusNoContent, "success")
}

// CreateSubNote godoc
// @Summary		Create subnote
// @Description	Create new subnote in current note
// @Tags 		note
// @ID			create-subnote
// @Accept		json
// @Produce		json
// @Param    	payload		body    	models.UpsertNoteRequestForSwagger  	true  	"note data"
// @Param		id			path		string									true	"note id"
// @Success		200			{object}	models.NoteForSwagger					true	"note"
// @Failure		400			{object}	responses.ErrorResponse					true	"error"
// @Failure		401
// @Failure		404		{object}	responses.ErrorResponse	true	"note not found"
// @Router		/api/note/{id}/add_subnote [post]
func (h *NoteHandler) CreateSubNote(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	jwtPayload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	noteIdString := mux.Vars(r)["id"]
	_, err := uuid.FromString(noteIdString)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, incorrectIdErr+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("note id must be a type of uuid"))
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

	subNote, err := h.client.CreateSubNote(r.Context(), &gen.CreateSubNoteRequest{
		UserId:   jwtPayload.Id.String(),
		NoteData: string(noteData),
		ParentId: noteIdString,
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusNotFound, err.Error())
		responses.WriteErrorMessage(w, http.StatusNotFound, errors.New("note not found"))
		return
	}

	resultNote, err := getNote(subNote.Note)
	if err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := responses.WriteResponseData(w, resultNote, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

func (h *NoteHandler) SubscribeOnUpdates(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	_, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	noteIdString := mux.Vars(r)["id"]
	noteID, err := uuid.FromString(noteIdString)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, incorrectIdErr+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("note id must be a type of uuid"))
		return
	}

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("fail to upgrade to websocket"))
		return
	}
	logger.Debug("connection upgraded: ", slog.Any("noteID", noteID))

	h.hub.AddClient(noteID, connection)

	logger.Debug("client disconnected: ", slog.Any("noteID", noteID))
}
