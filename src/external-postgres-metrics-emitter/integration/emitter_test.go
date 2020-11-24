package integration_test

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/lager"

	"github.com/starkandwayne/external-postgres-metrics-emitter-release/src/external-postgres-metrics-emitter/config"
	"github.com/starkandwayne/external-postgres-metrics-emitter-release/src/external-postgres-metrics-emitter/daemon"
	"github.com/starkandwayne/external-postgres-metrics-emitter-release/src/external-postgres-metrics-emitter/integration"
)

var _ = Describe("Emitter", func() {
	var (
		db          *sql.DB
		conf        config.Config
		loggregator *integration.DummyLoggregator
		logger      lager.Logger
		stop        chan bool
	)

	BeforeEach(func() {
		var err error
		conf, err = config.LoadConfig("./assets/config.yml")
		Expect(err).ToNot(HaveOccurred())
		dbConf := conf.DatabaseConfig

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=postgres sslmode=disable",
			dbConf.Host, dbConf.Port, dbConf.Username, dbConf.Password)
		db, err = sql.Open("postgres", psqlInfo)
		Expect(err).ToNot(HaveOccurred())
		loggregator, err = integration.NewDummyLoggregator()
		Expect(err).ToNot(HaveOccurred())
		go loggregator.Start()
		logger = lager.NewLogger("external-postgres-metrics-emitter")
		logger.RegisterSink(lager.NewWriterSink(GinkgoWriter, lager.DEBUG))
		logger.RegisterSink(lager.NewWriterSink(GinkgoWriter, lager.ERROR))

		stop = make(chan bool, 1)
	})

	It("End to end", func() {
		result, err := db.Exec("SELECT * FROM pg_stat_statements LIMIT 1;")
		Expect(err).ToNot(HaveOccurred())
		rows, err := result.RowsAffected()
		Expect(err).ToNot(HaveOccurred())
		Expect(rows).To(Equal(int64(1)))

		go daemon.Run(logger, []string{"daemon", "./assets/config.yml"}, stop)

		Eventually(func() []string {
			return loggregator.Received()
		}, "15s").Should(HaveLen(2))
	})

	AfterEach(func() {
		var err error
		err = db.Close()
		Expect(err).ToNot(HaveOccurred())
		loggregator.Stop()
		stop <- true
	})
})
