package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	mock_auth "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/delivery/grpc/gen/mocks"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	mock_hub "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/hub/mocks"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc/gen"
	mock_grpc "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc/gen/mocks"
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
		expectedData   []*gen.NoteModel
	}{
		{

			name:           successTestName,
			id:             uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
			username:       "testuser",
			expectedStatus: http.StatusOK,
			expectedData: []*gen.NoteModel{
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
			expectedData:   []*gen.NoteModel{},
		},
		{

			name:           "Test Error",
			id:             uuid.FromStringOrNil(""),
			username:       "",
			expectedStatus: http.StatusBadRequest,
			expectedData:   []*gen.NoteModel{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			//mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
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
				for i, note := range tt.expectedData {
					crTime, _ := time.Parse("2006-01-02 15:04:05 -0700 UTC", note.CreateTime)
					upTime, _ := time.Parse("2006-01-02 15:04:05 -0700 UTC", note.UpdateTime)
					expectedResult[i] = models.Note{
						Id:            uuid.FromStringOrNil(note.Id),
						Data:          []byte(note.Data),
						CreateTime:    crTime,
						UpdateTime:    upTime,
						Parent:        uuid.FromStringOrNil(note.OwnerId),
						OwnerId:       uuid.FromStringOrNil(note.OwnerId),
						Children:      []uuid.UUID{},
						Tags:          []string{},
						Collaborators: []uuid.UUID{},
					}

				}
				d, _ := json.Marshal(expectedResult)
				assert.Equal(t, w.Body.Bytes(), d)
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
		expectedData   models.Note
	}{
		// TODO: Add test cases.
		{
			name:           successTestName,
			expectedStatus: 200,
			noteId:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
			userId:         uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
			username:       "test_user",
			expectedData: models.Note{
				Id:            uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
				OwnerId:       uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
				UpdateTime:    time.Time{},
				CreateTime:    time.Time{},
				Data:          []byte(""),
				Children:      []uuid.UUID{},
				Tags:          []string{},
				Collaborators: []uuid.UUID{},
			},
		},
		{
			name:           "Test Unauthorized",
			expectedStatus: 401,
			noteId:         uuid.FromStringOrNil("c80e3ea8-0813-4731-b6ee-b41604c56f95"),
			userId:         uuid.FromStringOrNil(""),
			username:       "test_user",
			expectedData:   models.Note{},
		},
		{
			name:           "Test Error",
			expectedStatus: http.StatusNotFound,
			noteId:         uuid.FromStringOrNil(""),
			userId:         uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
			username:       "test_user",
			expectedData:   models.Note{},
		},
		{
			name:           "Test Bad Request",
			expectedStatus: http.StatusBadRequest,
			noteId:         uuid.FromStringOrNil(""),
			userId:         uuid.FromStringOrNil("a233ea8-0813-4731-b12e-b41604c56f95"),
			username:       "test_user",
			expectedData:   models.Note{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			//mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
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
				}).Return(&gen.GetNoteResponse{Note: &gen.NoteModel{
					Id:            tt.expectedData.Id.String(),
					CreateTime:    tt.expectedData.CreateTime.String(),
					UpdateTime:    tt.expectedData.UpdateTime.String(),
					Data:          string(tt.expectedData.Data),
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
			//mockUsecase := mock_note.NewMockNoteUsecase(ctrl)
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

	tests := []struct {
		name             string
		expectedStatus   int
		mockUsecase      func(mockClient *mock_grpc.MockNoteClient, mockAuth *mock_auth.MockAuthClient, mockHub *mock_hub.MockHubInterface)
		expectedResponse models.GetTagsResponse
	}{}
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
