package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	ParseBodyError       = "can`t parse delivery body: "
	WriteBodyError       = "can`t write response body: "
	JwtPayloadParseError = "can`t parse JWT payload from delivery context"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func WriteResponseData(w http.ResponseWriter, responseData interface{}, successStatusCode int) error {
	body, err := json.Marshal(responseData)
	if err != nil {
		return fmt.Errorf("error in marshalling response body: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))

	w.WriteHeader(successStatusCode)
	_, _ = w.Write(body)

	return nil
}

func WriteErrorMessage(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	_, _ = fmt.Fprintf(w, `{"message":"%s"}`, message)
}
