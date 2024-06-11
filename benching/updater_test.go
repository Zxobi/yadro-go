package benching

import (
	"context"
	"log/slog"
	"testing"
	"time"
	"yadro-go/internal/adapter/secondary/xkcd"
	"yadro-go/internal/core/domain"
	"yadro-go/internal/core/service"
	"yadro-go/internal/core/service/stemming"
	"yadro-go/test/logger"
)

func BenchmarkFetchParallel(b *testing.B) {
	b.Run("parallel_20", func(b *testing.B) {
		srv := newService(20)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if _, err := srv.Update(context.Background()); err != nil {
				b.Error(err)
			}
		}
	})
	b.Run("parallel_50", func(b *testing.B) {
		srv := newService(50)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if _, err := srv.Update(context.Background()); err != nil {
				b.Error(err)
			}
		}
	})
	b.Run("parallel_100", func(b *testing.B) {
		srv := newService(100)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if _, err := srv.Update(context.Background()); err != nil {
				b.Error(err)
			}
		}
	})
	b.Run("parallel_200", func(b *testing.B) {
		srv := newService(200)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if _, err := srv.Update(context.Background()); err != nil {
				b.Error(err)
			}
		}
	})
	b.Run("parallel_400", func(b *testing.B) {
		srv := newService(400)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if _, err := srv.Update(context.Background()); err != nil {
				b.Error(err)
			}
		}
	})
	b.Run("parallel_800", func(b *testing.B) {
		srv := newService(400)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if _, err := srv.Update(context.Background()); err != nil {
				b.Error(err)
			}
		}
	})
}

type comicStub struct {
}
type keywordStub struct {
}

func (d *comicStub) Comics(_ context.Context, _ []int) ([]*domain.Comic, error) {
	return make([]*domain.Comic, 0), nil
}

func (d *comicStub) All(_ context.Context) ([]*domain.Comic, error) {
	return make([]*domain.Comic, 0), nil
}

func (d *comicStub) Save(_ context.Context, _ []*domain.Comic) error {
	return nil
}

func (d *keywordStub) Keywords(_ context.Context, _ []string) ([]*domain.ComicKeyword, error) {
	return make([]*domain.ComicKeyword, 0), nil
}

func (d *keywordStub) Save(_ context.Context, _ []*domain.ComicKeyword) error {
	return nil
}

func newService(parallel int) *service.Updater {
	log := slog.New(logger.EmptyHandler{})
	stemmer := stemming.New()
	client := xkcd.NewHttpClient(log, "https://xkcd.com", time.Minute)

	return service.NewUpdater(log, stemmer, &comicStub{}, &keywordStub{}, client, 99999, parallel)
}
