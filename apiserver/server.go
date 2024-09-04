package apiserver

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/watchedsky-social/core/internal/database/query"
	"github.com/watchedsky-social/core/internal/handlers"
	"github.com/watchedsky-social/core/internal/utils"
)

func Run(ctx context.Context) error {
	port := uint16(8888)
	if ps, ok := os.LookupEnv("PORT"); ok {
		if tmp, err := strconv.ParseUint(ps, 10, 16); err == nil {
			if tmp > 1024 {
				port = uint16(tmp)
			}
		} else {
			log.Printf("port %q is not a valid uint16 between 1025 - 65535, defaulting to 8888", ps)
		}
	}

	app := fiber.New(utils.DefaultConfig)

	app.Use(
		utils.DefaultFiberMiddlewares(func(c *fiber.Ctx) bool { return query.Q.Available() })...,
	)

	api := app.Group("/api/v1")
	api.Get("/info", handlers.Info)
	api.Get("/typeahead", handlers.Typeahead)
	api.Get("/zones/visible", handlers.VisibleZones)
	api.Get("/zones/watchid", handlers.GetWatchID)
	api.Get("/alerts/recent", handlers.RecentAlerts)
	api.Get("/alerts/:alertid", handlers.AlertByID)

	go func() {
		app.Listen(fmt.Sprintf(":%d", port))
	}()

	<-ctx.Done()
	return app.Shutdown()
}
