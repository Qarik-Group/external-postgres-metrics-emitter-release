package forwarder

import (
	"strconv"

	"code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/lager"
	"github.com/starkandwayne/external-postgres-metrics-emitter-release/src/external-postgres-metrics-emitter/config"
	"github.com/starkandwayne/external-postgres-metrics-emitter-release/src/external-postgres-metrics-emitter/postgres"
)

type metricForwarder struct {
	client *loggregator.IngressClient
	logger lager.Logger
}

const METRICS_FORWARDER_ORIGIN = "external_postgres_metrics_emitter"

func NewMetricForwarder(logger lager.Logger, conf *config.Config) (*metricForwarder, error) {
	tlsConfig, err := loggregator.NewIngressTLSConfig(
		conf.LoggregatorConfig.TLS.CACertFile,
		conf.LoggregatorConfig.TLS.CertFile,
		conf.LoggregatorConfig.TLS.KeyFile,
	)
	if err != nil {
		logger.Error("could-not-create-TLS-config", err, lager.Data{"config": conf})
		return &metricForwarder{}, err
	}

	client, err := loggregator.NewIngressClient(
		tlsConfig,
		loggregator.WithAddr(conf.LoggregatorConfig.MetronAddress),
		loggregator.WithTag("origin", METRICS_FORWARDER_ORIGIN),
		loggregator.WithLogger(newLoggregatorGRPCLogger(logger.Session("loggregator"))),
	)

	if err != nil {
		logger.Error("could-not-create-loggregator-client", err, lager.Data{"config": conf})
		return &metricForwarder{}, err
	}

	return &metricForwarder{
		client: client,
		logger: logger,
	}, nil
}

func (mf *metricForwarder) EmitMetric(stat *postgres.StatementStat) {
	mf.logger.Debug("custom-metric-emit-request-received:", lager.Data{"stat": stat})

	options := []loggregator.EmitGaugeOption{
		loggregator.WithEnvelopeTag("query", stat.Query),
		loggregator.WithEnvelopeTag("query_id", strconv.FormatInt(stat.QueryID, 10)),
		loggregator.WithEnvelopeTag("source_host", stat.Host),
		loggregator.WithEnvelopeTag("source_db", stat.DbName),

		loggregator.WithGaugeValue("calls", float64(stat.Calls), "calls"),
		loggregator.WithGaugeValue("total_time", stat.TotalTime, "ms"),
		loggregator.WithGaugeValue("min_time", stat.MinTime, "ms"),
		loggregator.WithGaugeValue("max_time", stat.MaxTime, "ms"),
		loggregator.WithGaugeValue("mean_time", stat.MeanTime, "ms"),
		loggregator.WithGaugeValue("rows", float64(stat.Rows), "rows"),
		loggregator.WithGaugeValue("shared_blks_hit", float64(stat.SharedBlksHit), "blks"),
		loggregator.WithGaugeValue("shared_blks_read", float64(stat.SharedBlksRead), "blks"),
		loggregator.WithGaugeValue("shared_blks_dirtied", float64(stat.SharedBlksDirtied), "blks"),
		loggregator.WithGaugeValue("shared_blks_written", float64(stat.SharedBlksWritten), "blks"),
		loggregator.WithGaugeValue("local_blks_hit", float64(stat.LocalBlksHit), "blks"),
		loggregator.WithGaugeValue("local_blks_read", float64(stat.LocalBlksRead), "blks"),
		loggregator.WithGaugeValue("local_blks_dirtied", float64(stat.LocalBlksDirtied), "blks"),
		loggregator.WithGaugeValue("local_blks_written", float64(stat.LocalBlksWritten), "blks"),
		loggregator.WithGaugeValue("temp_blks_read", float64(stat.TempBlksRead), "blks"),
		loggregator.WithGaugeValue("temp_blks_written", float64(stat.TempBlksWritten), "blks"),

		// Total time the statement spent reading blocks, in milliseconds (if track_io_timing is enabled, otherwise zero)
		// https://www.postgresql.org/docs/11/runtime-config-statistics.html#GUC-TRACK-IO-TIMING
		// This parameter is off by default, because it will repeatedly query the operating system for the current time,
		// which may cause significant overhead on some platforms
		// loggregator.WithGaugeValue("blk_read_time", stat.BlkReadTime, "ms"),
		// loggregator.WithGaugeValue("blk_write_time", stat.BlkWriteTime, "ms"),
	}
	mf.client.EmitGauge(options...)
}
