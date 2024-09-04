package feed

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/watchedsky-social/core/internal/config"
)

func wellKnown(c *fiber.Ctx) error {
	feedCfg := config.Config.Feed
	if !strings.HasSuffix(feedCfg.ServiceDID, feedCfg.Hostname) {
		return c.SendStatus(http.StatusNotFound)
	}

	didJSON := map[string]any{
		"@context": []string{"https://www.w3.org/ns/did/v1"},
		"id":       feedCfg.ServiceDID,
		"service": []any{
			map[string]any{
				"id":              "#bsky_fg",
				"type":            "BskyFeedGenerator",
				"serviceEndpoint": fmt.Sprintf("https://%s", feedCfg.Hostname),
			},
		},
	}

	return c.JSON(didJSON)
}
