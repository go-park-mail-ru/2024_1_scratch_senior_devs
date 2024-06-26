package responses

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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
