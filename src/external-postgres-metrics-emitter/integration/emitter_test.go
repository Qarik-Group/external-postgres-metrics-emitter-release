package integration_test

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/starkandwayne/external-postgres-metrics-emitter-release/src/external-postgres-metrics-emitter/config"
)

var _ = Describe("Emitter", func() {
	var db *sql.DB
	var conf config.Config

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
	})

	FIt("End to end", func() {
		result, err := db.Exec("SELECT * FROM pg_stat_statements LIMIT 1;")
		Expect(err).ToNot(HaveOccurred())
		rows, err := result.RowsAffected()
		Expect(err).ToNot(HaveOccurred())
		Expect(rows).To(Equal(int64(1)))
		// run 1 scrape
		// assert metron dumped json
		// assert our query show up in dump
	})

	AfterEach(func() {
		var err error
		err = db.Close()
		Expect(err).ToNot(HaveOccurred())
	})
})
