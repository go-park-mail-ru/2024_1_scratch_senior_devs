package utils

import (
	"bytes"
	"context"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetRequestData(t *testing.T) {
	type dataType struct {
		Key string `json:"key"`
	}

	var tests = []struct {
		name         string
		requestBody  []byte
		expectedData dataType
		err          bool
	}{
		{
			name:         "GetRequestData_Success",
			requestBody:  []byte(`{"key":"value"}`),
			expectedData: dataType{Key: "value"},
			err:          false,
		},
		{
			name:         "GetRequestData_Fail",
			requestBody:  []byte(`{"key":"value"`),
			expectedData: dataType{},
			err:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := http.NewRequestWithContext(context.Background(), "POST", "http://127.0.0.1:8080/test", bytes.NewBuffer(tt.requestBody))
			data := dataType{}

			err := GetRequestData(request, &data)
			assert.Equal(t, tt.expectedData, data)
			if (err != nil) != tt.err {
				t.Error("error in error")
			}
		})
	}
}

func TestWriteResponseData_Success(t *testing.T) {
	type dataType struct {
		Key string `json:"key"`
	}

	var tests = []struct {
		name string
		data dataType
		err  bool
	}{
		{
			name: "WriteResponseData_Success",
			data: dataType{Key: "value"},
			err:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := WriteResponseData(w, tt.data, http.StatusOK)
			if (err != nil) != tt.err {
				t.Error("error in error")
			}
		})
	}
}

func TestWriteResponseData_Fail(t *testing.T) {
	type dataType struct {
		Key func() `json:"key"`
	}

	var tests = []struct {
		name string
		data dataType
		err  bool
	}{
		{
			name: "WriteResponseData_Fail",
			data: dataType{Key: func() {}},
			err:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := WriteResponseData(w, tt.data, http.StatusOK)
			if (err != nil) != tt.err {
				t.Error("error in error")
			}
		})
	}
}

func TestGenTokenCookie(t *testing.T) {
	var tests = []struct {
		name    string
		token   string
		expTime time.Time
	}{
		{
			name:    "GenTokenCookie_Success",
			token:   "abc123",
			expTime: time.Now(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cookie := GenTokenCookie(tt.token, tt.expTime)
			assert.Equal(t, cookie.Name, middleware.JwtCookie)
			assert.Equal(t, cookie.Value, tt.token)
		})
	}
}

func TestDelTokenCookie(t *testing.T) {
	var tests = []struct {
		name string
	}{
		{
			name: "DelTokenCookie_Success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cookie := DelTokenCookie()
			assert.Equal(t, cookie.Name, middleware.JwtCookie)
			assert.Equal(t, cookie.Value, "")
		})
	}
}

func TestWriteErrorMessage(t *testing.T) {
	var tests = []struct {
		name string
	}{
		{
			name: "WriteErrorMessage_Success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteErrorMessage(w, http.StatusOK, "abc")

			assert.Equal(t, w.Code, http.StatusOK)
		})
	}
}

func TestGFN(t *testing.T) {
	GFN()
}
