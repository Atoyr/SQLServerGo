package database

import (
	"context"
	"database/sql"
)

type Cpu struct {
	ID                      int64   `json:"id"`
	SystemIdle              int64   `json:"system_idle"`
	SqlProcessUtilization   int64   `json:"sql_process_utilization"`
	OtherProcessUtilization int64   `json:"other_process_utilization"`
	Timestamp               float64 `json:"timestamp"`
}

func GetCpuUsed(dbcontext *sql.DB) (Cpu, error) {
	ctx := context.Background()
	row := dbcontext.QueryRowContext(ctx, CpuUsedQuery)

	c := new(Cpu)
	if err := row.Scan(
		&c.ID,
		&c.SystemIdle,
		&c.SqlProcessUtilization,
		&c.Timestamp,
	); err != nil {
		return *c, err
	}
	c.OtherProcessUtilization = 100 - (c.SystemIdle + c.SqlProcessUtilization)
	return *c, nil
}
