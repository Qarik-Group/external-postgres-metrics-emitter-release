package main

import (
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/starkandwayne/external-postgres-metrics-emitter-release/src/external-postgres-metrics-emitter/deamon"
)

func main() {
	logger := lager.NewLogger("external-postgres-metrics-emitter")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	deamon.Run(logger, os.Args)
}
