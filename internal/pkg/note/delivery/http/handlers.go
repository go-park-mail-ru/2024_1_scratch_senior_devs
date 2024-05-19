package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/exportpdf"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/hub"
	"github.com/gorilla/websocket"

	authGen "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc/gen"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/paging"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/responses"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

const (
	TimeLayout     = "2006-01-02 15:04:05 -0700 UTC"
	RpcErrorPrefix = "rpc error: code = Unknown desc = "
)

type NoteHandler struct {
	client     gen.NoteClient
	authClient authGen.AuthClient
	hub        hub.HubInterface
}

func CreateNotesHandler(client gen.NoteClient, authClient authGen.AuthClient, hub hub.HubInterface) *NoteHandler {
	return &NoteHandler{
		client:     client,
		authClient: authClient,
		hub:        hub,
	}
}

const (
	incorrectIdErr = "incorrect id parameter"
)

var (
	upgrader = websocket.Upgrader{}
)

func getNote(note *gen.NoteModel) (models.Note, error) {
	if note == nil {
		return models.Note{}, errors.New("not found")
	}
	createTime, err := time.Parse(TimeLayout, note.CreateTime)
	if err != nil {
		return models.Note{}, err
	}

	updateTime, err := time.Parse(TimeLayout, note.UpdateTime)
	if err != nil {
		return models.Note{}, err
	}

	children := make([]uuid.UUID, len(note.Children))
	for i, child := range note.Children {
		children[i] = uuid.FromStringOrNil(child)
	}

	collaborators := make([]uuid.UUID, len(note.Collaborators))
	for i, collaborator := range note.Collaborators {
		collaborators[i] = uuid.FromStringOrNil(collaborator)
	}

	tags := make([]string, len(note.Tags))
	copy(tags, note.Tags)

	return models.Note{
		Id:            uuid.FromStringOrNil(note.Id),
		OwnerId:       uuid.FromStringOrNil(note.OwnerId),
		Data:          note.Data,
		CreateTime:    createTime,
		UpdateTime:    updateTime,
		Parent:        uuid.FromStringOrNil(note.Parent),
		Children:      children,
		Tags:          tags,
		Collaborators: collaborators,
		Icon:          note.Icon,
		Header:        note.Header,
		Favorite:      note.Favorite,
		Public:        note.Public,
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
	tagsString := r.URL.Query().Get("tags")

	tagsArray := slices.DeleteFunc(strings.Split(tagsString, "|"), func(e string) bool {
		return e == ""
	})

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
		Tags:   tagsArray,
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	data := make([]models.Note, len(protoData.Notes))
	for i, protoNote := range protoData.Notes {
		data[i], err = getNote(protoNote)
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

func (h *NoteHandler) GetPublicNote(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	noteIdString := mux.Vars(r)["id"]
	_, err := uuid.FromString(noteIdString)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, incorrectIdErr+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("note id must be a type of uuid"))
		return
	}

	protoNote, err := h.client.GetPublicNote(r.Context(), &gen.GetPublicNoteRequest{
		NoteId: noteIdString,
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

	h.hub.WriteToCache(r.Context(), models.CacheMessage{
		Type:        "updated",
		NoteId:      resultNote.Id,
		Username:    jwtPayload.Username,
		Created:     time.Now().UTC(),
		MessageInfo: resultNote.Data,
		SocketID:    payload.SocketID,
	})

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
		if err.Error()[len(RpcErrorPrefix):] == note.ErrTooManySubnotes {
			log.LogHandlerError(logger, http.StatusConflict, err.Error())
			responses.WriteErrorMessage(w, http.StatusConflict, err)
			return
		}

		if err.Error()[len(RpcErrorPrefix):] == note.ErrTooDeep {
			log.LogHandlerError(logger, http.StatusNotFound, err.Error())
			responses.WriteErrorMessage(w, http.StatusNotFound, err)
			return
		}

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

	jwtPayload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
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

	response, err := h.client.CheckPermissions(r.Context(), &gen.CheckPermissionsRequest{
		NoteId: noteIdString,
		UserId: jwtPayload.Id.String(),
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusNotFound, err.Error())
		responses.WriteErrorMessage(w, http.StatusNotFound, errors.New("not found"))
		return
	}
	if !response.Result {
		log.LogHandlerError(logger, http.StatusNotFound, "not owner and not collaborator")
		responses.WriteErrorMessage(w, http.StatusNotFound, errors.New("not found"))
		return
	}

	upgrader.Subprotocols = []string{r.Header.Get("Sec-WebSocket-Protocol")}
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("fail to upgrade to websocket"))
		return
	}

	logger.Info("connection upgraded: ", slog.Any("noteID", noteID))

	h.hub.AddClient(r.Context(), noteID, hub.NewCustomClient(connection))

	logger.Info("client disconnected: ", slog.Any("noteID", noteID))
}

func (h *NoteHandler) AddCollaborator(w http.ResponseWriter, r *http.Request) {
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

	payload := models.AddCollaboratorRequest{}
	if err := responses.GetRequestData(r, &payload); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, responses.ParseBodyError+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("incorrect data format"))
		return
	}

	guest, err := h.authClient.GetUserByUsername(r.Context(), &authGen.GetUserByUsernameRequest{Username: payload.Username})
	if err != nil {
		log.LogHandlerError(logger, http.StatusNotFound, err.Error())
		responses.WriteErrorMessage(w, http.StatusNotFound, errors.New("user not found"))
		return
	}

	if jwtPayload.Id == uuid.FromStringOrNil(guest.Id) {
		log.LogHandlerError(logger, http.StatusBadRequest, "tried to invite self to note")
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("вы реально думали, что можно вот так просто взять и пригласить в заметку самого себя?"))
		return
	}

	_, err = h.client.AddCollaborator(r.Context(), &gen.AddCollaboratorRequest{
		NoteId:  noteIdString,
		UserId:  jwtPayload.Id.String(),
		GuestId: guest.Id,
	})
	if err != nil {
		if err.Error()[len(RpcErrorPrefix):] == note.ErrAlreadyCollaborator {
			log.LogHandlerError(logger, http.StatusConflict, err.Error())
			responses.WriteErrorMessage(w, http.StatusConflict, err)
			return
		}

		if err.Error()[len(RpcErrorPrefix):] == note.ErrTooManyCollaborators {
			log.LogHandlerError(logger, http.StatusExpectationFailed, err.Error())
			responses.WriteErrorMessage(w, http.StatusExpectationFailed, err)
			return
		}

		log.LogHandlerError(logger, http.StatusNotFound, err.Error())
		responses.WriteErrorMessage(w, http.StatusNotFound, errors.New("not found"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.LogHandlerInfo(logger, http.StatusNoContent, "success")
}

func (h *NoteHandler) AddTag(w http.ResponseWriter, r *http.Request) {
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
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("note id must be a type of uuid"))
		return
	}

	var tagName models.TagRequest
	if err := responses.GetRequestData(r, &tagName); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, responses.ParseBodyError+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("incorrect data format"))
		return
	}

	tagName.Sanitize()
	if err = tagName.Validate(); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	currentNote, err := h.client.AddTag(r.Context(), &gen.TagRequest{
		TagName: tagName.TagName,
		NoteId:  noteIdString,
		UserId:  jwtPayload.Id.String(),
	})
	if err != nil {
		if err.Error()[len(RpcErrorPrefix):] == note.ErrTooManyTags {
			log.LogHandlerError(logger, http.StatusConflict, err.Error())
			responses.WriteErrorMessage(w, http.StatusConflict, err)
			return
		}

		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resultNote, err := getNote(currentNote.Note)
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

func (h *NoteHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
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
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("note id must be a type of uuid"))
		return
	}

	var tagName models.TagRequest
	if err := responses.GetRequestData(r, &tagName); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, responses.ParseBodyError+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("incorrect data format"))
		return
	}

	responseNote, err := h.client.DeleteTag(r.Context(), &gen.TagRequest{
		TagName: tagName.TagName,
		NoteId:  noteIdString,
		UserId:  jwtPayload.Id.String(),
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resultNote, err := getNote(responseNote.Note)
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

func (h *NoteHandler) RememberTag(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	jwtPayload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var tagName models.TagRequest
	if err := responses.GetRequestData(r, &tagName); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, responses.ParseBodyError+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("incorrect data format"))
		return
	}

	_, err := h.client.RememberTag(r.Context(), &gen.AllTagRequest{
		TagName: tagName.TagName,
		UserId:  jwtPayload.Id.String(),
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("tag already exists"))
		return
	}

	log.LogHandlerInfo(logger, http.StatusNoContent, "success")
	w.WriteHeader(http.StatusNoContent)
}

func (h *NoteHandler) ForgetTag(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	jwtPayload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var tagName models.TagRequest
	if err := responses.GetRequestData(r, &tagName); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, responses.ParseBodyError+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("incorrect data format"))
		return
	}

	_, err := h.client.ForgetTag(r.Context(), &gen.AllTagRequest{
		TagName: tagName.TagName,
		UserId:  jwtPayload.Id.String(),
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.LogHandlerInfo(logger, http.StatusNoContent, "success")
	w.WriteHeader(http.StatusNoContent)
}
func (h *NoteHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	jwtPayload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var tagName models.UpdateTagRequest
	if err := responses.GetRequestData(r, &tagName); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, responses.ParseBodyError+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("incorrect data format"))
		return
	}
	tagName.Sanitize()
	if err := tagName.Validate(); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := h.client.UpdateTag(r.Context(), &gen.UpdateTagRequest{
		OldTag: tagName.OldTag,
		NewTag: tagName.NewTag,
		UserId: jwtPayload.Id.String(),
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.LogHandlerInfo(logger, http.StatusNoContent, "success")
	w.WriteHeader(http.StatusNoContent)
}

func (h *NoteHandler) SetIcon(w http.ResponseWriter, r *http.Request) {
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
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("note id must be a type of uuid"))
		return
	}

	var iconRequest models.SetIconRequest
	if err := responses.GetRequestData(r, &iconRequest); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, responses.ParseBodyError+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("incorrect data format"))
		return
	}

	currentNote, err := h.client.SetIcon(r.Context(), &gen.SetIconRequest{
		NoteId: noteIdString,
		Icon:   iconRequest.Icon,
		UserId: jwtPayload.Id.String(),
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resultNote, err := getNote(currentNote.Note)
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

func (h *NoteHandler) SetHeader(w http.ResponseWriter, r *http.Request) {
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
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("note id must be a type of uuid"))
		return
	}

	var headerRequest models.SetHeaderRequest
	if err := responses.GetRequestData(r, &headerRequest); err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, responses.ParseBodyError+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("incorrect data format"))
		return
	}

	currentNote, err := h.client.SetHeader(r.Context(), &gen.SetHeaderRequest{
		NoteId: noteIdString,
		Header: headerRequest.Header,
		UserId: jwtPayload.Id.String(),
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resultNote, err := getNote(currentNote.Note)
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

func (h *NoteHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {
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
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("note id must be a type of uuid"))
		return
	}

	protoNote, err := h.client.AddFav(r.Context(), &gen.ChangeFlagRequest{

		NoteId: noteIdString,
		UserId: jwtPayload.Id.String(),
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
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

func (h *NoteHandler) DeleteFavorite(w http.ResponseWriter, r *http.Request) {
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
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("note id must be a type of uuid"))
		return
	}

	protoNote, err := h.client.DelFav(r.Context(), &gen.ChangeFlagRequest{
		NoteId: noteIdString,
		UserId: jwtPayload.Id.String(),
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
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

func (h *NoteHandler) GetTags(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	jwtPayload, ok := r.Context().Value(config.PayloadContextKey).(models.JwtPayload)
	if !ok {
		log.LogHandlerError(logger, http.StatusUnauthorized, responses.JwtPayloadParseError)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	result, err := h.client.GetTags(r.Context(), &gen.GetTagsRequest{
		UserId: jwtPayload.Id.String(),
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tags := make([]string, len(result.Tags))
	copy(tags, result.Tags)
	response := models.GetTagsResponse{Tags: tags}
	if err := responses.WriteResponseData(w, response, http.StatusOK); err != nil {
		log.LogHandlerError(logger, http.StatusInternalServerError, responses.WriteBodyError+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.LogHandlerInfo(logger, http.StatusOK, "success")
}

func (h *NoteHandler) SetPublic(w http.ResponseWriter, r *http.Request) {
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
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("note id must be a type of uuid"))
		return
	}

	currentNote, err := h.client.SetPublic(r.Context(), &gen.AccessModeRequest{
		NoteId: noteIdString,
		UserId: jwtPayload.Id.String(),
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resultNote, err := getNote(currentNote.Note)
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

func (h *NoteHandler) SetPrivate(w http.ResponseWriter, r *http.Request) {
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
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("note id must be a type of uuid"))
		return
	}

	currentNote, err := h.client.SetPrivate(r.Context(), &gen.AccessModeRequest{
		NoteId: noteIdString,
		UserId: jwtPayload.Id.String(),
	})
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resultNote, err := getNote(currentNote.Note)
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

func (h *NoteHandler) ExportToPDF(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.LogHandlerError(logger, http.StatusBadRequest, "can`t read request body: "+err.Error())
		responses.WriteErrorMessage(w, http.StatusBadRequest, errors.New("can`t read request body"))
		return
	}

	resultPDF, noteTitle, err := exportpdf.GeneratePDF(string(payload))
	if err != nil {
		if err.Error() == "invalid input HTML" {
			log.LogHandlerError(logger, http.StatusBadRequest, err.Error())
			responses.WriteErrorMessage(w, http.StatusBadRequest, err)
			return
		}

		if err.Error() == "internal error while parsing processed HTML" {
			log.LogHandlerError(logger, http.StatusInternalServerError, err.Error())
			responses.WriteErrorMessage(w, http.StatusInternalServerError, err)
			return
		}

		log.LogHandlerInfo(logger, http.StatusUnavailableForLegalReasons, err.Error())
		responses.WriteErrorMessage(w, http.StatusUnavailableForLegalReasons, errors.New("прости, но эта заметка слишком сложная для конвертации в PDF. Мы усиленно работаем, чтобы решить эту проблему"))
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", noteTitle))
	_, _ = w.Write(resultPDF)
	w.WriteHeader(http.StatusOK)
	log.LogHandlerInfo(logger, http.StatusOK, "success")
}
