package config

import (
	"github.com/spf13/viper"
	"math"
	"time"
)

const (
	optDbFile          = "db_file"
	optIndexFile       = "index_file"
	optSourceUrl       = "source_url"
	optReqTimeout      = "req_timeout_sec"
	optFetchLimit      = "fetch_limit"
	optParallel        = "parallel"
	optScanTimeout     = "scan_timeout"
	optScanLimit       = "scan_limit"
	optPort            = "port"
	optSchedulerHour   = "scheduler_hour"
	optSchedulerMinute = "scheduler_minute"
)

type Config struct {
	DbFile          string
	IndexFile       string
	Url             string
	UpdateCron      string
	FetchLimit      int
	Parallel        int
	ScanLimit       int
	Port            int
	SchedulerHour   int
	SchedulerMinute int
	ReqTimeout      time.Duration
	ScanTimeout     time.Duration
}

func ReadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.SetDefault(optDbFile, "database.json")
	viper.SetDefault(optIndexFile, "index.json")
	viper.SetDefault(optSourceUrl, "https://xkcd.com")
	viper.SetDefault(optSchedulerHour, 3)
	viper.SetDefault(optSchedulerMinute, 0)
	viper.SetDefault(optReqTimeout, math.MaxInt)
	viper.SetDefault(optFetchLimit, math.MaxInt)
	viper.SetDefault(optParallel, 1)
	viper.SetDefault(optScanTimeout, math.MaxInt)
	viper.SetDefault(optScanLimit, 10)
	viper.SetDefault(optPort, 20202)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		DbFile:          viper.GetString(optDbFile),
		IndexFile:       viper.GetString(optIndexFile),
		Url:             viper.GetString(optSourceUrl),
		FetchLimit:      viper.GetInt(optFetchLimit),
		ScanLimit:       viper.GetInt(optScanLimit),
		Parallel:        viper.GetInt(optParallel),
		Port:            viper.GetInt(optPort),
		SchedulerHour:   viper.GetInt(optSchedulerHour),
		SchedulerMinute: viper.GetInt(optSchedulerMinute),
		ReqTimeout:      viper.GetDuration(optReqTimeout),
		ScanTimeout:     viper.GetDuration(optScanTimeout),
	}, nil
}
