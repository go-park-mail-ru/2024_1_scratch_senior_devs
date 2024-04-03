package protection

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/gorilla/mux"

	"github.com/golang-jwt/jwt/v5"
	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

func GenJwtToken(user models.User, lifeTime time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.Id,
		"usr": user.Username,
		"exp": time.Now().Add(lifeTime).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func parseJwtToken(token string) (*jwt.Token, error) {
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

func CreateJwtMiddleware(logger *slog.Logger, cfg config.JwtConfig) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwtLogger := logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

			header := r.Header.Get("Authorization")
			if header == "" {
				log.LogHandlerError(jwtLogger, http.StatusUnauthorized, "empty Authorization header")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			headerParts := strings.Split(header, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				log.LogHandlerError(jwtLogger, http.StatusUnauthorized, "wrong Authorization header format")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			token := headerParts[1]

			jwtCookie, err := r.Cookie(cfg.JwtCookie)
			if err != nil {
				log.LogHandlerError(jwtLogger, http.StatusUnauthorized, "no jwt cookie: "+err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if jwtCookie.Value != token {
				log.LogHandlerError(jwtLogger, http.StatusUnauthorized, "tokens in cookie and header are different")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			claims, err := parseJwtToken(token)
			if err != nil {
				log.LogHandlerError(jwtLogger, http.StatusUnauthorized, "invalid token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			timeExp, err := claims.Claims.GetExpirationTime()
			if err != nil {
				log.LogHandlerError(jwtLogger, http.StatusUnauthorized, "invalid token expiration time")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if timeExp.Before(time.Now().UTC()) {
				log.LogHandlerError(jwtLogger, http.StatusUnauthorized, "token expired")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			payload, err := parseJwtPayloadFromClaims(claims)
			if err != nil {
				log.LogHandlerError(jwtLogger, http.StatusUnauthorized, "invalid token payload")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			jwtLogger.Info("success")

			ctx := context.WithValue(r.Context(), config.PayloadContextKey, payload)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func ReadAndCloseBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			_, _ = io.ReadAll(r.Body)
			defer r.Body.Close()
		}()

		next.ServeHTTP(w, r)
	})
}
