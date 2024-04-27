package cli

import (
	"flag"
	"math"
)

var DefaultPort = 20202

const flagO = "o"
const flagN = "n"
const flagC = "c"
const flagS = "s"
const flagI = "i"
const flagPort = "port"

type Options struct {
	O    bool
	N    int
	C    string
	S    string
	I    bool
	Port int
}

func ReadCliOptions() (opt Options) {
	flag.BoolVar(&opt.O, flagO, false, "print to stdout")
	flag.BoolVar(&opt.I, flagI, false, "use index")
	flag.IntVar(&opt.N, flagN, math.MaxInt, "shorten output to [n] records")
	flag.IntVar(&opt.Port, flagPort, DefaultPort, "port for webserver")
	flag.StringVar(&opt.C, flagC, ".", "path to search for config")
	flag.StringVar(&opt.S, flagS, "", "query string")
	flag.Parse()
	return
}
