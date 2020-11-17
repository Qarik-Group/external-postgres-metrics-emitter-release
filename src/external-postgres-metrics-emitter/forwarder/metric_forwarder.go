package forwarder

import (
	"code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/lager"
	"github.com/starkandwayne/external-postgres-metrics-emitter-release/src/external-postgres-metrics-emitter/config"
)

type MetricForwarder interface {
	// EmitMetric(*management.QueueInfo)
}

type metricForwarder struct {
	client *loggregator.IngressClient
	logger lager.Logger
}

const METRICS_FORWARDER_ORIGIN = "autoscaler_metrics_forwarder"

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

// func (mf *metricForwarder) EmitMetric(info *postgres.QueueInfo) {
// 	mf.logger.Debug("custom-metric-emit-request-received:", lager.Data{"info": info})

// 	options := []loggregator.EmitGaugeOption{
// 		loggregator.WithGaugeValue(strings.Replace(info.Name, "-", "_", -1)+"_messages_ready", info.MessagesReady, "msgs"),
// 	}
// 	mf.client.EmitGauge(options...)
// }
