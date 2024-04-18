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

var (
	scanner     *service.Scanner
	querySmall  = "I'm following your questions"
	queryMedium = "The dedicated follower carried a bottle of water to quench his thirst during the long hike"
	queryLarge  = "The quick brown fox jumps over the lazy dog. " +
		"Today is a beautiful day with clear skies and gentle breezes. " +
		"I plan to take a leisurely walk in the park later. " +
		"Meanwhile, I'll grab a cup of coffee and catch up on some reading. " +
		"Life is good when you take time to appreciate the little things."
)

func init() {
	log := slog.New(logger.EmptyHandler{})
	fileDb, err := database.NewFileDatabase(log, "database.json", "database_index.json")
	if err != nil {
		panic(err)
	}

	client := xkcd.NewHttpClient("https://xkcd.com", time.Minute)
	srv := service.NewComicsService(log, client, fileDb, 99999, 200)

	if err = srv.Fetch(context.Background()); err != nil {
		panic(err)
	}

	scanner = service.NewScanner(log, fileDb, fileDb)
}

func BenchmarkScanNoIndexQuerySmall(b *testing.B) {
	scanner.Scan(context.Background(), querySmall, false)
}

func BenchmarkScanNoIndexQueryMedium(b *testing.B) {
	scanner.Scan(context.Background(), queryMedium, false)
}

func BenchmarkScanNoIndexQueryLarge(b *testing.B) {
	scanner.Scan(context.Background(), queryLarge, false)
}

func BenchmarkScanIndexQuerySmall(b *testing.B) {
	scanner.Scan(context.Background(), querySmall, true)
}

func BenchmarkScanIndexQueryMedium(b *testing.B) {
	scanner.Scan(context.Background(), queryMedium, true)
}

func BenchmarkScanIndexQueryLarge(b *testing.B) {
	scanner.Scan(context.Background(), queryLarge, true)
}
