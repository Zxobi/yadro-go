package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	str, err := parseStr()
	if err != nil {
		exitWithErr(err)
	}

	stemmed, err := stem(str)
	if err != nil {
		exitWithErr(err)
	}
	fmt.Println(strings.Join(stemmed, " "))
}

func exitWithErr(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "Error occured: %v\n", err)
	os.Exit(1)
}
