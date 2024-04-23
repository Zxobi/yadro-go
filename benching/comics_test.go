package benching

import (
	"context"
	"log/slog"
	"testing"
	"time"
	"yadro-go/benching/logger"
	"yadro-go/pkg/database"
	"yadro-go/pkg/service"
	"yadro-go/pkg/xkcd"
)

func BenchmarkFetchParallel(b *testing.B) {
	b.Run("parallel_20", func(b *testing.B) {
		srv := newService(20)
		for i := 0; i < b.N; i++ {
			if err := srv.Fetch(context.Background()); err != nil {
				b.Error(err)
			}
		}
	})
	b.Run("parallel_50", func(b *testing.B) {
		srv := newService(50)
		for i := 0; i < b.N; i++ {
			if err := srv.Fetch(context.Background()); err != nil {
				b.Error(err)
			}
		}
	})
	b.Run("parallel_100", func(b *testing.B) {
		srv := newService(100)
		for i := 0; i < b.N; i++ {
			if err := srv.Fetch(context.Background()); err != nil {
				b.Error(err)
			}
		}
	})
	b.Run("parallel_200", func(b *testing.B) {
		srv := newService(200)
		for i := 0; i < b.N; i++ {
			if err := srv.Fetch(context.Background()); err != nil {
				b.Error(err)
			}
		}
	})
	b.Run("parallel_400", func(b *testing.B) {
		srv := newService(400)
		for i := 0; i < b.N; i++ {
			if err := srv.Fetch(context.Background()); err != nil {
				b.Error(err)
			}
		}
	})
	b.Run("parallel_800", func(b *testing.B) {
		srv := newService(400)
		for i := 0; i < b.N; i++ {
			if err := srv.Fetch(context.Background()); err != nil {
				b.Error(err)
			}
		}
	})
}

type dbStub struct {
}

func (d *dbStub) Records() database.RecordMap {
	return make(database.RecordMap)
}

func (d *dbStub) Save(_ database.RecordMap) error {
	return nil
}

func newService(parallel int) *service.ComicsService {
	log := slog.New(logger.EmptyHandler{})
	client := xkcd.NewHttpClient("https://xkcd.com", time.Minute)

	return service.NewComicsService(log, client, &dbStub{}, 99999, parallel)
}
