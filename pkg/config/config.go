package config

import (
	"github.com/spf13/viper"
	"math"
	"time"
)

const (
	optDbFile        = "db_file"
	optSourceUrl     = "source_url"
	optReqTimeout    = "req_timeout_sec"
	optFetchLimit    = "fetch_limit"
	optSaveBatchSize = "save_batch_size"
)

type Config struct {
	DbFile        string
	Url           string
	ReqTimeout    time.Duration
	FetchLimit    int
	SaveBatchSize int
}

func ReadConfig(path string) (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.SetDefault(optDbFile, "database.json")
	viper.SetDefault(optSourceUrl, "https://xkcd.com")
	viper.SetDefault(optReqTimeout, math.MaxInt)
	viper.SetDefault(optFetchLimit, math.MaxInt)
	viper.SetDefault(optSaveBatchSize, math.MaxInt)

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	return Config{
		DbFile:        viper.GetString(optDbFile),
		Url:           viper.GetString(optSourceUrl),
		ReqTimeout:    viper.GetDuration(optReqTimeout) * time.Second,
		FetchLimit:    viper.GetInt(optFetchLimit),
		SaveBatchSize: viper.GetInt(optSaveBatchSize),
	}, nil
}
