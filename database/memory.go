package database

import (
	"context"
	"database/sql"
)

// Database is SQLServer instrance database
type Memory struct {
	PhysicalMemory    int    `json:"physical_memory"`
	UsedMemory        int    `json:"used_memory"`
	AvailableMemory   int    `json:"available_memory"`
	TotalPageFile     int    `json:"total_page_file"`
	UsedPageFile      int    `json:"used_page_file"`
	AvailablePageFile int    `json:"available_page_file"`
	SystemChace       int    `json:"system_chace"`
	SystemMemoryState string `json:"system_memory_state"`
}

const MemoryInfoQuery string = `
  use master
  SELECT 
    total_physical_memory_kb/1024 AS [Physical Memory (MB)], 
    available_physical_memory_kb/1024 AS [Available Memory (MB)], 
    total_page_file_kb/1024 AS [Total Page File (MB)], 
    available_page_file_kb/1024 AS [Available Page File (MB)], 
    system_cache_kb/1024 AS [System Cache (MB)],
    system_memory_state_desc AS [System Memory State]
  FROM sys.dm_os_sys_memory WITH (NOLOCK) OPTION (RECOMPILE);
`

func GetMemory(dbcontext *sql.DB) (Memory, error) {
	ctx := context.Background()
	row := dbcontext.QueryRowContext(ctx, MemoryInfoQuery)

	m := new(Memory)
	if err := row.Scan(
		&m.PhysicalMemory,
		&m.AvailableMemory,
		&m.TotalPageFile,
		&m.AvailablePageFile,
		&m.SystemChace,
		&m.SystemMemoryState); err != nil {
		return *m, err
	}
	m.UsedMemory = m.PhysicalMemory - m.AvailableMemory
	m.UsedPageFile = m.TotalPageFile - m.AvailablePageFile
	return *m, nil

}
