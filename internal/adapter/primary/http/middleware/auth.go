package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"yadro-go/internal/adapter/primary"
	"yadro-go/internal/adapter/primary/http/protocol"
	"yadro-go/internal/core/service"
	"yadro-go/pkg/logger"
)

type AuthMiddleware struct {
	log  *slog.Logger
	auth primary.Auth
}

func NewAuthMiddleware(log *slog.Logger, auth primary.Auth) *AuthMiddleware {
	return &AuthMiddleware{log: log, auth: auth}
}

func (a *AuthMiddleware) WithAuth(role int, next http.HandlerFunc) http.HandlerFunc {
	const op = "middleware.WithAuth"
	log := a.log.With(slog.String("op", op))

	return func(w http.ResponseWriter, request *http.Request) {
		bearer := strings.Split(request.Header.Get("Authorization"), "Bearer ")
		if len(bearer) != 2 {
			responseUnauthorized(w)
			return
		}

		token := bearer[1]

		user, err := a.auth.Authenticate(request.Context(), token)
		if err != nil {
			if errors.Is(err, service.ErrBadToken) {
				responseUnauthorized(w)
				return
			}

			log.Error("failed to authenticate", logger.Err(err))
			protocol.ResponseError(w, http.StatusInternalServerError, "internal error")
			return
		}

		if !user.HasRole(role) {
			protocol.ResponseError(w, http.StatusForbidden, "forbidden")
			return
		}

		next(w, request)
	}
}

func responseUnauthorized(w http.ResponseWriter) {
	w.Header().Add("WWW-Authenticate", "Bearer")
	protocol.ResponseError(w, http.StatusUnauthorized, "unauthorized")
}
