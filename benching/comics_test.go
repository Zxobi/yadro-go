package benching

import (
	"context"
	"log/slog"
	"testing"
	"time"
	"yadro-go/pkg/database"
	"yadro-go/pkg/service"
	"yadro-go/pkg/xkcd"
)

func BenchmarkFetchParallel20(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fetch(20)
	}
}

func BenchmarkFetchParallel50(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fetch(50)
	}
}

func BenchmarkFetchParallel100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fetch(100)
	}
}

func BenchmarkFetchParallel200(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fetch(200)
	}
}

func BenchmarkFetchParallel400(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fetch(400)
	}
}

func BenchmarkFetchParallel800(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fetch(800)
	}
}

type dbStub struct {
}

func (d *dbStub) Records() database.RecordMap {
	return make(database.RecordMap)
}

func (d *dbStub) Save(_ database.RecordMap) error {
	return nil
}

func fetch(parallel int) {
	client := xkcd.NewHttpClient("https://xkcd.com", time.Minute)
	srv := service.NewComicsService(slog.Default(), client, &dbStub{}, 99999, parallel)

	if err := srv.Fetch(context.Background()); err != nil {
		panic(err)
	}
}
