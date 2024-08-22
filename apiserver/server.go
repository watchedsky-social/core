package apiserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/watchedsky-social/core/internal/database/query"
	"github.com/watchedsky-social/core/internal/handlers"
)

func headerFunc(headerName string, responseHeader bool) logger.LogFunc {
	return func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
		h := c.GetReqHeaders()
		if responseHeader {
			h = c.GetRespHeaders()
		}

		if h != nil {
			headers := http.Header(h)
			val := headers.Get(headerName)

			return output.WriteString(val)
		}

		return output.WriteString("-")
	}
}

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

	app := fiber.New(fiber.Config{
		Prefork:           false,
		StrictRouting:     false,
		CaseSensitive:     true,
		UnescapePath:      true,
		EnablePrintRoutes: true,
		GETOnly:           true,
		BodyLimit:         -1,
		ServerHeader:      "watchedsky",
		AppName:           "WatchedSky",
		Network:           fiber.NetworkTCP,
	})

	// '$remote_addr - $remote_user [$time_local] '    '"$request" $status $body_bytes_sent '    '"$http_referer" "$http_user_agent"';
	app.Use(
		recover.New(),
		compress.New(),
		cors.New(),
		requestid.New(),
		logger.New(logger.Config{
			CustomTags: map[string]logger.LogFunc{
				"http_content_length": headerFunc("content-length", true),
				"http_referer":        headerFunc("referer", false),
				"http_user_agent":     headerFunc("user-agent", false),
			},
			Format:     `${ip} - - [${time}] "${method} ${path}" ${status} ${http_content_length} "${http_referer}" "${http_user_agent}"`,
			TimeFormat: time.RFC3339,
		}),
		helmet.New(),
		healthcheck.New(healthcheck.Config{
			LivenessProbe: func(c *fiber.Ctx) bool {
				return query.Q.Available()
			},
			ReadinessProbe: func(c *fiber.Ctx) bool {
				r, e := query.Mapsearch.WithContext(c.UserContext()).PrefixSearch("44813")
				return e != nil && len(r) > 0
			},
		}),
	)

	api := app.Group("/api/v1")
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
