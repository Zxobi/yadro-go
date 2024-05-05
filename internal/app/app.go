package app

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log/slog"
	nethttp "net/http"
	"os/signal"
	"strconv"
	"syscall"
	"yadro-go/internal/adapter/primary/http"
	"yadro-go/internal/adapter/secondary/repository"
	"yadro-go/internal/adapter/secondary/xkcd"
	"yadro-go/internal/core/service"
	"yadro-go/internal/core/service/stemming"
	"yadro-go/pkg/config"
	"yadro-go/pkg/httpserver"
	logutil "yadro-go/pkg/logger"
	"yadro-go/pkg/sqlite"
)

func Run(logger *slog.Logger, cfg *config.Config) error {
	const op = "app.Run"
	log := logger.With(slog.String("op", op))

	sqliteDb := sqlite.SQLite{}
	db, err := sqliteDb.Connect(cfg.Dsn)
	if err != nil {
		log.Error("failed to connect to sqlite", logutil.Err(err))
		return err
	}
	defer db.Close()

	log.Info("migrations running")

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Error("failed to create SQLite driver", logutil.Err(err))
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+cfg.Migrations, "sqlite3", driver)
	if err != nil {
		log.Error("failed to create migration instance", logutil.Err(err))
		return err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Error("failed to apply migrations", logutil.Err(err))
		return err
	}

	log.Info("migrations done")

	comicsRepo := repository.NewComicRepository(log, db)
	keywordsRepo := repository.NewKeywordRepository(log, db)
	stemmer := stemming.New()

	client := xkcd.NewHttpClient(cfg.Url, cfg.ReqTimeout)
	updater := service.NewUpdater(logger, stemmer, comicsRepo, keywordsRepo, client, cfg.FetchLimit, cfg.Parallel)
	scanner := service.NewScanner(logger, stemmer, comicsRepo, keywordsRepo)

	handler := nethttp.NewServeMux()
	http.NewRouter(logger, handler, scanner, updater, http.ScanTimeout(cfg.ScanTimeout), http.ScanLimit(cfg.ScanLimit))
	server := httpserver.New(logger, handler, httpserver.Port(strconv.Itoa(cfg.Port)))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go updater.StartScheduler(ctx, cfg.SchedulerHour, cfg.SchedulerMinute)
	go server.Start()

	select {
	case <-ctx.Done():
		log.Info("stopping: stop signal received")
	case err = <-server.Notify():
		log.Error("stopping: httpServer notify", logutil.Err(err))
	}

	if err = server.Shutdown(); err != nil {
		log.Error("server shutdown error", logutil.Err(err))
		return err
	}

	return err
}
