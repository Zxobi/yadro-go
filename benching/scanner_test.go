package benching

import (
	"context"
	"log/slog"
	"testing"
	"time"
	"yadro-go/benching/logger"
	"yadro-go/internal/database"
	"yadro-go/internal/service"
	"yadro-go/internal/xkcd"
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
	fileDb, err := database.NewFileDatabase(log, "database.json", "index.json")
	if err != nil {
		panic(err)
	}

	client := xkcd.NewHttpClient("https://xkcd.com", time.Minute)
	srv := service.NewUpdater(log, client, fileDb, 99999, 200)

	if _, err = srv.Update(context.Background()); err != nil {
		panic(err)
	}

	scanner = service.NewScanner(log, fileDb, fileDb)
}

func BenchmarkScanNoIndex(b *testing.B) {
	b.Run("query_small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			scanner.Scan(context.Background(), querySmall, false)
		}
	})
	b.Run("query_medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			scanner.Scan(context.Background(), queryMedium, false)
		}
	})
	b.Run("query_large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			scanner.Scan(context.Background(), queryLarge, false)
		}
	})
}

func BenchmarkScanIndex(b *testing.B) {
	b.Run("query_small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			scanner.Scan(context.Background(), querySmall, true)
		}
	})
	b.Run("query_medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			scanner.Scan(context.Background(), queryMedium, true)
		}
	})
	b.Run("query_large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			scanner.Scan(context.Background(), queryLarge, true)
		}
	})
}
