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

	config, err := config.LoadConfig(configFilePath)
	if err != nil {
		logger.Fatal("Reading config file", err, lager.Data{
			"emitter-config-file-path": configFilePath,
		})
	}

	metricsClient, err := forwarder.NewMetricForwarder(logger, &config)
	if err != nil {
		logger.Fatal("Couldn't create metric-forwarder", err)
	}

	db, err := postgres.Connect(config.DatabaseConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database", err)
	}

	ctx := context.Background()

	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				ctx, _ := context.WithCancel(ctx)

				stats, err := db.GetStatsAndReset(ctx)
				if err != nil {
					logger.Error("Failed to get stats from database", err)
				}

				for _, stat := range stats {
					metricsClient.EmitMetric(&stat)
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
