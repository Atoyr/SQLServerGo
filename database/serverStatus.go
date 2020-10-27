package database

import (
	"context"
	"database/sql"
	"time"
)

type ServerStatus struct {
	InstanceStartTime   time.Time `json:"InstanceStartTime"`
}

const ServerStatusQuery string = `
    SELECT create_date  AS SQLServerStartTime FROM sys.databases WHERE name = 'tempdb';
  `

func GetServerStatus(dbcontext *sql.DB) (ServerStatus, error) {
	ctx := context.Background()
	row := dbcontext.QueryRowContext(ctx, ServerStatusQuery)

	ss := new(ServerStatus)
	if err := row.Scan(
		&ss.InstanceStartTime,
		); err != nil {
		return *ss, err
	}
	return *ss, nil
}
