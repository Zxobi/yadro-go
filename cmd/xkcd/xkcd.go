package main

import (
	"encoding/json"
	"fmt"
	"os"
	"yadro-go/pkg/cli"
	"yadro-go/pkg/config"
	"yadro-go/pkg/database"
	"yadro-go/pkg/service"
	"yadro-go/pkg/xkcd"
)

func main() {
	cliOpt := cli.ReadCliOptions()
	cfg, err := config.ReadConfig(cliOpt.C)
	if err != nil {
		exitWithErr(err)
	}

	client := xkcd.NewHttpClient(cfg.Url, cfg.ReqTimeout)
	db, err := database.NewFileDatabase(cfg.DbFile)
	if err != nil {
		exitWithErr(err)
	}
	srv := service.NewComicsService(client, db, cfg.FetchLimit)

	if err = srv.Fetch(); err != nil {
		exitWithErr(err)
		return
	}

	if cliOpt.O && cliOpt.N > 0 {
		if err = printDb(db, cliOpt.N); err != nil {
			exitWithErr(err)
		}
	}
}

func exitWithErr(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "Error occured: %v\n", err)
	os.Exit(1)
}

func printDb(db service.IDatabase, n int) error {
	records := db.Read()

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
