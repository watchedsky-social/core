package config_test

import (
	"os"
	"testing"

	"github.com/watchedsky-social/core/internal/config"
)

func TestEnvConfig(t *testing.T) {
	environ := map[string]string{
		"WATCHEDSKY_DB_PASSWORD":                 "FAKE",
		"WATCHEDSKY_ALERTS_BLUESKY_ID":           "skeetbot.watchedsky.social",
		"WATCHEDSKY_ALERTS_BLUESKY_APP_PASSWORD": "FAKE",
	}

	for k, v := range environ {
		os.Setenv(k, v)
	}

	config.Load()

}
