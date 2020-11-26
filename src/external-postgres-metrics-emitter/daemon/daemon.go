package daemon

import (
	"context"
	"errors"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/starkandwayne/external-postgres-metrics-emitter-release/src/external-postgres-metrics-emitter/config"
	"github.com/starkandwayne/external-postgres-metrics-emitter-release/src/external-postgres-metrics-emitter/forwarder"
	"github.com/starkandwayne/external-postgres-metrics-emitter-release/src/external-postgres-metrics-emitter/postgres"
)

func Run(logger lager.Logger, args []string, stop chan bool) {
	if len(args) < 2 {
		logger.Fatal("Missing argument - specify path to config file", errors.New("Missing config file path"))
	}

	configFilePath := args[1]

	conf, err := config.LoadConfig(configFilePath)
	if err != nil {
		logger.Fatal("Reading config file", err, lager.Data{
			"emitter-config-file-path": configFilePath,
		})
	}

	metricsClient, err := forwarder.NewMetricForwarder(logger, &conf)
	if err != nil {
		logger.Fatal("Couldn't create metric-forwarder", err)
	}

	var dbs []*postgres.Client
	for _, dbConf := range conf.DatabaseConfigs {
		db, err := postgres.Connect(dbConf)
		if err != nil {
			logger.Fatal("Failed to connect to database", err)
		}
		dbs = append(dbs, db)
	}

	ctx := context.Background()

	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				ctx, _ := context.WithCancel(ctx)

				for _, db := range dbs {
					stats, err := db.GetStatsAndReset(ctx)
					if err != nil {
						logger.Error("Failed to get stats from database", err)
					}

					for _, stat := range stats {
						metricsClient.EmitMetric(&stat)
					}
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	<-stop
	close(quit)
}
