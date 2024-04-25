package app

import (
	"context"
	"log/slog"
	nethttp "net/http"
	"os/signal"
	"strconv"
	"syscall"
	"yadro-go/internal/controller/http"
	"yadro-go/internal/database"
	"yadro-go/internal/service"
	"yadro-go/internal/xkcd"
	"yadro-go/pkg/config"
	"yadro-go/pkg/httpserver"
	logutil "yadro-go/pkg/logger"
)

func Run(logger *slog.Logger, cfg *config.Config) error {
	const op = "app.Run"
	log := logger.With(slog.String(op, op))

	db, err := database.NewFileDatabase(logger, cfg.DbFile, cfg.IndexFile)
	if err != nil {
		log.Error("failed to create database", slog.Any("err", err))
		return err
	}

	client := xkcd.NewHttpClient(cfg.Url, cfg.ReqTimeout)
	updater := service.NewUpdater(logger, client, db, cfg.FetchLimit, cfg.Parallel)
	scanner := service.NewScanner(logger, db, db)

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

	err = server.Shutdown()
	if err != nil {
		log.Error("server shutdown error", logutil.Err(err))
		return err
	}

	return nil
}
