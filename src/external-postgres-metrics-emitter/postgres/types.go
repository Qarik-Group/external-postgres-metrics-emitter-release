package postgres

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type Client struct {
	db *sql.DB
}

// userid              | oid              |           |          |
// dbid                | oid              |           |          |
// queryid             | bigint           |           |          |
// query               | text             |           |          |
// calls               | bigint           |           |          |
// total_time          | double precision |           |          |
// min_time            | double precision |           |          |
// max_time            | double precision |           |          |
// mean_time           | double precision |           |          |
// stddev_time         | double precision |           |          |
// rows                | bigint           |           |          |
// shared_blks_hit     | bigint           |           |          |
// shared_blks_read    | bigint           |           |          |
// shared_blks_dirtied | bigint           |           |          |
// shared_blks_written | bigint           |           |          |
// local_blks_hit      | bigint           |           |          |
// local_blks_read     | bigint           |           |          |
// local_blks_dirtied  | bigint           |           |          |
// local_blks_written  | bigint           |           |          |
// temp_blks_read      | bigint           |           |          |
// temp_blks_written   | bigint           |           |          |
// blk_read_time       | double precision |           |          |
// blk_write_time      | double precision |           |          |
type StatementStat struct {
	Query   string `sql:"query"`
	QueryID int    `sql:"queryid"`
	Calls   int    `sql:"calls"`
	Timestamp         time.Time `sql:"-"`
}
