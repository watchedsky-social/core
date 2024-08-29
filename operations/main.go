package operations

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/watchedsky-social/core/alerts"
	"github.com/watchedsky-social/core/apiserver"
	"github.com/watchedsky-social/core/internal/config"
	"github.com/watchedsky-social/core/internal/database/query"
	"github.com/watchedsky-social/core/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type command func(context.Context) error

var commandMap map[string]command = map[string]command{
	"api-server":      apiserver.Run,
	"alert-collector": alerts.CollectNewAlerts,
	"alert-vacuum":    alerts.CleanOldAlerts,
	"hydrate-state":   alerts.InitialHydration,
	"post-alerts":     alerts.SkeetNewAlerts,
}

func install() error {
	// this has no effect under go run but i guess you can run it that way, idfc
	me, err := filepath.Abs(os.Args[0])
	if err != nil {
		return err
	}

	loc := filepath.Dir(me)
	for k := range commandMap {
		if err = os.Symlink(me, filepath.Join(loc, k)); err != nil {
			return err
		}
	}

	return nil
}

func loadEnv() {
	files := []string{".env"}
	if envFile := os.Getenv("WATCHEDSKY_ENV_FILE"); envFile != "" {
		files = append(files, envFile)
	}

	files = utils.Filter(files, func(f string) bool {
		_, err := os.Stat(f)
		return err == nil
	})

	if err := godotenv.Load(files...); err != nil {
		log.Fatal(err)
	}
}

func Main() {
	if len(os.Args) == 2 && os.Args[1] == "install" {
		if err := install(); err != nil {
			log.Fatal(err)
		}

		return
	}

	loadEnv()
	config.Load()

	dbArgs := config.Config.DB

	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require TimeZone=UTC",
		dbArgs.Host, dbArgs.User, dbArgs.Password, dbArgs.Name)), &gorm.Config{
		Logger: logger.Default,
	})
	if err != nil {
		log.Fatal(err)
	}

	query.SetDefault(db)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	me := filepath.Base(os.Args[0])
	if me == "main" && len(os.Args) > 1 {
		// we're in go run mode, take the next arg
		me = os.Args[1]
	}

	cmd, ok := commandMap[me]
	if !ok {
		log.Fatalf("%q is not a recognized mode", me)
	}

	if err = cmd(ctx); err != nil {
		log.Fatal(err)
	}
}
