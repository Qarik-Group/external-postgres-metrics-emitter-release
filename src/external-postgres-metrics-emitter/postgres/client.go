package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kisielk/sqlstruct"
	_ "github.com/lib/pq"
	"github.com/starkandwayne/external-postgres-metrics-emitter-release/src/external-postgres-metrics-emitter/config"
)

const enableStatsStatementQuery = "CREATE EXTENSION IF NOT EXISTS pg_stat_statements"
const resetStatsStatementQuery = "SELECT pg_stat_statements_reset()"

func Connect(config config.DatabaseConfig) (*Client, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=postgres sslmode=disable",
		config.Host, config.Port, config.Username, config.Password)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(enableStatsStatementQuery)
	if err != nil {
		return nil, err
	}

	host := fmt.Sprintf("%s:%d", config.Host, config.Port)

	return &Client{db, host}, nil
}

func (c *Client) GetStatsAndReset(ctx context.Context) ([]StatementStat, error) {
	rows, err := c.db.QueryContext(ctx, "SELECT * FROM pg_stat_statements JOIN (SELECT oid, datname FROM pg_database) AS db_name ON pg_stat_statements.dbid = db_name.oid")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	_, err = c.db.ExecContext(ctx, resetStatsStatementQuery)
	if err != nil {
		return nil, err
	}

	stats := make([]StatementStat, 0)
	for rows.Next() {
		var stat StatementStat
		if err := sqlstruct.Scan(&stat, rows); err != nil {
			return nil, err
		}
		if stat.Query == resetStatsStatementQuery {
			continue
		}
		stat.Host = c.host
		stats = append(stats, stat)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}
