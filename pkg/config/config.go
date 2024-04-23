package config

import (
	"github.com/spf13/viper"
	"math"
	"time"
)

const (
	optDbFile      = "db_file"
	optIndexFile   = "index_file"
	optSourceUrl   = "source_url"
	optReqTimeout  = "req_timeout_sec"
	optFetchLimit  = "fetch_limit"
	optParallel    = "parallel"
	optScanTimeout = "scan_timeout"
	optScanLimit   = "scan_limit"
)

type Config struct {
	DbFile      string
	IndexFile   string
	Url         string
	FetchLimit  int
	Parallel    int
	ScanLimit   int
	ReqTimeout  time.Duration
	ScanTimeout time.Duration
}

func ReadConfig(path string) (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.SetDefault(optDbFile, "database.json")
	viper.SetDefault(optIndexFile, "index.json")
	viper.SetDefault(optSourceUrl, "https://xkcd.com")
	viper.SetDefault(optReqTimeout, math.MaxInt)
	viper.SetDefault(optFetchLimit, math.MaxInt)
	viper.SetDefault(optParallel, 1)
	viper.SetDefault(optScanTimeout, math.MaxInt)
	viper.SetDefault(optScanLimit, 10)

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	return Config{
		DbFile:      viper.GetString(optDbFile),
		IndexFile:   viper.GetString(optIndexFile),
		Url:         viper.GetString(optSourceUrl),
		FetchLimit:  viper.GetInt(optFetchLimit),
		ScanLimit:   viper.GetInt(optScanLimit),
		Parallel:    viper.GetInt(optParallel),
		ReqTimeout:  viper.GetDuration(optReqTimeout),
		ScanTimeout: viper.GetDuration(optScanTimeout),
	}, nil
}
