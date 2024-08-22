package alerts

import (
	"context"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/watchedsky-social/core/internal/config"
	"github.com/watchedsky-social/core/internal/database/query"
)

func CleanOldAlerts(ctx context.Context) error {
	var cfg config.AlertConfig
	if err := envconfig.Process("vacuum", &cfg); err != nil {
		return err
	}

	dao := query.Alert
	r, err := dao.WithContext(ctx).Where(dao.Sent.Lt(time.Now().Add(-cfg.MaxAge))).Delete()
	if err != nil {
		return err
	}

	log.Printf("Deleted %d old alerts", r.RowsAffected)
	return nil
}
