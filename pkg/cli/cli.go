package cli

import (
	"flag"
)

const oFlag = "o"
const nFlag = "n"

type Options struct {
	O bool
	N int
}

func ReadCliOptions() Options {
	oF := flag.Bool(oFlag, false, "print to stdout")
	nF := flag.Int(nFlag, 0, "shorten output to [n] records")
	flag.Parse()

	return Options{O: *oF, N: *nF}
}
