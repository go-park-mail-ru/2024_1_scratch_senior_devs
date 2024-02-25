package authmw

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/satori/uuid"
)

const (
	secret    = "Только никому не говори..."
	JwtCookie = "YouNoteJWT"
)

func GenToken(user *models.User, lifeTime time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.Id,
		"usr": user.Username,
		"exp": time.Now().Add(lifeTime).Unix(),
	})
	return token.SignedString([]byte(secret))
}

func parseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
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

		timeExp, err := claims.Claims.GetExpirationTime() //получаем из токена время просрока
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if timeExp.Before(time.Now().UTC()) { //если токен просрочен

			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		payloadMap := claims.Claims.(jwt.MapClaims)
		payload := models.JwtPayload{
			Id:       payloadMap["id"].(uuid.UUID),
			Username: payloadMap["usr"].(string),
		}
		ctx := context.WithValue(r.Context(), "payload", payload)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
