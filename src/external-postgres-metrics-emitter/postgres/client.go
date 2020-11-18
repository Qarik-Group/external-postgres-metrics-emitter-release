package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kisielk/sqlstruct"
	_ "github.com/lib/pq"
	"github.com/starkandwayne/external-postgres-metrics-emitter-release/src/external-postgres-metrics-emitter/config"
)

const resetStatsStatementQuery = "SELECT pg_stat_statements_reset()"

func Connect(config config.DatabaseConfig) (*Client, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=postgres sslmode=disable",
		config.Host, config.Port, config.Username, config.Password)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return &Client{db}, nil
}

func (c *Client) GetStatsAndReset(ctx context.Context) ([]StatementStat, error) {
	rows, err := c.db.QueryContext(ctx, fmt.Sprintf(
		"SELECT %s FROM pg_stat_statements",
		sqlstruct.Columns(StatementStat{})))
	if err != nil {
		return nil, err
	}

	_, err = c.db.QueryContext(ctx, resetStatsStatementQuery)
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
		stats = append(stats, stat)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}
