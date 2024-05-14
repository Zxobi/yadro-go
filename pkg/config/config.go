package config

import (
	"github.com/spf13/viper"
	"math"
	"time"
)

const (
	optDsn             = "dsn"
	optSourceUrl       = "source_url"
	optMigrations      = "migrations"
	optReqTimeout      = "req_timeout_sec"
	optFetchLimit      = "fetch_limit"
	optParallel        = "parallel"
	optScanTimeout     = "scan_timeout"
	optScanLimit       = "scan_limit"
	optPort            = "port"
	optSchedulerHour   = "scheduler_hour"
	optSchedulerMinute = "scheduler_minute"
	optTokenSecret     = "token_secret"
	optTokenTTL        = "token_ttl"
)

type Config struct {
	Dsn             string
	Url             string
	Migrations      string
	TokenSecret     string
	FetchLimit      int
	Parallel        int
	ScanLimit       int
	Port            int
	SchedulerHour   int
	SchedulerMinute int
	ReqTimeout      time.Duration
	ScanTimeout     time.Duration
	TokenTTL        time.Duration
}

func ReadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.SetDefault(optSourceUrl, "https://xkcd.com")
	viper.SetDefault(optSchedulerHour, 3)
	viper.SetDefault(optSchedulerMinute, 0)
	viper.SetDefault(optReqTimeout, math.MaxInt)
	viper.SetDefault(optFetchLimit, math.MaxInt)
	viper.SetDefault(optParallel, 1)
	viper.SetDefault(optScanTimeout, math.MaxInt)
	viper.SetDefault(optScanLimit, 10)
	viper.SetDefault(optPort, 20202)
	viper.SetDefault(optDsn, "database.db")
	viper.SetDefault(optMigrations, "migrations")
	viper.SetDefault(optTokenSecret, "token-secret")
	viper.SetDefault(optTokenTTL, 1*time.Hour)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		Dsn:             viper.GetString(optDsn),
		Url:             viper.GetString(optSourceUrl),
		Migrations:      viper.GetString(optMigrations),
		TokenSecret:     viper.GetString(optTokenSecret),
		FetchLimit:      viper.GetInt(optFetchLimit),
		ScanLimit:       viper.GetInt(optScanLimit),
		Parallel:        viper.GetInt(optParallel),
		Port:            viper.GetInt(optPort),
		SchedulerHour:   viper.GetInt(optSchedulerHour),
		SchedulerMinute: viper.GetInt(optSchedulerMinute),
		ReqTimeout:      viper.GetDuration(optReqTimeout),
		ScanTimeout:     viper.GetDuration(optScanTimeout),
		TokenTTL:        viper.GetDuration(optTokenTTL),
	}, nil
}
