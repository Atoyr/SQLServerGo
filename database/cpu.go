package database

import (
	"context"
	"database/sql"
)

const CpuUsedQuery string = `
  SELECT TOP 1
      record.value('(./Record/@id)[1]', 'int') AS record_id
    , record.value('(./Record/SchedulerMonitorEvent/SystemHealth/SystemIdle)[1]', 'int') AS [SystemIdle]
    , record.value('(./Record/SchedulerMonitorEvent/SystemHealth/ProcessUtilization)[1]', 'int') AS [SQLProcessUtilization]
    , [timestamp]
  FROM (
    SELECT [timestamp]
    , convert(XML, record) AS [record]
    FROM sys.dm_os_ring_buffers with(nolock)
    WHERE ring_buffer_type = N'RING_BUFFER_SCHEDULER_MONITOR'
    AND record LIKE '%<systemHealth>%'
  ) AS x
  ORDER BY timestamp desc
`

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
