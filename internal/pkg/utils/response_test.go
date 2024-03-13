package utils

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
