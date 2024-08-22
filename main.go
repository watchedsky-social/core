package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/watchedsky-social/core/operations"
)

var Version = "dev"

func main() {
	files := []string{".env"}
	if envFile, ok := os.LookupEnv("WATCHEDSKY_ENV_FILE"); ok {
		files = append(files, envFile)
	}
	godotenv.Load(files...)

	operations.Main()
}
