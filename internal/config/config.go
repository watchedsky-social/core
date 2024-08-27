package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

var Version = "dev"
var BuildID = "local"

type DBConfig struct {
	Host     string `default:"pg.lab.verysmart.house"`
	User     string `default:"watchedsky-social"`
	Name     string `default:"watchedsky-social"`
	Password string `required:"true"`
}

type BlueskyConfig struct {
	ID          string
	AppPassword string `envconfig:"APP_PASSWORD"`
}

type AlertConfig struct {
	MaxAge       time.Duration `split_keys:"true" default:"372h"`
	HydrationDir string        `split_keys:"true" default:"/app/hydration"`
	Bluesky      BlueskyConfig
}

type WatchedSkyConfig struct {
	DB     DBConfig `required:"true"`
	Alerts AlertConfig
}

var Config *WatchedSkyConfig

func Load() {
	if Config == nil {
		Config = new(WatchedSkyConfig)

		err := envconfig.Process("watchedsky", Config)
		if err != nil {
			panic(err)
		}
	}
}
