package main

import (
	"runtime"
	"log"
	_ "./loginserver"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Print("Server stopped.")
}
