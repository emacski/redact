package main

import (
	"log"
	"runtime"
)

var version = "dev"

func init() {
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()
}

func main() {
	log.SetFlags(0)
	rootCmd.Execute()
}
