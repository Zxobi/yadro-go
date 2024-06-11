package logger

import (
	"context"
	"log/slog"
)

type EmptyHandler struct {
}

func (e EmptyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return false
}

func (e EmptyHandler) Handle(ctx context.Context, record slog.Record) error {
	return nil
}

func (e EmptyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return e
}

func (e EmptyHandler) WithGroup(name string) slog.Handler {
	return e
}
