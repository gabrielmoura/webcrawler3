package tasks

import (
	"context"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"webCrawler3/infra/db/repoDb"
	"webCrawler3/infra/log"
)

func Handle(ctx context.Context) {
	c := cron.New()
	c.AddFunc("@every 1m", func() {
		log.Logger.Info("Running health check")
		errs := repoDb.PrepareHealth(ctx) // Do something
		if len(errs) > 0 {
			for _, err := range errs {
				log.Logger.Error("Error running health check", zap.Error(err))
			}
		}
	})
	c.Start()
}
