package database

import (
	"context"
	"database/sql"
)

// BufferCache is SQLServer Buffer Cache
type BufferCache struct {
	BufferCacheRate string `json:"buffer_cache_rate"`
}

const BufferCacheInfoQuery string = `
use master
  SELECT ( opc.cntr_value * 1.0 / bcr.cntr_value ) * 100 AS BufferCacheRate
  FROM sys.dm_os_performance_counters opc
  JOIN (SELECT cntr_value,
  object_name 
  FROM sys.dm_os_performance_counters  
  WHERE counter_name = 'Buffer cache hit ratio base'
  AND object_name LIKE '%Buffer Manager%') bcr 
  ON  opc.OBJECT_NAME = bcr.OBJECT_NAME
  WHERE opc.counter_name = 'Buffer cache hit ratio'
  AND opc.OBJECT_NAME LIKE '%Buffer Manager%'
`

func GetBufferCacheRate(dbcontext *sql.DB) (BufferCache, error) {
	ctx := context.Background()
	row := dbcontext.QueryRowContext(ctx, BufferCacheInfoQuery)

	bf := new(BufferCache)
	if err := row.Scan(
		&bf.BufferCacheRate); err != nil {
		return *bf, err
	}
	return *bf, nil

}
