package args

import (
	"flag"
)

var (
	Development bool
)

func init() {
	flag.BoolVar(&Development, "dev", false, "")
	flag.Parse()
}
