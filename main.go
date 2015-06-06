package main

import (
	"runtime"

	"github.com/khlieng/name_pending/commands"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	commands.Execute()
}
