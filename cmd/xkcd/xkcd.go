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
	cfg, err := config.ReadConfig()
	if err != nil {
		exitWithErr(err)
	}

	client := xkcd.NewHttpClient(cfg.Url, cfg.ReqTimeout)
	db := database.NewFileDatabase(cfg.DbFile)
	if err = db.Init(); err != nil {
		exitWithErr(err)
	}
	srv := service.NewComicsService(client, db, cfg.FetchLimit)

	if err = srv.Fetch(); err != nil {
		exitWithErr(err)
		return
	}

	if cliOpt.O {
		printDb(cliOpt, db)
	}
}

func exitWithErr(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "Error occured: %v\n", err)
	os.Exit(1)
}

func printDb(opt cli.Options, db database.IDatabase) {
	records, err := db.Read()
	if err != nil {
		exitWithErr(err)
	}

	if opt.N > 0 && opt.N < len(records) {
		err = printNRecords(records, opt.N)
	} else {
		err = printRecords(records)
	}

	if err != nil {
		exitWithErr(err)
	}
}

func printRecords(records database.RecordMap) error {
	data, err := json.MarshalIndent(records, "", "\t")
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}

func printNRecords(records database.RecordMap, n int) error {
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
