package middleware

import (
	"log/slog"
	"net/http"
	"yadro-go/internal/adapter/primary/http/hanlder"
	"yadro-go/internal/adapter/primary/http/protocol"
	"yadro-go/internal/core/domain"
)

type RateLimiter interface {
	Take(id string) bool
}

type RpsLimitMiddleware struct {
	log         *slog.Logger
	rateLimiter RateLimiter
}

func NewRpcLimitMiddleware(log *slog.Logger, rateLimiter RateLimiter) *RpsLimitMiddleware {
	return &RpsLimitMiddleware{
		log:         log,
		rateLimiter: rateLimiter,
	}
}

func (rl *RpsLimitMiddleware) WithRpsLimit(next hanlder.AuthenticatedHandlerFunc) hanlder.AuthenticatedHandlerFunc {
	const op = "middleware.WithRpsLimit"
	log := rl.log.With(slog.String("op", op))

	return func(w http.ResponseWriter, request *http.Request, user *domain.User) {
		if !rl.rateLimiter.Take(user.Username) {
			log.Warn("rate limit exceeded", slog.String("uname", user.Username))
			protocol.ResponseError(w, http.StatusTooManyRequests, "request limit exceeded")
			return
		}

		next(w, request, user)
	}
}
