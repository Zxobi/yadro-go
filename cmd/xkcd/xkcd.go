package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"yadro-go/pkg/cli"
	"yadro-go/pkg/config"
	"yadro-go/pkg/database"
	"yadro-go/pkg/service"
	"yadro-go/pkg/xkcd"
)

func main() {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	cliOpt := cli.ReadCliOptions()
	cfg, err := config.ReadConfig(cliOpt.C)
	if err != nil {
		exitWithErr(err)
	}

	log.Debug("config loaded", slog.Any("config", cfg))

	client := xkcd.NewHttpClient(cfg.Url, cfg.ReqTimeout)
	db, err := database.NewFileDatabase(log, cfg.DbFile, cfg.IndexFile)
	if err != nil {
		log.Error("failed to create database", slog.Any("err", err))
		exitWithErr(err)
	}
	srv := service.NewComicsService(log, client, db, cfg.FetchLimit, cfg.Parallel)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	if err = srv.Fetch(ctx); err != nil {
		log.Error("failed to fetch records", slog.Any("err", err))
		exitWithErr(err)
		return
	}

	if cliOpt.O && cliOpt.N > 0 {
		if err = printDb(db, cliOpt.N); err != nil {
			exitWithErr(err)
		}
	}

	if cliOpt.S != "" {
		scanner := service.NewScanner(log, db, db)
		scanCtx, scanCancel := context.WithTimeout(ctx, cfg.ScanTimeout)
		defer scanCancel()

		matches := scanner.Scan(scanCtx, cliOpt.S, cliOpt.I)
		if len(matches) == 0 {
			return
		} else if len(matches) > cfg.ScanLimit {
			matches = matches[:cfg.ScanLimit]
		}

		fmt.Println(strings.Join(matches, "\n"))
	}
}

func exitWithErr(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "Error occured: %v\n", err)
	os.Exit(1)
}

func printDb(db service.RecordRepository, n int) error {
	records := db.Records()

	if n >= len(records) {
		return printRecords(records)
	}

	printMap := make(database.RecordMap, n)
	count := 0
	for k, v := range records {
		if count == n {
			break
		}
		printMap[k] = v
		count++
	}

	return printRecords(printMap)
}

func printRecords(records database.RecordMap) error {
	data, err := json.MarshalIndent(records, "", "\t")
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}
