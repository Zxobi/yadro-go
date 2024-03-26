package main

import (
	"errors"
	"flag"
)

const strFlag = "s"

func parseStr() (string, error) {
	str := flag.String(strFlag, "", "string to stem")

	flag.Parse()

	if err := validateFlag(); err != nil {
		return "", err
	}
	return *str, nil
}

func validateFlag() error {
	strSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == strFlag {
			strSet = true
		}
	})

	if !strSet {
		return errors.New("[-s] flag is required")
	}
	return nil
}
