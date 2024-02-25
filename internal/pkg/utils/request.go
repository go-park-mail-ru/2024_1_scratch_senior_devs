package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
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

func WriteResponseData(w http.ResponseWriter, responseData interface{}) error {
	body, err := json.Marshal(responseData)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))

	_, err = w.Write(body)
	if err != nil {
		return err
	}

	return nil
}

func GenTokenCookie(token string, expTime time.Time) *http.Cookie {
	return &http.Cookie{
		Secure:   true,
		Value:    token,
		HttpOnly: true,
		Expires:  expTime,
	}
}
