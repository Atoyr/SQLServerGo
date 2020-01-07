package database

import (
	"context"
	"database/sql"
)

// Request is DB request session
type Request struct {
	Spid                      string
	BlkSpid                   string
	ElapsedSec                int
	DbName                    string
	HostName                  string
	ProgramName               string
	Status                    string
	CommandText               string
	CurrentRunningStmt        string
	TimeSec                   int
	WaitResource              string
	WaitType                  string
	LastWaitType              string
	WaitTimeMsec              int64
	OpenTransactionCount      string
	Command                   string
	PercentComplete           int
	CpuTime                   int
	TransactionIsolationLevel string
	GrantedQueryMemoryKbyte   int
	QueryPlanXml              string
}

func GetRequests(dbcontext *sql.DB) ([]Request, error) {
	ctx := context.Background()
	rs := make([]Request, 0, 10)
	tsql := `
    SELECT 
       der.session_id as spid
      ,der.blocking_session_id as blk_spid
      ,datediff(s, der.start_time, GETDATE()) as elapsed_sec
      ,DB_NAME(der.database_id) AS db_name
      ,des.host_name
      ,des.program_name
      ,der.status -- Status of the request. (background / running / runnable / sleeping / suspended)
      ,dest.text as command_text
      ,REPLACE(REPLACE(REPLACE(SUBSTRING(dest.text, 
      (der.statement_start_offset / 2) + 1, 
      ((CASE der.statement_end_offset
      WHEN -1 THEN DATALENGTH(dest.text)
      ELSE der.statement_end_offset
      END - der.statement_start_offset) / 2) + 1),CHAR(13), ' '), CHAR(10), ' '), CHAR(9), ' ') AS current_running_stmt
      ,datediff(s, der.start_time, GETDATE()) as time_sec
      ,wait_resource 
      ,ISNULL(wait_type,'')
      ,last_wait_type
      ,der.wait_time  as wait_time_ms
      ,der.open_transaction_count
      ,der.command
      ,der.percent_complete
      ,der.cpu_time
      ,(case der.transaction_isolation_level
        when 0 then 'Unspecified'
        when 1 then 'ReadUncomitted'
        when 2 then 'ReadCommitted'
        when 3 then 'Repeatable'
        when 4 then 'Serializable'
        when 5 then 'Snapshot'
      else cast(der.transaction_isolation_level as varchar) end) as transaction_isolation_level
      ,der.granted_query_memory * 8 as granted_query_memory_kb 
      ,ISNULL(deqp.query_plan ,'')
  FROM
      sys.dm_exec_requests der
  JOIN sys.dm_exec_sessions des ON des.session_id = der.session_id
  OUTER APPLY sys.dm_exec_sql_text(sql_handle) AS dest
  OUTER APPLY sys.dm_exec_query_plan(plan_handle) AS deqp
  WHERE
      des.is_user_process = 1`
	tsql = tsql + `
  ORDER BY
    datediff(s, der.start_time, GETDATE()) DESC
  `
	rows, err := dbcontext.QueryContext(ctx, tsql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		r := new(Request)
		if err := rows.Scan(
			r.Spid,
			r.BlkSpid,
			r.ElapsedSec,
			r.DbName,
			r.HostName,
			r.ProgramName,
			r.Status,
			r.CommandText,
			r.CurrentRunningStmt,
			r.TimeSec,
			r.WaitResource,
			r.WaitType,
			r.LastWaitType,
			r.WaitTimeMsec,
			r.OpenTransactionCount,
			r.Command,
			r.PercentComplete,
			r.CpuTime,
			r.TransactionIsolationLevel,
			r.GrantedQueryMemoryKbyte,
			r.QueryPlanXml,
		); err != nil {
			return nil, err
		}
		rs = append(rs, *r)
	}
	return rs, nil
}
