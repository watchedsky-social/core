package utils

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

var DefaultConfig = fiber.Config{
	Prefork:           false,
	StrictRouting:     false,
	CaseSensitive:     true,
	UnescapePath:      true,
	EnablePrintRoutes: false,
	GETOnly:           true,
	BodyLimit:         -1,
	ServerHeader:      "watchedsky",
	AppName:           "WatchedSky",
	Network:           fiber.NetworkTCP,
}

func DefaultFiberMiddlewares(healthCheck healthcheck.HealthChecker) []any {
	return []any{
		recover.New(),
		compress.New(),
		cors.New(),
		requestid.New(),
		logger.New(logger.Config{
			CustomTags: map[string]logger.LogFunc{
				"http_content_length": LogHeaderFunc("content-length", true),
				"http_referer":        LogHeaderFunc("referer", false),
				"http_user_agent":     LogHeaderFunc("user-agent", false),
			},
			Format:     "${ip} - - [${time}] \"${method} ${path}\" ${status} ${http_content_length} \"${http_referer}\" \"${http_user_agent}\"\n",
			TimeFormat: time.RFC3339,
			Next: func(c *fiber.Ctx) bool {
				return c.Path() == "/livez" || c.Path() == "/readyz"
			},
		}),
		helmet.New(),
		healthcheck.New(healthcheck.Config{
			LivenessProbe:  healthCheck,
			ReadinessProbe: healthCheck,
		}),
	}
}

func LogHeaderFunc(headerName string, responseHeader bool) logger.LogFunc {
	return func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
		h := c.GetReqHeaders()
		if responseHeader {
			h = c.GetRespHeaders()
		}

		if h != nil {
			headers := http.Header(h)
			val := headers.Get(headerName)
			if val == "" {
				val = "-"
			}

			return output.WriteString(val)
		}

		return output.WriteString("-")
	}
}
