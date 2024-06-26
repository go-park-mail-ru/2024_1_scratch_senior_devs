package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	authGen "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/delivery/grpc/gen"
	mock_auth "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/delivery/grpc/gen/mocks"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	mock_hub "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/hub/mocks"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/protection"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc/gen"
	mock_grpc "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc/gen/mocks"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/exportpdf"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNoteHandler_GetAllNotes(t *testing.T) {
	const successTestName = "Test Success"

	tests := []struct {
		name           string
		expectedStatus int
		id             uuid.UUID
		username       string
		expectedData   []*gen.NoteResponseModel
	}{
		{

			name:           successTestName,
			id:             uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
			username:       "testuser",
			expectedStatus: http.StatusOK,
			expectedData: []*gen.NoteResponseModel{
				{
					Id:            "c80e3ea8-0813-4731-b6ee-b41604c56f95",
					OwnerId:       "a233ea8-0813-4731-b12e-b41604c56f95",
					UpdateTime:    time.Time{}.String(),
					CreateTime:    time.Time{}.String(),
					Data:          "",
					Parent:        "a233ea8-0813-4731-b12e-b41604c56f95",
					Children:      []string{},
					Tags:          []string{},
					Collaborators: []string{},
				}, {

					Id:            "c80e3ea8-0813-4731-b12e-b41604c56f95",
					OwnerId:       "a233ea8-0813-4731-b12e-b41604c56f95",
					UpdateTime:    time.Time{}.String(),
					CreateTime:    time.Time{}.String(),
					Data:          "",
					Parent:        "a233ea8-0813-4731-b12e-b41604c56f95",
					Children:      []string{},
					Tags:          []string{},
					Collaborators: []string{},
				},
			},
		},
		{

			name:           "Test Unauthorized",
			id:             uuid.FromStringOrNil(""),
			username:       "",
			expectedStatus: http.StatusUnauthorized,
			expectedData:   []*gen.NoteResponseModel{},
		},
		{

			name:           "Test Error",
			id:             uuid.FromStringOrNil(""),
			username:       "",
			expectedStatus: http.StatusBadRequest,
			expectedData:   []*gen.NoteResponseModel{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mock_grpc.NewMockNoteClient(ctrl)
			mockAuthClient := mock_auth.NewMockAuthClient(ctrl)
			mockHub := mock_hub.NewMockHubInterface(ctrl)
			defer ctrl.Finish()
			req := httptest.NewRequest("GET", "http://example.com/api/handler", nil)
			w := httptest.NewRecorder()
			ctx := context.WithValue(req.Context(), config.PayloadContextKey, models.JwtPayload{Id: tt.id, Username: tt.username})

			if tt.name == "Test Unauthorized" {
				ctx = context.WithValue(req.Context(), config.PayloadContextKey, &gen.GetAllResponse{})

			}
			if tt.name == successTestName {
				mockClient.EXPECT().GetAllNotes(ctx, &gen.GetAllRequest{Count: int64(10), Offset: int64(0), Title: "", UserId: tt.id.String(), Tags: []string{}}, gomock.Any()).Return(&gen.GetAllResponse{Notes: tt.expectedData}, nil)
			}
			if tt.name == "Test Error" {
				mockClient.EXPECT().GetAllNotes(ctx, &gen.GetAllRequest{Count: int64(10), Offset: int64(0), Title: "", UserId: tt.id.String(), Tags: []string{}}, gomock.Any()).Return(&gen.GetAllResponse{Notes: tt.expectedData}, errors.New("error"))

			}
			req = req.WithContext(ctx)

			h := CreateNotesHandler(mockClient, mockAuthClient, mockHub)
			h.GetAllNotes(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.name == successTestName {
				expectedResult := make([]models.Note, len(tt.expectedData))
				for i, expectedNote := range tt.expectedData {
					crTime, _ := time.Parse("2006-01-02 15:04:05 -0700 UTC", expectedNote.CreateTime)
					upTime, _ := time.Parse("2006-01-02 15:04:05 -0700 UTC", expectedNote.UpdateTime)
					expectedResult[i] = models.Note{
						Id:            uuid.FromStringOrNil(expectedNote.Id),
						Data:          expectedNote.Data,
						CreateTime:    crTime,
						UpdateTime:    upTime,
						Parent:        uuid.FromStringOrNil(expectedNote.OwnerId),
						OwnerId:       uuid.FromStringOrNil(expectedNote.OwnerId),
						Children:      []uuid.UUID{},
						Tags:          []string{},
						Collaborators: []uuid.UUID{},
					}

				}
			}
		})
	}
}

func TestNoteHandler_GetNote(t *testing.T) {
	const successTestName = "Test Success"

	tests := []struct {
		name string

		expectedStatus int
		userId         uuid.UUID
		noteId         uuid.UUID
		username       string
		expectedData   models.NoteResponse
	}{
		{
			name:           successTestName,
			expectedStatus: 200,
			noteId:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
			userId:         uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
			username:       "test_user",
			expectedData: models.NoteResponse{
				Note: models.Note{
					Id:            uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
					OwnerId:       uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
					UpdateTime:    time.Time{},
					CreateTime:    time.Time{},
					Data:          "",
					Children:      []uuid.UUID{},
					Tags:          []string{},
					Collaborators: []uuid.UUID{},
				},
			},
		},
		{
			name:           "Test Unauthorized",
			expectedStatus: 401,
			noteId:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
			userId:         uuid.FromStringOrNil(""),
			username:       "test_user",
			expectedData:   models.NoteResponse{},
		},
		{
			name:           "Test Error",
			expectedStatus: http.StatusNotFound,
			noteId:         uuid.FromStringOrNil(""),
			userId:         uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
			username:       "test_user",
			expectedData:   models.NoteResponse{},
		},
		{
			name:           "Test Bad Request",
			expectedStatus: http.StatusBadRequest,
			noteId:         uuid.FromStringOrNil(""),
			userId:         uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
			username:       "test_user",
			expectedData:   models.NoteResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mock_grpc.NewMockNoteClient(ctrl)
			mockAuthClient := mock_auth.NewMockAuthClient(ctrl)
			mockHub := mock_hub.NewMockHubInterface(ctrl)
			defer ctrl.Finish()
			req := httptest.NewRequest("GET", "http://example.com/api/note/c80e3ea8-0813-4731-b6ee-b41604c56f95", nil)
			w := httptest.NewRecorder()
			ctx := context.WithValue(req.Context(), config.PayloadContextKey, models.JwtPayload{Id: tt.userId, Username: tt.username})
			if tt.name == "Test Unauthorized" {
				ctx = context.WithValue(req.Context(), config.PayloadContextKey, "")

			}
			if tt.name == successTestName {
				mockClient.EXPECT().GetNote(gomock.Any(), &gen.GetNoteRequest{
					Id:     tt.noteId.String(),
					UserId: tt.userId.String(),
				}).Return(&gen.GetNoteResponseResponse{Note: &gen.NoteResponseModel{
					Id:            tt.expectedData.Id.String(),
					CreateTime:    tt.expectedData.CreateTime.String(),
					UpdateTime:    tt.expectedData.UpdateTime.String(),
					Data:          tt.expectedData.Data,
					OwnerId:       tt.userId.String(),
					Parent:        tt.expectedData.Parent.String(),
					Children:      []string{},
					Tags:          []string{},
					Collaborators: []string{},
				}}, nil)
			}
			if tt.name == "Test Error" {
				mockClient.EXPECT().GetNote(gomock.Any(), &gen.GetNoteRequest{
					Id:     tt.noteId.String(),
					UserId: tt.userId.String(),
				}).Return(nil, errors.New("error"))

			}
			req = req.WithContext(ctx)
			if tt.name != "Test Bad Request" {
				req = mux.SetURLVars(req, map[string]string{"id": tt.noteId.String()})
			}
			h := CreateNotesHandler(mockClient, mockAuthClient, mockHub)
			h.GetNote(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.name == successTestName {
				d, _ := json.Marshal(tt.expectedData)
				assert.Equal(t, w.Body.Bytes(), d)
			}
		})
	}
}

func TestNoteHandler_AddNote(t *testing.T) {
	id := uuid.NewV4()
	currTime := time.Now().UTC()

	var tests = []struct {
		name           string
		requestBody    string
		noteData       []byte
		usecaseErr     bool
		expectedStatus int
	}{
		{
			name:           "Test_NoteHandler_AddNote_Success",
			requestBody:    `{"data":{"title": "my note"}}`,
			noteData:       []byte(`{"title":"my note"}`),
			usecaseErr:     false,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Test_NoteHandler_AddNote_Fail_1",
			requestBody:    `{"data":{"title": "my note"}`,
			noteData:       []byte(`{"title":"my note"}`),
			usecaseErr:     true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Test_NoteHandler_AddNote_Fail_2",
			requestBody:    `{"data":{"title": "my note"}}`,
			noteData:       []byte(`{"title":"my note"}`),
			usecaseErr:     true,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mock_grpc.NewMockNoteClient(ctrl)
			mockAuthClient := mock_auth.NewMockAuthClient(ctrl)
			mockHub := mock_hub.NewMockHubInterface(ctrl)

			defer ctrl.Finish()

			if tt.name != "Test_NoteHandler_AddNote_Fail_1" {
				call := mockClient.EXPECT().AddNote(gomock.Any(), &gen.AddNoteRequest{Data: string(tt.noteData), UserId: id.String()})
				if tt.usecaseErr {
					call.Return(&gen.AddNoteResponse{}, errors.New("usecase error"))
				} else {
					call.Return(&gen.AddNoteResponse{
						Note: &gen.NoteModel{
							Id:         id.String(),
							Data:       string(tt.noteData),
							CreateTime: currTime.String(),
							UpdateTime: currTime.String(),
							OwnerId:    id.String(),
						},
					}, nil)
				}
			}

			r := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewBufferString(tt.requestBody))
			ctx := context.WithValue(r.Context(), config.PayloadContextKey, models.JwtPayload{
				Id:       id,
				Username: "username",
			})
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			handler := CreateNotesHandler(mockClient, mockAuthClient, mockHub)
			handler.AddNote(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestNoteHandler_DeleteNote(t *testing.T) {
	userId := uuid.NewV4()
	noteId := uuid.NewV4()
	tests := []struct {
		name           string
		expectedStatus int
		mockUsecase    func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface)
	}{
		{
			name:           "Test_DeleteNote_Success",
			expectedStatus: http.StatusNoContent,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {

				mockClient.EXPECT().DeleteNote(gomock.Any(), &gen.DeleteNoteRequest{
					Id:     noteId.String(),
					UserId: userId.String(),
				}).Return(&gen.DeleteNoteResponse{}, nil)
			},
		},
		{
			name:           "Test_DeleteNote_Unauthorized",
			expectedStatus: http.StatusUnauthorized,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {

			},
		},
		{
			name:           "Test_DeleteNote_Fail_BadRequest",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {

			},
		},
		{
			name:           "Test_DeleteNote_Fail_NotFound",
			expectedStatus: http.StatusNotFound,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().DeleteNote(gomock.Any(), &gen.DeleteNoteRequest{
					Id:     noteId.String(),
					UserId: userId.String(),
				}).Return(&gen.DeleteNoteResponse{}, errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mock_grpc.NewMockNoteClient(ctrl)
			mockAuthClient := mock_auth.NewMockAuthClient(ctrl)
			mockHub := mock_hub.NewMockHubInterface(ctrl)

			defer ctrl.Finish()

			r := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewReader([]byte{}))
			ctx := context.Background()
			if tt.expectedStatus != http.StatusUnauthorized {
				ctx = context.WithValue(r.Context(), config.PayloadContextKey, models.JwtPayload{
					Id:       userId,
					Username: "username",
				})
			}
			w := httptest.NewRecorder()
			r = r.WithContext(ctx)
			if tt.expectedStatus == http.StatusBadRequest {
				r = mux.SetURLVars(r, map[string]string{"id": ""})

			} else {
				r = mux.SetURLVars(r, map[string]string{"id": noteId.String()})

			}

			handler := CreateNotesHandler(mockClient, mockAuthClient, mockHub)
			tt.mockUsecase(mockClient, mockAuthClient, mockHub)
			handler.DeleteNote(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestNoteHandler_GetTags(t *testing.T) {
	userId := uuid.NewV4()

	tests := []struct {
		name             string
		expectedStatus   int
		mockUsecase      func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface)
		expectedResponse models.GetTagsResponse
	}{
		{
			name:           "Test_GetTags_Success",
			expectedStatus: http.StatusOK,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().GetTags(gomock.Any(), gomock.Any()).Return(&gen.GetTagsResponse{
					Tags: []string{"first", "second"},
				}, nil)
			},
			expectedResponse: models.GetTagsResponse{
				Tags: []string{"first", "second"},
			},
		},
		{
			name:           "Test_GetTags_Unauthorized",
			expectedStatus: http.StatusUnauthorized,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {

			},
			expectedResponse: models.GetTagsResponse{},
		},
		{
			name:           "Test_GetTags_BadRequest",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().GetTags(gomock.Any(), gomock.Any()).Return(&gen.GetTagsResponse{}, errors.New("error"))
			},
			expectedResponse: models.GetTagsResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mock_grpc.NewMockNoteClient(ctrl)
			mockAuthClient := mock_auth.NewMockAuthClient(ctrl)
			mockHub := mock_hub.NewMockHubInterface(ctrl)

			defer ctrl.Finish()

			r := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewReader([]byte{}))
			ctx := context.Background()
			if tt.expectedStatus != http.StatusUnauthorized {
				ctx = context.WithValue(r.Context(), config.PayloadContextKey, models.JwtPayload{
					Id:       userId,
					Username: "username",
				})
			}
			w := httptest.NewRecorder()
			r = r.WithContext(ctx)

			handler := CreateNotesHandler(mockClient, mockAuthClient, mockHub)
			tt.mockUsecase(mockClient, mockAuthClient, mockHub)
			handler.GetTags(w, r)

			data, _ := json.Marshal(tt.expectedResponse)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, data, w.Body.Bytes())
			}
		})
	}
}

func TestNoteHandler_AddTag(t *testing.T) {
	userId := uuid.NewV4()
	noteId := uuid.NewV4()

	tests := []struct {
		requestBody      []byte
		name             string
		expectedStatus   int
		mockUsecase      func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface)
		expectedResponse models.Note
	}{
		{
			requestBody:    []byte("{\"tag_name\":\"tag\"}"),
			name:           "Test_AddTag_Success",
			expectedStatus: http.StatusOK,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().AddTag(gomock.Any(), &gen.TagRequest{
					TagName: "tag",
					NoteId:  noteId.String(),
					UserId:  userId.String(),
				}).Return(&gen.GetNoteResponse{Note: &gen.NoteModel{
					Id:            noteId.String(),
					OwnerId:       userId.String(),
					Tags:          []string{"tag"},
					CreateTime:    time.Time{}.String(),
					UpdateTime:    time.Time{}.String(),
					Parent:        uuid.UUID{}.String(),
					Collaborators: []string{},
					Children:      []string{},
				}}, nil)
			},
			expectedResponse: models.Note{
				Id:            noteId,
				OwnerId:       userId,
				UpdateTime:    time.Time{},
				CreateTime:    time.Time{},
				Data:          "",
				Children:      []uuid.UUID{},
				Tags:          []string{"tag"},
				Collaborators: []uuid.UUID{},
			},
		},
		{
			requestBody:    []byte(""),
			name:           "Test_AddTag_Unauthorized",
			expectedStatus: http.StatusUnauthorized,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
			expectedResponse: models.Note{},
		},
		{
			requestBody:    []byte(""),
			name:           "Test_AddTag_BadRequest",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
			expectedResponse: models.Note{},
		},
		{
			requestBody:    []byte("{\"tag_name\":\"tag\"}"),
			name:           "Test_AddTag_BadRequestClientErr",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().AddTag(gomock.Any(), &gen.TagRequest{
					TagName: "tag",
					NoteId:  noteId.String(),
					UserId:  userId.String(),
				}).Return(&gen.GetNoteResponse{}, errors.New("rpc error: code = Unknown desc = error"))
			},
			expectedResponse: models.Note{},
		},
		{
			requestBody:    []byte("{\"tag_name\":\"taggggggggggggggggggggggggggggggggggggg\"}"),
			name:           "Test_AddTag_BadRequest_TagInvalid",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {

			},
			expectedResponse: models.Note{},
		},
		{
			requestBody:    []byte("{\"tag_name\":\"tag\"}"),
			name:           "Test_AddTag_BadRequest_gteNoteErr",
			expectedStatus: http.StatusInternalServerError,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().AddTag(gomock.Any(), &gen.TagRequest{
					TagName: "tag",
					NoteId:  noteId.String(),
					UserId:  userId.String(),
				}).Return(&gen.GetNoteResponse{Note: &gen.NoteModel{
					Id:      noteId.String(),
					OwnerId: userId.String(),
					Tags:    []string{"tag"},

					Parent:        uuid.UUID{}.String(),
					Collaborators: []string{},
					Children:      []string{},
				}}, nil)
			},
			expectedResponse: models.Note{},
		},
		{
			requestBody:    []byte("{\"tag_name\":\"tag\"}"),
			name:           "Test_AddTag_TooManyErr",
			expectedStatus: http.StatusConflict,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().AddTag(gomock.Any(), &gen.TagRequest{
					TagName: "tag",
					NoteId:  noteId.String(),
					UserId:  userId.String(),
				}).Return(&gen.GetNoteResponse{}, errors.New(RpcErrorPrefix+note.ErrTooManyTags))
			},
			expectedResponse: models.Note{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mock_grpc.NewMockNoteClient(ctrl)
			mockAuthClient := mock_auth.NewMockAuthClient(ctrl)
			mockHub := mock_hub.NewMockHubInterface(ctrl)

			defer ctrl.Finish()

			r := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewReader(tt.requestBody))
			ctx := context.Background()
			if tt.expectedStatus != http.StatusUnauthorized {
				ctx = context.WithValue(r.Context(), config.PayloadContextKey, models.JwtPayload{
					Id:       userId,
					Username: "username",
				})
			}
			w := httptest.NewRecorder()
			r = r.WithContext(ctx)
			if tt.name == "Test_AddTag_BadRequest" {
				r = mux.SetURLVars(r, map[string]string{"id": ""})

			} else {
				r = mux.SetURLVars(r, map[string]string{"id": noteId.String()})

			}

			handler := CreateNotesHandler(mockClient, mockAuthClient, mockHub)
			tt.mockUsecase(mockClient, mockAuthClient, mockHub)
			handler.AddTag(w, r)

			data, _ := json.Marshal(tt.expectedResponse)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, data, w.Body.Bytes())
			}
		})
	}
}

func TestNoteHandler_DeleteTag(t *testing.T) {
	userId := uuid.NewV4()
	noteId := uuid.NewV4()

	tests := []struct {
		requestBody      []byte
		name             string
		expectedStatus   int
		mockUsecase      func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface)
		expectedResponse models.Note
	}{
		{
			requestBody:    []byte("{\"tag_name\":\"tag\"}"),
			name:           "Test_DeleteTag_Success",
			expectedStatus: http.StatusOK,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().DeleteTag(gomock.Any(), &gen.TagRequest{
					TagName: "tag",
					NoteId:  noteId.String(),
					UserId:  userId.String(),
				}).Return(&gen.GetNoteResponse{Note: &gen.NoteModel{
					Id:            noteId.String(),
					OwnerId:       userId.String(),
					Tags:          []string{},
					CreateTime:    time.Time{}.String(),
					UpdateTime:    time.Time{}.String(),
					Parent:        uuid.UUID{}.String(),
					Collaborators: []string{},
					Children:      []string{},
				}}, nil)
			},
			expectedResponse: models.Note{
				Id:            noteId,
				OwnerId:       userId,
				UpdateTime:    time.Time{},
				CreateTime:    time.Time{},
				Data:          "",
				Children:      []uuid.UUID{},
				Tags:          []string{},
				Collaborators: []uuid.UUID{},
			},
		},
		{
			requestBody:    []byte(""),
			name:           "Test_DeleteTag_Unauthorized",
			expectedStatus: http.StatusUnauthorized,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
			expectedResponse: models.Note{},
		},
		{
			requestBody:    []byte(""),
			name:           "Test_DeleteTag_BadRequest",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
			expectedResponse: models.Note{},
		},
		{
			requestBody:    []byte("{\"tag_name\":\"tag\"}"),
			name:           "Test_DeleteTag_BadRequestClientErr",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().DeleteTag(gomock.Any(), &gen.TagRequest{
					TagName: "tag",
					NoteId:  noteId.String(),
					UserId:  userId.String(),
				}).Return(&gen.GetNoteResponse{}, errors.New("rpc error: code = Unknown desc = error"))
			},
			expectedResponse: models.Note{},
		},

		{
			requestBody:    []byte("{\"tag_name\":\"tag\"}"),
			name:           "Test_DeleteTag_BadRequest_gteNoteErr",
			expectedStatus: http.StatusInternalServerError,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().DeleteTag(gomock.Any(), &gen.TagRequest{
					TagName: "tag",
					NoteId:  noteId.String(),
					UserId:  userId.String(),
				}).Return(&gen.GetNoteResponse{Note: nil}, nil)
			},
			expectedResponse: models.Note{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mock_grpc.NewMockNoteClient(ctrl)
			mockAuthClient := mock_auth.NewMockAuthClient(ctrl)
			mockHub := mock_hub.NewMockHubInterface(ctrl)

			defer ctrl.Finish()

			r := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewReader(tt.requestBody))
			ctx := context.Background()
			if tt.expectedStatus != http.StatusUnauthorized {
				ctx = context.WithValue(r.Context(), config.PayloadContextKey, models.JwtPayload{
					Id:       userId,
					Username: "username",
				})
			}
			w := httptest.NewRecorder()
			r = r.WithContext(ctx)
			if tt.name == "Test_DeleteTag_BadRequest" {
				r = mux.SetURLVars(r, map[string]string{"id": ""})

			} else {
				r = mux.SetURLVars(r, map[string]string{"id": noteId.String()})

			}

			handler := CreateNotesHandler(mockClient, mockAuthClient, mockHub)
			tt.mockUsecase(mockClient, mockAuthClient, mockHub)
			handler.DeleteTag(w, r)

			data, _ := json.Marshal(tt.expectedResponse)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, data, w.Body.Bytes())
			}
		})
	}
}

func TestNoteHandler_CreateSubNote(t *testing.T) {
	userId := uuid.NewV4()
	noteId := uuid.NewV4()
	childId := uuid.NewV4()

	tests := []struct {
		requestBody      []byte
		name             string
		expectedStatus   int
		mockUsecase      func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface)
		expectedResponse models.Note
	}{
		{
			requestBody:    []byte("{\"data\":\"\"}"),
			name:           "Test_CreateSubNote_Success",
			expectedStatus: http.StatusOK,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().CreateSubNote(gomock.Any(), &gen.CreateSubNoteRequest{
					UserId:   userId.String(),
					NoteData: "\"\"",
					ParentId: noteId.String(),
				}).Return(&gen.CreateSubNoteResponse{
					Note: &gen.NoteModel{
						Id:            childId.String(),
						OwnerId:       userId.String(),
						Tags:          []string{},
						CreateTime:    time.Time{}.String(),
						UpdateTime:    time.Time{}.String(),
						Parent:        noteId.String(),
						Collaborators: []string{},
						Children:      []string{},
					},
				}, nil)
			},
			expectedResponse: models.Note{
				Id:            childId,
				OwnerId:       userId,
				UpdateTime:    time.Time{},
				CreateTime:    time.Time{},
				Parent:        noteId,
				Data:          "",
				Children:      []uuid.UUID{},
				Tags:          []string{},
				Collaborators: []uuid.UUID{},
			},
		},
		{
			requestBody:    []byte(""),
			name:           "Test_CreateSubNote_Unauthorized",
			expectedStatus: http.StatusUnauthorized,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
			expectedResponse: models.Note{},
		},
		{
			requestBody:    []byte(""),
			name:           "Test_CreateSubnote_BadRequest",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
			expectedResponse: models.Note{},
		},
		{
			requestBody:    []byte("dsfdgf"),
			name:           "Test_CreateSubnote_BadRequest2",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
			expectedResponse: models.Note{},
		},
		{
			requestBody:    []byte("{\"data\":\"\"}"),
			name:           "Test_CreateSubnote_NotFound",
			expectedStatus: http.StatusNotFound,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().CreateSubNote(gomock.Any(), &gen.CreateSubNoteRequest{
					UserId:   userId.String(),
					NoteData: "\"\"",
					ParentId: noteId.String(),
				}).Return(&gen.CreateSubNoteResponse{}, errors.New("rpc error: code = Unknown desc = error"))
			},
			expectedResponse: models.Note{},
		},
		{
			requestBody:    []byte("{\"data\":\"\"}"),
			name:           "Test_CreateSubnote_TooManySubnotes",
			expectedStatus: http.StatusConflict,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().CreateSubNote(gomock.Any(), &gen.CreateSubNoteRequest{
					UserId:   userId.String(),
					NoteData: "\"\"",
					ParentId: noteId.String(),
				}).Return(&gen.CreateSubNoteResponse{}, errors.New(RpcErrorPrefix+note.ErrTooManySubnotes))
			},
			expectedResponse: models.Note{},
		},
		{
			requestBody:    []byte("{\"data\":\"\"}"),
			name:           "Test_CreateSubnote_TooDeep",
			expectedStatus: http.StatusNotFound,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().CreateSubNote(gomock.Any(), &gen.CreateSubNoteRequest{
					UserId:   userId.String(),
					NoteData: "\"\"",
					ParentId: noteId.String(),
				}).Return(&gen.CreateSubNoteResponse{}, errors.New(RpcErrorPrefix+note.ErrTooDeep))
			},
			expectedResponse: models.Note{},
		},

		{
			requestBody:    []byte("{\"data\":\"\"}"),
			name:           "Test_DeleteTag_BadRequest_getNoteErr",
			expectedStatus: http.StatusInternalServerError,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().CreateSubNote(gomock.Any(), &gen.CreateSubNoteRequest{
					UserId:   userId.String(),
					NoteData: "\"\"",
					ParentId: noteId.String(),
				}).Return(&gen.CreateSubNoteResponse{Note: nil}, nil)
			},
			expectedResponse: models.Note{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mock_grpc.NewMockNoteClient(ctrl)
			mockAuthClient := mock_auth.NewMockAuthClient(ctrl)
			mockHub := mock_hub.NewMockHubInterface(ctrl)

			defer ctrl.Finish()

			r := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewReader(tt.requestBody))
			ctx := context.Background()
			if tt.expectedStatus != http.StatusUnauthorized {
				ctx = context.WithValue(r.Context(), config.PayloadContextKey, models.JwtPayload{
					Id:       userId,
					Username: "username",
				})
			}
			w := httptest.NewRecorder()
			r = r.WithContext(ctx)
			if tt.name == "Test_CreateSubnote_BadRequest" {
				r = mux.SetURLVars(r, map[string]string{"id": ""})

			} else {
				r = mux.SetURLVars(r, map[string]string{"id": noteId.String()})

			}

			handler := CreateNotesHandler(mockClient, mockAuthClient, mockHub)
			tt.mockUsecase(mockClient, mockAuthClient, mockHub)
			handler.CreateSubNote(w, r)

			data, _ := json.Marshal(tt.expectedResponse)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, data, w.Body.Bytes())
			}
		})
	}
}

func TestNoteHandler_UpdateNote(t *testing.T) {
	userId := uuid.NewV4()
	noteId := uuid.NewV4()

	tests := []struct {
		requestBody      []byte
		name             string
		expectedStatus   int
		mockUsecase      func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface)
		expectedResponse models.Note
	}{
		{
			requestBody:    []byte("{\"data\":\"new\"}"),
			name:           "Test_UpdateNote_Success",
			expectedStatus: http.StatusOK,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().UpdateNote(gomock.Any(), &gen.UpdateNoteRequest{
					UserId: userId.String(),
					Data:   "\"new\"",
					Id:     noteId.String(),
				}).Return(&gen.UpdateNoteResponse{
					Note: &gen.NoteModel{
						Id:            noteId.String(),
						OwnerId:       userId.String(),
						Tags:          []string{},
						CreateTime:    time.Time{}.String(),
						UpdateTime:    time.Time{}.String(),
						Parent:        uuid.UUID{}.String(),
						Collaborators: []string{},
						Children:      []string{},
						Data:          "\"new\"",
					},
				}, nil)
				mockHub.EXPECT().WriteToCache(gomock.Any(), gomock.Any())

			},
			expectedResponse: models.Note{
				Id:            noteId,
				OwnerId:       userId,
				UpdateTime:    time.Time{},
				CreateTime:    time.Time{},
				Parent:        uuid.UUID{},
				Data:          "new",
				Children:      []uuid.UUID{},
				Tags:          []string{},
				Collaborators: []uuid.UUID{},
			},
		},
		{
			requestBody:    []byte(""),
			name:           "Test_UpdateNote_Unauthorized",
			expectedStatus: http.StatusUnauthorized,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
			expectedResponse: models.Note{},
		},
		{
			requestBody:    []byte(""),
			name:           "Test_UpdateNote_BadRequest",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
			expectedResponse: models.Note{},
		},

		{
			requestBody:    []byte("{\"data\":\"\"}"),
			name:           "Test_UpdateNote_NotFound",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().UpdateNote(gomock.Any(), &gen.UpdateNoteRequest{
					UserId: userId.String(),
					Data:   "\"\"",
					Id:     noteId.String(),
				}).Return(&gen.UpdateNoteResponse{}, errors.New("rpc error: code = Unknown desc = error"))
			},
			expectedResponse: models.Note{},
		},

		{
			requestBody:    []byte("{\"data\":\"\"}"),
			name:           "Test_UpdateNote_BadRequest_getNoteErr",
			expectedStatus: http.StatusInternalServerError,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().UpdateNote(gomock.Any(), &gen.UpdateNoteRequest{
					UserId: userId.String(),
					Data:   "\"\"",
					Id:     noteId.String(),
				}).Return(&gen.UpdateNoteResponse{Note: nil}, nil)
			},
			expectedResponse: models.Note{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mock_grpc.NewMockNoteClient(ctrl)
			mockAuthClient := mock_auth.NewMockAuthClient(ctrl)
			mockHub := mock_hub.NewMockHubInterface(ctrl)

			defer ctrl.Finish()

			r := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewReader(tt.requestBody))
			ctx := context.Background()
			if tt.expectedStatus != http.StatusUnauthorized {
				ctx = context.WithValue(r.Context(), config.PayloadContextKey, models.JwtPayload{
					Id:       userId,
					Username: "username",
				})
			}
			w := httptest.NewRecorder()
			r = r.WithContext(ctx)
			if tt.name == "Test_UpdateNote_BadRequest" {
				r = mux.SetURLVars(r, map[string]string{"id": ""})

			} else {
				r = mux.SetURLVars(r, map[string]string{"id": noteId.String()})

			}

			handler := CreateNotesHandler(mockClient, mockAuthClient, mockHub)
			tt.mockUsecase(mockClient, mockAuthClient, mockHub)
			handler.UpdateNote(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)

		})
	}
}

func TestNoteHandler_AddCollaborator(t *testing.T) {
	userId := uuid.NewV4()
	noteId := uuid.NewV4()
	guestId := uuid.NewV4()

	tests := []struct {
		requestBody    []byte
		name           string
		expectedStatus int
		mockUsecase    func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface)
	}{
		{
			requestBody:    []byte("{\"username\":\"guestuser\"}"),
			name:           "Test_AddCollaborator_Success",
			expectedStatus: http.StatusNoContent,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockAuth.EXPECT().GetUserByUsername(gomock.Any(), &authGen.GetUserByUsernameRequest{
					Username: "guestuser",
				}).Return(&authGen.User{
					Id:         guestId.String(),
					Username:   "guestuser",
					CreateTime: time.Time{}.String(),
				}, nil)
				mockClient.EXPECT().AddCollaborator(gomock.Any(), &gen.AddCollaboratorRequest{
					NoteId:  noteId.String(),
					UserId:  userId.String(),
					GuestId: guestId.String(),
				}).Return(&gen.AddCollaboratorResponse{}, nil)
				mockHub.EXPECT().WriteToCacheMain(gomock.Any(), gomock.Any(), gomock.Any())
			},
		},
		{
			requestBody:    []byte(""),
			name:           "Test_AddCollaborator_Unauthorized",
			expectedStatus: http.StatusUnauthorized,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},
		{
			requestBody:    []byte(""),
			name:           "Test_AddCollaborator_BadRequest",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},

		{
			requestBody:    []byte("{\"username\":\"guestuser\"}"),
			name:           "Test_AddCollaborator_NotFound",
			expectedStatus: http.StatusNotFound,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockAuth.EXPECT().GetUserByUsername(gomock.Any(), &authGen.GetUserByUsernameRequest{
					Username: "guestuser",
				}).Return(&authGen.User{
					Id:         guestId.String(),
					Username:   "guestuser",
					CreateTime: time.Time{}.String(),
				}, nil)
				mockClient.EXPECT().AddCollaborator(gomock.Any(), &gen.AddCollaboratorRequest{
					NoteId:  noteId.String(),
					UserId:  userId.String(),
					GuestId: guestId.String(),
				}).Return(&gen.AddCollaboratorResponse{}, errors.New("rpc error: code = Unknown desc = error"))
			},
		},
		{
			requestBody:    []byte("{\"username\":\"guestuser\"}"),
			name:           "Test_AddCollaborator_GetUserError",
			expectedStatus: http.StatusNotFound,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {

				mockAuth.EXPECT().GetUserByUsername(gomock.Any(), &authGen.GetUserByUsernameRequest{
					Username: "guestuser",
				}).Return(&authGen.User{}, errors.New("rpc error: code = Unknown desc = error"))

			},
		},
		{
			requestBody:    []byte("{\"username\":\"username\"}"),
			name:           "Test_AddCollaborator_AddedHimselfErr",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockAuth.EXPECT().GetUserByUsername(gomock.Any(), &authGen.GetUserByUsernameRequest{
					Username: "username",
				}).Return(&authGen.User{
					Id:         userId.String(),
					Username:   "username",
					CreateTime: time.Time{}.String(),
				}, nil)

			},
		},
		{
			requestBody:    []byte("{\"username\":\"guestuser\"}"),
			name:           "Test_AddCollaborator_AlreadyCollaborator",
			expectedStatus: http.StatusConflict,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockAuth.EXPECT().GetUserByUsername(gomock.Any(), &authGen.GetUserByUsernameRequest{
					Username: "guestuser",
				}).Return(&authGen.User{
					Id:         guestId.String(),
					Username:   "guestuser",
					CreateTime: time.Time{}.String(),
				}, nil)
				mockClient.EXPECT().AddCollaborator(gomock.Any(), &gen.AddCollaboratorRequest{
					NoteId:  noteId.String(),
					UserId:  userId.String(),
					GuestId: guestId.String(),
				}).Return(&gen.AddCollaboratorResponse{}, errors.New(RpcErrorPrefix+note.ErrAlreadyCollaborator))
			},
		},
		{
			requestBody:    []byte("{\"username\":\"guestuser\"}"),
			name:           "Test_AddCollaborator_TooManyCollaboratorsErr",
			expectedStatus: http.StatusExpectationFailed,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockAuth.EXPECT().GetUserByUsername(gomock.Any(), &authGen.GetUserByUsernameRequest{
					Username: "guestuser",
				}).Return(&authGen.User{
					Id:         guestId.String(),
					Username:   "guestuser",
					CreateTime: time.Time{}.String(),
				}, nil)
				mockClient.EXPECT().AddCollaborator(gomock.Any(), &gen.AddCollaboratorRequest{
					NoteId:  noteId.String(),
					UserId:  userId.String(),
					GuestId: guestId.String(),
				}).Return(&gen.AddCollaboratorResponse{}, errors.New(RpcErrorPrefix+note.ErrTooManyCollaborators))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mock_grpc.NewMockNoteClient(ctrl)
			mockAuthClient := mock_auth.NewMockAuthClient(ctrl)
			mockHub := mock_hub.NewMockHubInterface(ctrl)

			defer ctrl.Finish()

			r := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewReader(tt.requestBody))
			ctx := context.Background()
			if tt.expectedStatus != http.StatusUnauthorized {
				ctx = context.WithValue(r.Context(), config.PayloadContextKey, models.JwtPayload{
					Id:       userId,
					Username: "username",
				})
			}
			w := httptest.NewRecorder()
			r = r.WithContext(ctx)
			if tt.name == "Test_AddCollaborator_BadRequest" {
				r = mux.SetURLVars(r, map[string]string{"id": ""})

			} else {
				r = mux.SetURLVars(r, map[string]string{"id": noteId.String()})

			}

			handler := CreateNotesHandler(mockClient, mockAuthClient, mockHub)
			tt.mockUsecase(mockClient, mockAuthClient, mockHub)
			handler.AddCollaborator(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)

		})
	}
}

func TestExportToPDF(t *testing.T) {
	_ = os.Setenv("ATTACHES_BASE_PATH", "/mnt/c/projects/Go/YouNote_data/attaches")

	t.Run("TestExportToPDF_Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockClient := mock_grpc.NewMockNoteClient(ctrl)
		mockAuthClient := mock_auth.NewMockAuthClient(ctrl)
		mockHub := mock_hub.NewMockHubInterface(ctrl)
		defer ctrl.Finish()

		handler := CreateNotesHandler(mockClient, mockAuthClient, mockHub)

		req, err := http.NewRequest("POST", "/export_to_pdf", bytes.NewBufferString(exportpdf.TestNoteHTMLInput))
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(context.Background(), config.PayloadContextKey, models.JwtPayload{
			Id:       uuid.NewV4(),
			Username: "username",
		})
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		r := mux.NewRouter()
		r.Use(protection.ReadAndCloseBody)
		r.Handle("/export_to_pdf", http.HandlerFunc(handler.ExportToPDF)).Methods(http.MethodPost, http.MethodOptions)
		http.HandlerFunc(handler.ExportToPDF).ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		contentType := rr.Header().Get("Content-Type")
		if contentType != "application/pdf" {
			t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/pdf")
		}

		err = os.WriteFile("exported.pdf", rr.Body.Bytes(), 0644)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestNoteHandler_RememberTag(t *testing.T) {
	userId := uuid.NewV4()
	noteId := uuid.NewV4()

	tests := []struct {
		requestBody    []byte
		name           string
		expectedStatus int
		mockUsecase    func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface)
	}{
		{
			requestBody:    []byte("{\"tag_name\":\"tag\"}"),
			name:           "Test_RememberTag_Success",
			expectedStatus: http.StatusNoContent,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().RememberTag(gomock.Any(), &gen.AllTagRequest{
					TagName: "tag",
					UserId:  userId.String(),
				}).Return(&gen.EmptyResponse{}, nil)
			},
		},
		{
			requestBody:    []byte(""),
			name:           "Test_RememberTag_Unauthorized",
			expectedStatus: http.StatusUnauthorized,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},
		{
			requestBody:    []byte(""),
			name:           "Test_RememberTag_BadRequest",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},
		{
			requestBody:    []byte("{\"tag_name\":\"tag\"}"),
			name:           "Test_RememberTag_BadRequestClientErr",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().RememberTag(gomock.Any(), &gen.AllTagRequest{
					TagName: "tag",

					UserId: userId.String(),
				}).Return(&gen.EmptyResponse{}, errors.New("rpc error: code = Unknown desc = error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mock_grpc.NewMockNoteClient(ctrl)
			mockAuthClient := mock_auth.NewMockAuthClient(ctrl)
			mockHub := mock_hub.NewMockHubInterface(ctrl)

			defer ctrl.Finish()

			r := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewReader(tt.requestBody))
			ctx := context.Background()
			if tt.expectedStatus != http.StatusUnauthorized {
				ctx = context.WithValue(r.Context(), config.PayloadContextKey, models.JwtPayload{
					Id:       userId,
					Username: "username",
				})
			}
			w := httptest.NewRecorder()
			r = r.WithContext(ctx)
			if tt.name == "Test_RememberTag_BadRequest" {
				r = mux.SetURLVars(r, map[string]string{"id": ""})

			} else {
				r = mux.SetURLVars(r, map[string]string{"id": noteId.String()})

			}

			handler := CreateNotesHandler(mockClient, mockAuthClient, mockHub)
			tt.mockUsecase(mockClient, mockAuthClient, mockHub)
			handler.RememberTag(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)

		})
	}
}

func TestNoteHandler_ForgetTag(t *testing.T) {
	userId := uuid.NewV4()
	noteId := uuid.NewV4()

	tests := []struct {
		requestBody    []byte
		name           string
		expectedStatus int
		mockUsecase    func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface)
	}{
		{
			requestBody:    []byte("{\"tag_name\":\"tag\"}"),
			name:           "Test_ForgetTag_Success",
			expectedStatus: http.StatusNoContent,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().ForgetTag(gomock.Any(), &gen.AllTagRequest{
					TagName: "tag",
					UserId:  userId.String(),
				}).Return(&gen.EmptyResponse{}, nil)
			},
		},
		{
			requestBody:    []byte(""),
			name:           "Test_ForgetTag_Unauthorized",
			expectedStatus: http.StatusUnauthorized,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},
		{
			requestBody:    []byte(""),
			name:           "Test_ForgetTag_BadRequest",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},
		{
			requestBody:    []byte("{\"tag_name\":\"tag\"}"),
			name:           "Test_ForgetTag_BadRequestClientErr",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().ForgetTag(gomock.Any(), &gen.AllTagRequest{
					TagName: "tag",

					UserId: userId.String(),
				}).Return(&gen.EmptyResponse{}, errors.New("rpc error: code = Unknown desc = error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mock_grpc.NewMockNoteClient(ctrl)
			mockAuthClient := mock_auth.NewMockAuthClient(ctrl)
			mockHub := mock_hub.NewMockHubInterface(ctrl)

			defer ctrl.Finish()

			r := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewReader(tt.requestBody))
			ctx := context.Background()
			if tt.expectedStatus != http.StatusUnauthorized {
				ctx = context.WithValue(r.Context(), config.PayloadContextKey, models.JwtPayload{
					Id:       userId,
					Username: "username",
				})
			}
			w := httptest.NewRecorder()
			r = r.WithContext(ctx)
			if tt.name == "Test_ForgetTag_BadRequest" {
				r = mux.SetURLVars(r, map[string]string{"id": ""})

			} else {
				r = mux.SetURLVars(r, map[string]string{"id": noteId.String()})

			}

			handler := CreateNotesHandler(mockClient, mockAuthClient, mockHub)
			tt.mockUsecase(mockClient, mockAuthClient, mockHub)
			handler.ForgetTag(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)

		})
	}
}

func TestNoteHandler_UpdateTag(t *testing.T) {
	userId := uuid.NewV4()
	noteId := uuid.NewV4()

	tests := []struct {
		requestBody    []byte
		name           string
		expectedStatus int
		mockUsecase    func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface)
	}{
		{
			requestBody:    []byte("{\"old_name\":\"old\",\"new_name\":\"new\"}"),
			name:           "Test_UpdateTag_Success",
			expectedStatus: http.StatusNoContent,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().UpdateTag(gomock.Any(), &gen.UpdateTagRequest{
					OldTag: "old",
					NewTag: "new",
					UserId: userId.String(),
				}).Return(&gen.EmptyResponse{}, nil)
			},
		},
		{
			requestBody:    []byte(""),
			name:           "Test_UpdateTag_Unauthorized",
			expectedStatus: http.StatusUnauthorized,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},
		{
			requestBody:    []byte(""),
			name:           "Test_UpdateTag_BadRequest",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},
		{
			requestBody:    []byte("{\"old_name\":\"old\",\"new_name\":\"new1234567891234567890\"}"),
			name:           "Test_UpdateTag_Invalid_Request",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},
		{
			requestBody:    []byte("{\"old_name\":\"old\",\"new_name\":\"new\"}"),
			name:           "Test_UpdateTag_BadRequestClientErr",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().UpdateTag(gomock.Any(), &gen.UpdateTagRequest{
					OldTag: "old",
					NewTag: "new",

					UserId: userId.String(),
				}).Return(&gen.EmptyResponse{}, errors.New("rpc error: code = Unknown desc = error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mock_grpc.NewMockNoteClient(ctrl)
			mockAuthClient := mock_auth.NewMockAuthClient(ctrl)
			mockHub := mock_hub.NewMockHubInterface(ctrl)

			defer ctrl.Finish()

			r := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewReader(tt.requestBody))
			ctx := context.Background()
			if tt.expectedStatus != http.StatusUnauthorized {
				ctx = context.WithValue(r.Context(), config.PayloadContextKey, models.JwtPayload{
					Id:       userId,
					Username: "username",
				})
			}
			w := httptest.NewRecorder()
			r = r.WithContext(ctx)
			if tt.name == "Test_UpdateTag_BadRequest" {
				r = mux.SetURLVars(r, map[string]string{"id": ""})

			} else {
				r = mux.SetURLVars(r, map[string]string{"id": noteId.String()})

			}

			handler := CreateNotesHandler(mockClient, mockAuthClient, mockHub)
			tt.mockUsecase(mockClient, mockAuthClient, mockHub)
			handler.UpdateTag(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)

		})
	}
}

func TestNoteHandler_SetIcon(t *testing.T) {
	userId := uuid.NewV4()
	noteId := uuid.NewV4()

	tests := []struct {
		requestBody    []byte
		name           string
		expectedStatus int
		mockUsecase    func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface)
	}{
		{
			requestBody:    []byte("{\"icon\":\"icon\"}"),
			name:           "Test_SetIcon_Success",
			expectedStatus: http.StatusOK,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().SetIcon(gomock.Any(), &gen.SetIconRequest{
					NoteId: noteId.String(),
					Icon:   "icon",
					UserId: userId.String(),
				}).Return(&gen.GetNoteResponse{
					Note: &gen.NoteModel{
						Id:            noteId.String(),
						OwnerId:       userId.String(),
						Tags:          []string{},
						CreateTime:    time.Time{}.String(),
						UpdateTime:    time.Time{}.String(),
						Parent:        uuid.UUID{}.String(),
						Collaborators: []string{},
						Children:      []string{},
						Icon:          "icon",
					},
				}, nil)
			},
		},
		{
			requestBody:    []byte(""),
			name:           "Test_SetIcon_Unauthorized",
			expectedStatus: http.StatusUnauthorized,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},
		{
			requestBody:    []byte(""),
			name:           "Test_SetIcon_BadRequest",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},
		{
			requestBody:    []byte("{\"icon\":\"icon\"}"),
			name:           "Test_SetIcon_BadRequestClientErr",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().SetIcon(gomock.Any(), &gen.SetIconRequest{
					NoteId: noteId.String(),
					Icon:   "icon",
					UserId: userId.String(),
				}).Return(&gen.GetNoteResponse{}, errors.New("rpc error: code = Unknown desc = error"))
			},
		},
		{
			requestBody:    []byte("{\"icon\":\"icon\"}"),
			name:           "Test_SetIcon_GetNoteErr",
			expectedStatus: http.StatusInternalServerError,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().SetIcon(gomock.Any(), &gen.SetIconRequest{
					NoteId: noteId.String(),
					Icon:   "icon",
					UserId: userId.String(),
				}).Return(&gen.GetNoteResponse{Note: nil}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mock_grpc.NewMockNoteClient(ctrl)
			mockAuthClient := mock_auth.NewMockAuthClient(ctrl)
			mockHub := mock_hub.NewMockHubInterface(ctrl)

			defer ctrl.Finish()

			r := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewReader(tt.requestBody))
			ctx := context.Background()
			if tt.expectedStatus != http.StatusUnauthorized {
				ctx = context.WithValue(r.Context(), config.PayloadContextKey, models.JwtPayload{
					Id:       userId,
					Username: "username",
				})
			}
			w := httptest.NewRecorder()
			r = r.WithContext(ctx)
			if tt.name == "Test_SetIcon_BadRequest" {
				r = mux.SetURLVars(r, map[string]string{"id": ""})

			} else {
				r = mux.SetURLVars(r, map[string]string{"id": noteId.String()})

			}

			handler := CreateNotesHandler(mockClient, mockAuthClient, mockHub)
			tt.mockUsecase(mockClient, mockAuthClient, mockHub)
			handler.SetIcon(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)

		})
	}
}

func TestNoteHandler_SetHeader(t *testing.T) {
	userId := uuid.NewV4()
	noteId := uuid.NewV4()

	tests := []struct {
		requestBody    []byte
		name           string
		expectedStatus int
		mockUsecase    func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface)
	}{
		{
			requestBody:    []byte("{\"header\":\"header\"}"),
			name:           "Test_SetHeader_Success",
			expectedStatus: http.StatusOK,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().SetHeader(gomock.Any(), &gen.SetHeaderRequest{
					NoteId: noteId.String(),
					Header: "header",
					UserId: userId.String(),
				}).Return(&gen.GetNoteResponse{
					Note: &gen.NoteModel{
						Id:            noteId.String(),
						OwnerId:       userId.String(),
						Tags:          []string{},
						CreateTime:    time.Time{}.String(),
						UpdateTime:    time.Time{}.String(),
						Parent:        uuid.UUID{}.String(),
						Collaborators: []string{},
						Children:      []string{},
						Header:        "header",
					},
				}, nil)
			},
		},
		{
			requestBody:    []byte(""),
			name:           "Test_SetHeader_Unauthorized",
			expectedStatus: http.StatusUnauthorized,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},
		{
			requestBody:    []byte(""),
			name:           "Test_SetHeader_BadRequest",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},
		{
			requestBody:    []byte("{\"header\":\"header\"}"),
			name:           "Test_SetHeader_BadRequestClientErr",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().SetHeader(gomock.Any(), &gen.SetHeaderRequest{
					NoteId: noteId.String(),
					Header: "header",
					UserId: userId.String(),
				}).Return(&gen.GetNoteResponse{}, errors.New("rpc error: code = Unknown desc = error"))
			},
		},
		{
			requestBody:    []byte("{\"header\":\"header\"}"),
			name:           "Test_SetHeader_GetNoteErr",
			expectedStatus: http.StatusInternalServerError,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().SetHeader(gomock.Any(), &gen.SetHeaderRequest{
					NoteId: noteId.String(),
					Header: "header",
					UserId: userId.String(),
				}).Return(&gen.GetNoteResponse{Note: nil}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mock_grpc.NewMockNoteClient(ctrl)
			mockAuthClient := mock_auth.NewMockAuthClient(ctrl)
			mockHub := mock_hub.NewMockHubInterface(ctrl)

			defer ctrl.Finish()

			r := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewReader(tt.requestBody))
			ctx := context.Background()
			if tt.expectedStatus != http.StatusUnauthorized {
				ctx = context.WithValue(r.Context(), config.PayloadContextKey, models.JwtPayload{
					Id:       userId,
					Username: "username",
				})
			}
			w := httptest.NewRecorder()
			r = r.WithContext(ctx)
			if tt.name == "Test_SetHeader_BadRequest" {
				r = mux.SetURLVars(r, map[string]string{"id": ""})

			} else {
				r = mux.SetURLVars(r, map[string]string{"id": noteId.String()})

			}

			handler := CreateNotesHandler(mockClient, mockAuthClient, mockHub)
			tt.mockUsecase(mockClient, mockAuthClient, mockHub)
			handler.SetHeader(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)

		})
	}
}

func TestNoteHandler_AddFavorite(t *testing.T) {
	userId := uuid.NewV4()
	noteId := uuid.NewV4()

	tests := []struct {
		name           string
		expectedStatus int
		mockUsecase    func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface)
	}{
		{
			name:           "Test_AddFavorite_Success",
			expectedStatus: http.StatusOK,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().AddFav(gomock.Any(), &gen.ChangeFlagRequest{
					NoteId: noteId.String(),
					UserId: userId.String(),
				}).Return(&gen.GetNoteResponse{
					Note: &gen.NoteModel{
						Id:            noteId.String(),
						OwnerId:       userId.String(),
						Tags:          []string{},
						CreateTime:    time.Time{}.String(),
						UpdateTime:    time.Time{}.String(),
						Parent:        uuid.UUID{}.String(),
						Collaborators: []string{},
						Children:      []string{},
					},
				}, nil)
			},
		},
		{
			name:           "Test_AddFavorite_Unauthorized",
			expectedStatus: http.StatusUnauthorized,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},
		{
			name:           "Test_AddFavorite_BadRequest",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},
		{
			name:           "Test_AddFavorite_BadRequestClientErr",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().AddFav(gomock.Any(), &gen.ChangeFlagRequest{
					NoteId: noteId.String(),
					UserId: userId.String(),
				}).Return(&gen.GetNoteResponse{}, errors.New("rpc error: code = Unknown desc = error"))
			},
		},
		{
			name:           "Test_AddFavorite_GetNoteErr",
			expectedStatus: http.StatusInternalServerError,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().AddFav(gomock.Any(), &gen.ChangeFlagRequest{
					NoteId: noteId.String(),
					UserId: userId.String(),
				}).Return(&gen.GetNoteResponse{Note: nil}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mock_grpc.NewMockNoteClient(ctrl)
			mockAuthClient := mock_auth.NewMockAuthClient(ctrl)
			mockHub := mock_hub.NewMockHubInterface(ctrl)

			defer ctrl.Finish()

			r := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewReader([]byte{}))
			ctx := context.Background()
			if tt.expectedStatus != http.StatusUnauthorized {
				ctx = context.WithValue(r.Context(), config.PayloadContextKey, models.JwtPayload{
					Id:       userId,
					Username: "username",
				})
			}
			w := httptest.NewRecorder()
			r = r.WithContext(ctx)
			if tt.name == "Test_AddFavorite_BadRequest" {
				r = mux.SetURLVars(r, map[string]string{"id": ""})

			} else {
				r = mux.SetURLVars(r, map[string]string{"id": noteId.String()})

			}

			handler := CreateNotesHandler(mockClient, mockAuthClient, mockHub)
			tt.mockUsecase(mockClient, mockAuthClient, mockHub)
			handler.AddFavorite(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)

		})
	}
}

func TestNoteHandler_DeleteFavorite(t *testing.T) {
	userId := uuid.NewV4()
	noteId := uuid.NewV4()

	tests := []struct {
		name           string
		expectedStatus int
		mockUsecase    func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface)
	}{
		{
			name:           "Test_DeleteFavorite_Success",
			expectedStatus: http.StatusOK,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().DelFav(gomock.Any(), &gen.ChangeFlagRequest{
					NoteId: noteId.String(),
					UserId: userId.String(),
				}).Return(&gen.GetNoteResponse{
					Note: &gen.NoteModel{
						Id:            noteId.String(),
						OwnerId:       userId.String(),
						Tags:          []string{},
						CreateTime:    time.Time{}.String(),
						UpdateTime:    time.Time{}.String(),
						Parent:        uuid.UUID{}.String(),
						Collaborators: []string{},
						Children:      []string{},
					},
				}, nil)
			},
		},
		{
			name:           "Test_DeleteFavorite_Unauthorized",
			expectedStatus: http.StatusUnauthorized,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},
		{
			name:           "Test_DeleteFavorite_BadRequest",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},
		{
			name:           "Test_DeleteFavorite_BadRequestClientErr",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().DelFav(gomock.Any(), &gen.ChangeFlagRequest{
					NoteId: noteId.String(),
					UserId: userId.String(),
				}).Return(&gen.GetNoteResponse{}, errors.New("rpc error: code = Unknown desc = error"))
			},
		},
		{
			name:           "Test_DeleteFavorite_GetNoteErr",
			expectedStatus: http.StatusInternalServerError,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().DelFav(gomock.Any(), &gen.ChangeFlagRequest{
					NoteId: noteId.String(),
					UserId: userId.String(),
				}).Return(&gen.GetNoteResponse{Note: nil}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mock_grpc.NewMockNoteClient(ctrl)
			mockAuthClient := mock_auth.NewMockAuthClient(ctrl)
			mockHub := mock_hub.NewMockHubInterface(ctrl)

			defer ctrl.Finish()

			r := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewReader([]byte{}))
			ctx := context.Background()
			if tt.expectedStatus != http.StatusUnauthorized {
				ctx = context.WithValue(r.Context(), config.PayloadContextKey, models.JwtPayload{
					Id:       userId,
					Username: "username",
				})
			}
			w := httptest.NewRecorder()
			r = r.WithContext(ctx)
			if tt.name == "Test_DeleteFavorite_BadRequest" {
				r = mux.SetURLVars(r, map[string]string{"id": ""})

			} else {
				r = mux.SetURLVars(r, map[string]string{"id": noteId.String()})

			}

			handler := CreateNotesHandler(mockClient, mockAuthClient, mockHub)
			tt.mockUsecase(mockClient, mockAuthClient, mockHub)
			handler.DeleteFavorite(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)

		})
	}
}

func TestNoteHandler_SetPublic(t *testing.T) {
	userId := uuid.NewV4()
	noteId := uuid.NewV4()

	tests := []struct {
		name           string
		expectedStatus int
		mockUsecase    func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface)
	}{
		{
			name:           "Test_SetPublic_Success",
			expectedStatus: http.StatusOK,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().SetPublic(gomock.Any(), &gen.AccessModeRequest{
					NoteId: noteId.String(),
					UserId: userId.String(),
				}).Return(&gen.GetNoteResponse{
					Note: &gen.NoteModel{
						Id:            noteId.String(),
						OwnerId:       userId.String(),
						Tags:          []string{},
						CreateTime:    time.Time{}.String(),
						UpdateTime:    time.Time{}.String(),
						Parent:        uuid.UUID{}.String(),
						Collaborators: []string{},
						Children:      []string{},
					},
				}, nil)
			},
		},
		{
			name:           "Test_SetPublic_Unauthorized",
			expectedStatus: http.StatusUnauthorized,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},
		{
			name:           "Test_SetPublic_BadRequest",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},
		{
			name:           "Test_SetPublic_BadRequestClientErr",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().SetPublic(gomock.Any(), &gen.AccessModeRequest{
					NoteId: noteId.String(),
					UserId: userId.String(),
				}).Return(&gen.GetNoteResponse{}, errors.New("rpc error: code = Unknown desc = error"))
			},
		},
		{
			name:           "Test_SetPublic_GetNoteErr",
			expectedStatus: http.StatusInternalServerError,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().SetPublic(gomock.Any(), &gen.AccessModeRequest{
					NoteId: noteId.String(),
					UserId: userId.String(),
				}).Return(&gen.GetNoteResponse{Note: nil}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mock_grpc.NewMockNoteClient(ctrl)
			mockAuthClient := mock_auth.NewMockAuthClient(ctrl)
			mockHub := mock_hub.NewMockHubInterface(ctrl)

			defer ctrl.Finish()

			r := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewReader([]byte{}))
			ctx := context.Background()
			if tt.expectedStatus != http.StatusUnauthorized {
				ctx = context.WithValue(r.Context(), config.PayloadContextKey, models.JwtPayload{
					Id:       userId,
					Username: "username",
				})
			}
			w := httptest.NewRecorder()
			r = r.WithContext(ctx)
			if tt.name == "Test_SetPublic_BadRequest" {
				r = mux.SetURLVars(r, map[string]string{"id": ""})

			} else {
				r = mux.SetURLVars(r, map[string]string{"id": noteId.String()})

			}

			handler := CreateNotesHandler(mockClient, mockAuthClient, mockHub)
			tt.mockUsecase(mockClient, mockAuthClient, mockHub)
			handler.SetPublic(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)

		})
	}
}

func TestNoteHandler_SetPrivate(t *testing.T) {

	userId := uuid.NewV4()
	noteId := uuid.NewV4()

	tests := []struct {
		name           string
		expectedStatus int
		mockUsecase    func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface)
	}{
		{
			name:           "Test_SetPrivate_Success",
			expectedStatus: http.StatusOK,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().SetPrivate(gomock.Any(), &gen.AccessModeRequest{
					NoteId: noteId.String(),
					UserId: userId.String(),
				}).Return(&gen.GetNoteResponse{
					Note: &gen.NoteModel{
						Id:            noteId.String(),
						OwnerId:       userId.String(),
						Tags:          []string{},
						CreateTime:    time.Time{}.String(),
						UpdateTime:    time.Time{}.String(),
						Parent:        uuid.UUID{}.String(),
						Collaborators: []string{},
						Children:      []string{},
					},
				}, nil)
			},
		},
		{
			name:           "Test_SetPrivate_Unauthorized",
			expectedStatus: http.StatusUnauthorized,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},
		{
			name:           "Test_SetPrivate_BadRequest",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
			},
		},
		{
			name:           "Test_SetPrivate_BadRequestClientErr",
			expectedStatus: http.StatusBadRequest,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().SetPrivate(gomock.Any(), &gen.AccessModeRequest{
					NoteId: noteId.String(),
					UserId: userId.String(),
				}).Return(&gen.GetNoteResponse{}, errors.New("rpc error: code = Unknown desc = error"))
			},
		},
		{
			name:           "Test_SetPrivate_GetNoteErr",
			expectedStatus: http.StatusInternalServerError,
			mockUsecase: func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface) {
				mockClient.EXPECT().SetPrivate(gomock.Any(), &gen.AccessModeRequest{
					NoteId: noteId.String(),
					UserId: userId.String(),
				}).Return(&gen.GetNoteResponse{Note: nil}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mock_grpc.NewMockNoteClient(ctrl)
			mockAuthClient := mock_auth.NewMockAuthClient(ctrl)
			mockHub := mock_hub.NewMockHubInterface(ctrl)

			defer ctrl.Finish()

			r := httptest.NewRequest("POST", "http://example.com/api/handler", bytes.NewReader([]byte{}))
			ctx := context.Background()
			if tt.expectedStatus != http.StatusUnauthorized {
				ctx = context.WithValue(r.Context(), config.PayloadContextKey, models.JwtPayload{
					Id:       userId,
					Username: "username",
				})
			}
			w := httptest.NewRecorder()
			r = r.WithContext(ctx)
			if tt.name == "Test_SetPrivate_BadRequest" {
				r = mux.SetURLVars(r, map[string]string{"id": ""})

			} else {
				r = mux.SetURLVars(r, map[string]string{"id": noteId.String()})

			}

			handler := CreateNotesHandler(mockClient, mockAuthClient, mockHub)
			tt.mockUsecase(mockClient, mockAuthClient, mockHub)
			handler.SetPrivate(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)

		})
	}
}
