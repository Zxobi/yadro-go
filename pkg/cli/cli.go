package cli

import (
	"flag"
)

const (
	DefaultPort = 20202

	flagC    = "c"
	flagPort = "port"
)

type Options struct {
	C    string
	Port int
}

func ReadCliOptions() (opt Options) {
	flag.IntVar(&opt.Port, flagPort, DefaultPort, "port for webserver")
	flag.StringVar(&opt.C, flagC, ".", "path to search for config")
	flag.Parse()
	return
}
