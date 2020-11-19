package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Client struct {
	db   *sql.DB
	host string
}

type StatementStat struct {
	UserID            int
	DbID              int
	DbName            string `sql:"datname"`
	QueryID           int64
	Query             string
	Calls             int64
	TotalTime         float64 `sql:"total_time"`
	MinTime           float64 `sql:"min_time"`
	MaxTime           float64 `sql:"max_time"`
	MeanTime          float64 `sql:"mean_time"`
	StdDevTime        float64 `sql:"stddev_time"`
	Rows              int64
	SharedBlksHit     int64   `sql:"shared_blks_hit"`
	SharedBlksRead    int64   `sql:"shared_blks_read"`
	SharedBlksDirtied int64   `sql:"shared_blks_dirtied"`
	SharedBlksWritten int64   `sql:"shared_blks_written"`
	LocalBlksHit      int64   `sql:"local_blks_hit"`
	LocalBlksRead     int64   `sql:"local_blks_read"`
	LocalBlksDirtied  int64   `sql:"local_blks_dirtied"`
	LocalBlksWritten  int64   `sql:"local_blks_written"`
	TempBlksRead      int64   `sql:"temp_blks_read"`
	TempBlksWritten   int64   `sql:"temp_blks_written"`
	BlkReadTime       float64 `sql:"blk_read_time"`
	BlkWriteTime      float64 `sql:"blk_write_time"`
	Host              string  `sql:"-"`
}
