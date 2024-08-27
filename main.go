package main

import (
	"log"

	"github.com/watchedsky-social/core/operations"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	operations.Main()
}
