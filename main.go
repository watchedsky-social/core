package main

import (
	"log"

	"github.com/watchedsky-social/core/operations"
)

var Version = "dev"

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	operations.Main()
}
