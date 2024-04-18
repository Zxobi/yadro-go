package cli

import (
	"flag"
	"math"
)

const oFlag = "o"
const nFlag = "n"
const cFlag = "c"
const sFlag = "s"
const iFlag = "i"

type Options struct {
	O bool
	N int
	C string
	S string
	I bool
}

func ReadCliOptions() (opt Options) {
	flag.BoolVar(&opt.O, oFlag, false, "print to stdout")
	flag.BoolVar(&opt.I, iFlag, false, "use index")
	flag.IntVar(&opt.N, nFlag, math.MaxInt, "shorten output to [n] records")
	flag.StringVar(&opt.C, cFlag, ".", "path to search for config")
	flag.StringVar(&opt.S, sFlag, "", "query string")
	flag.Parse()
	return
}
