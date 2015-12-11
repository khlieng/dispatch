package main

import (
	"runtime"

	"github.com/khlieng/dispatch/commands"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	commands.Execute()
}
