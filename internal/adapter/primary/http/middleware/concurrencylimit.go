package middleware

import (
	"log/slog"
	"net/http"
	"yadro-go/internal/adapter/primary/http/protocol"
	"yadro-go/pkg/sync"
)

type ConcurrencyLimitMiddleware struct {
	log       *slog.Logger
	semaphore *sync.Semaphore
}

func NewConcurrencyLimitMiddleware(log *slog.Logger, limit int) *ConcurrencyLimitMiddleware {
	return &ConcurrencyLimitMiddleware{
		log:       log,
		semaphore: sync.NewSemaphore(limit),
	}
}

func (c *ConcurrencyLimitMiddleware) WithConcurrencyLimit(next http.Handler) http.Handler {
	const op = "middleware.WithConcurrencyLimit"
	log := c.log.With(slog.String("op", op))

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !c.semaphore.TryAcquire() {
			log.Warn("concurrency limit exceeded")
			protocol.ResponseError(w, http.StatusServiceUnavailable, "service unavailable, try again later")
			return
		}
		defer c.semaphore.Release()

		next.ServeHTTP(w, req)
	})
}
