package httpserver

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"
	"yadro-go/pkg/logger"
)

const (
	defaultReadTimeout     = 30 * time.Second
	defaultWriteTimeout    = 30 * time.Second
	defaultAddr            = ":20202"
	defaultShutdownTimeout = 30 * time.Second
)

type Server struct {
	log             *slog.Logger
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func New(log *slog.Logger, handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         defaultAddr,
	}
	s := &Server{
		log:             log,
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Start() {
	const op = "httpserver.Start"
	log := s.log.With(slog.String("op", op))

	l, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		log.Error("failed to start server", logger.Err(err))
		s.notify <- err
		close(s.notify)
		return
	}

	log.Info("http server running", slog.String("addr", l.Addr().String()))

	s.notify <- s.server.Serve(l)
	close(s.notify)
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	const op = "httpserver.Shutdown"
	log := s.log.With(slog.String("op", op))

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	log.Info("http server stopping")
	return s.server.Shutdown(ctx)
}
