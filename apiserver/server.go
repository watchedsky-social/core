package apiserver

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/watchedsky-social/core/internal/database/query"
	"github.com/watchedsky-social/core/internal/handlers"
	"github.com/watchedsky-social/core/internal/metrics"
	"github.com/watchedsky-social/core/internal/utils"
)

var codeMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "apiserver_http_status_code",
	Help: "Total counts of each path by status code, with query string removed",
},
	[]string{"path", "code"},
)

var spanMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "apiserver_total_response_time_us",
	Help: "Total response time of each path by status code, with query string removed, in microseceonds",
},
	[]string{"path", "code"},
)

func registerMetrics() {

	metrics.Registry.MustRegister(
		codeMetric,
		spanMetric,
	)
}

func metricMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()

	elapsed := time.Since(start)
	path := string(c.Request().URI().Path())
	code := fmt.Sprintf("%d", c.Response().StatusCode())

	codeMetric.With(prometheus.Labels{"path": path, "code": code}).Add(1)
	spanMetric.With(prometheus.Labels{"path": path, "code": code}).Add(float64(elapsed / time.Microsecond))

	return err
}

func Run(ctx context.Context) error {
	registerMetrics()
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

	app.Use(metricMiddleware)

	api := app.Group("/api/v1")
	api.Get("/info", handlers.Info)
	api.Get("/typeahead", handlers.Typeahead)
	api.Get("/zones/visible", handlers.VisibleZones)
	api.Get("/zones/watchid", handlers.GetWatchID)
	api.Get("/alerts/recent", handlers.RecentAlerts)
	api.Get("/alerts/:alertid", handlers.AlertByID)

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{
		Registry: metrics.Registry,
	})))

	go func() {
		app.Listen(fmt.Sprintf(":%d", port))
	}()

	<-ctx.Done()
	return app.Shutdown()
}
