package authmw

import (
	"fmt"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/satori/uuid"
	"net/http"
	"strings"
	"time"
)

const secret = "Только не кому не говори..."

type authService struct {
	Id       uuid.UUID
	Username string
	LifeTime time.Duration
	secret   string
}

func (service *authService) GenToken(user *models.User, lifeTime time.Duration) (string, error) {
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
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		token := headerParts[1]
		claims, err := parseToken(token)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if claims.Claims["exp"] < time.Now().Unix() {
			// ...
		}

		next.ServeHTTP(w, r)
	})
}
