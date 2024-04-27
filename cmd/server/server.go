package main

import (
	"fmt"
	"log/slog"
	"os"
	"yadro-go/internal/app"
	"yadro-go/pkg/cli"
	"yadro-go/pkg/config"
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

	if cliOpt.Port != cli.DefaultPort {
		cfg.Port = cliOpt.Port
	}

	log.Debug("config loaded", slog.Any("config", cfg))

	if err = app.Run(log, cfg); err != nil {
		exitWithErr(err)
	}
}

func exitWithErr(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "Error occured: %v\n", err)
	os.Exit(1)
}
