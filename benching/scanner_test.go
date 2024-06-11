package benching

import (
	"context"
	"log/slog"
	"testing"
	"time"
	"yadro-go/internal/adapter/secondary/repository"
	"yadro-go/internal/adapter/secondary/xkcd"
	"yadro-go/internal/core/service"
	"yadro-go/internal/core/service/stemming"
	logutil "yadro-go/pkg/logger"
	"yadro-go/pkg/sqlite"
	"yadro-go/test/logger"
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

	sqliteDb := sqlite.SQLite{}
	db, err := sqliteDb.Connect("test.db")
	if err != nil {
		log.Error("failed to connect to sqlite", logutil.Err(err))
		panic(err)
	}

	comicsRepo := repository.NewComicRepository(log, db)
	keywordsRepo := repository.NewKeywordRepository(log, db)
	stemmer := stemming.New()

	client := xkcd.NewHttpClient(log, "https://xkcd.com", time.Minute)
	updater := service.NewUpdater(log, stemmer, comicsRepo, keywordsRepo, client, 2000, 200)

	if _, err = updater.Update(context.Background()); err != nil {
		panic(err)
	}

	scanner = service.NewScanner(log, stemmer, comicsRepo, keywordsRepo)
}

func BenchmarkScanNoIndex(b *testing.B) {
	b.Run("query_small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := scanner.Scan(context.Background(), querySmall, false); err != nil {
				b.Error(err)
			}
		}
	})
	b.Run("query_medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := scanner.Scan(context.Background(), queryMedium, false); err != nil {
				b.Error(err)
			}
		}
	})
	b.Run("query_large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := scanner.Scan(context.Background(), queryLarge, false); err != nil {
				b.Error(err)
			}
		}
	})
}

func BenchmarkScanIndex(b *testing.B) {
	b.Run("query_small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := scanner.Scan(context.Background(), querySmall, true); err != nil {
				b.Error(err)
			}
		}
	})
	b.Run("query_medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := scanner.Scan(context.Background(), queryMedium, true); err != nil {
				b.Error(err)
			}
		}
	})
	b.Run("query_large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := scanner.Scan(context.Background(), queryLarge, true); err != nil {
				b.Error(err)
			}
		}
	})
}
