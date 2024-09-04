package feed

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/watchedsky-social/core/internal/config"
	"github.com/watchedsky-social/core/internal/database/query"
	"github.com/watchedsky-social/core/internal/utils"
)

func ServeCustomAlgorithmFeed(ctx context.Context) error {
	server := fiber.New(utils.DefaultConfig)

	server.Use(
		utils.DefaultFiberMiddlewares(func(c *fiber.Ctx) bool { return query.Q.Available() })...,
	)

	server.Get("/.well-known/did.json", wellKnown)
	server.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	server.Get("/xrpc/app.bsky.feed.getFeedSkeleton", handleFeedAlgorithm)

	go func() {
		<-ctx.Done()
		server.Shutdown()
	}()

	return server.Listen(fmt.Sprintf(":%d", config.Config.Feed.Port))
}
