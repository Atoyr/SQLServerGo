package database

import (
  "fmt"
	"context"
	"database/sql"
)

type TableSize struct {
  TableName string `json:"table_name"`
  Rows int `json:"rows"`
  ReservedBytes int64 `json:"reserved_bytes"`
  DataBytes int64 `json:"data_bytes"`
  IndexBytes int64 `json:"index_bytes"`
  UnusedBytes int64 `json:"unused_bytes"`
}


const TableSizeQuery string = `
SELECT 
   o.name as table_name
  ,row_count as rows
  ,sum(reserved_page_count) * 8 as reserved
  ,data_count * 8 as data
  ,(sum(used_page_count) - data_count) * 8 as  index_size
  ,case when sum(reserved_page_count) > sum(used_page_count) then (sum(reserved_page_count) - sum(used_page_count)) * 8 else 0 end as  unused
FROM sys.dm_db_partition_stats AS DP
join (
  SELECT 
  object_id
  ,in_row_data_page_count + lob_used_page_count + row_overflow_used_page_count as data_count
  FROM sys.dm_db_partition_stats
  WHERE
  index_id < 2
) as data_partition
on DP.object_id = data_partition.object_id
join sys.objects as o
on DP.object_id = o.object_id  
JOIN sys.indexes AS I
ON DP.object_id = I.object_id
AND DP.index_id = I.index_id
WHERE
o.type = N'U'
GROUP BY 
o.name
,row_count
,data_count
order by o.name
`

func GetTableSize(dbcontext *sql.DB, database string) ([]TableSize, error) {
	ctx := context.Background()
	rows, err := dbcontext.QueryContext(ctx, addUseDatabase(database, TableSizeQuery))
	if err != nil {
		return nil, err
	}

	tsizes := make([]TableSize, 0)
	defer rows.Close()
	for rows.Next() {
		ts := TableSize{}
		if err := rows.Scan(
			&ts.TableName,
			&ts.Rows,
			&ts.ReservedBytes,
			&ts.DataBytes,
			&ts.IndexBytes,
			&ts.UnusedBytes,
		); err != nil {
			return nil, err
		}
		 tsizes = append(tsizes, ts)
	}
	return tsizes, nil
}

func addUseDatabase(database, query string) string{
  return fmt.Sprintf("use %s \n %s",database,query)
}
