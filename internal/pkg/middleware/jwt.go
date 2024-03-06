package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

const (
	JwtCookie = "YouNoteJWT"
)

func GenToken(user models.User, lifeTime time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.Id,
		"usr": user.Username,
		"exp": time.Now().Add(lifeTime).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func parseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}

func parseJwtPayloadFromClaims(claims *jwt.Token) (models.JwtPayload, error) {
	payloadMap, ok := claims.Claims.(jwt.MapClaims)
	if !ok {
		return models.JwtPayload{}, errors.New("invalid format (claims)")
	}
	stringUserId, ok := payloadMap["id"].(string)
	if !ok {
		return models.JwtPayload{}, errors.New("invalid format (id)")
	}
	username, ok := payloadMap["usr"].(string)
	if !ok {
		return models.JwtPayload{}, errors.New("invalid format (usr)")
	}
	userId, err := uuid.FromString(stringUserId)
	if err != nil {
		return models.JwtPayload{}, errors.New("invalid format (id)")
	}

	return models.JwtPayload{
		Id:       userId,
		Username: username,
	}, nil
}

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := headerParts[1]

		cookie, err := r.Cookie(JwtCookie)
		if err != nil || cookie.Value != token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := parseToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		timeExp, err := claims.Claims.GetExpirationTime()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if timeExp.Before(time.Now().UTC()) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		payload, err := parseJwtPayloadFromClaims(claims)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), models.PayloadContextKey, payload)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
