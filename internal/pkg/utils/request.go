package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware"
	"io"
	"net/http"
	"strconv"
	"time"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

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
		Name:     middleware.JwtCookie,
		Secure:   false, // TODO false --> true
		Value:    token,
		HttpOnly: false, // TODO false --> true
		Expires:  expTime,
		Path:     "/",
	}
}

func DelTokenCookie() *http.Cookie {
	return &http.Cookie{
		Name:   middleware.JwtCookie,
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}
}

func WriteErrorMessage(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	_, _ = fmt.Fprintf(w, `{"message":"%s"}`, message)
}

func GetHash(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes)
}
