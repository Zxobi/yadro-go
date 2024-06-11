package service

import (
	"context"
	"yadro-go/internal/core/domain"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type Stemmer interface {
	StemString(str string) []string
	StemComic(comic *domain.Comic) []string
}

type ComicProvider interface {
	GetById(id int) (*domain.Comic, error)
}

type ComicRepository interface {
	Comics(ctx context.Context, nums []int) ([]*domain.Comic, error)
	All(ctx context.Context) ([]*domain.Comic, error)
	Save(ctx context.Context, comics []*domain.Comic) error
}

type KeywordRepository interface {
	Keywords(ctx context.Context, keywords []string) ([]*domain.ComicKeyword, error)
	Save(ctx context.Context, keywords []*domain.ComicKeyword) error
}

type UserRepository interface {
	UserByUsername(ctx context.Context, username string) (*domain.User, error)
}

type TokenManager interface {
	Token(username string) (string, error)
	Verify(token string) (string, error)
}
