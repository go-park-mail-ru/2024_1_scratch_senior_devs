package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/authmw"
)

func GetRequestData(r *http.Request, requestData interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		return err
	}

	return nil
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

func GenTokenCookie(token string, expTime time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     authmw.JwtCookie,
		Secure:   false,
		Value:    token,
		HttpOnly: false,
		Expires:  expTime,
	}
}

func DelTokenCookie() *http.Cookie {
	return &http.Cookie{
		Name:   authmw.JwtCookie,
		Value:  "",
		MaxAge: -1,
	}
}

func WriteIncorrectFormatError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(`{"message":"incorrect request body format"}`))
}
