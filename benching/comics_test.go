package benching

import (
	"context"
	"testing"
	"time"
	"yadro-go/pkg/database"
	"yadro-go/pkg/service"
	"yadro-go/pkg/xkcd"
)

func BenchmarkFetchParallel20(b *testing.B) {
	fetch(20)
}

func BenchmarkFetchParallel50(b *testing.B) {
	fetch(50)
}

func BenchmarkFetchParallel100(b *testing.B) {
	fetch(100)
}

func BenchmarkFetchParallel200(b *testing.B) {
	fetch(200)
}

func BenchmarkFetchParallel400(b *testing.B) {
	fetch(400)
}

func BenchmarkFetchParallel800(b *testing.B) {
	fetch(800)
}

type dbStub struct {
}

func (d *dbStub) Read() database.RecordMap {
	return make(database.RecordMap)
}

func (d *dbStub) Write(_ database.RecordMap) error {
	return nil
}

func fetch(parallel int) {
	client := xkcd.NewHttpClient("https://xkcd.com", time.Minute)
	srv := service.NewComicsService(client, &dbStub{}, 2000, parallel)

	if err := srv.Fetch(context.Background()); err != nil {
		panic(err)
	}
}
