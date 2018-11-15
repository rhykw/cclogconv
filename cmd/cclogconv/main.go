package main

import (
	"github.com/rhykw/cclogconv"
	"os"
)

func main() {
	cl := cclogconv.CCLogConv{Out: os.Stdout, Err: os.Stderr}
	exitCode := cl.Run(os.Args[1:])
	os.Exit(exitCode)
}
