package primary

import (
	"context"
	"yadro-go/internal/core/domain"
)

type QueryScanner interface {
	Scan(ctx context.Context, query string, useIndex bool) ([]string, error)
}

type Updater interface {
	Update(ctx context.Context) (int, error)
}

type Auth interface {
	Login(ctx context.Context, username string, password string) (string, error)
	Authenticate(ctx context.Context, token string) (*domain.User, error)
}
