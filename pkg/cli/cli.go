package cli

import (
	"flag"
	"math"
)

const oFlag = "o"
const nFlag = "n"
const cFlag = "c"

type Options struct {
	O bool
	N int
	C string
}

func ReadCliOptions() (opt Options) {
	flag.BoolVar(&opt.O, oFlag, false, "print to stdout")
	flag.IntVar(&opt.N, nFlag, math.MaxInt, "shorten output to [n] records")
	flag.StringVar(&opt.C, cFlag, ".", "path to search for config")
	flag.Parse()
	return
}
