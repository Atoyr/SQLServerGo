package database

import (
	"context"
	"database/sql"
	"time"
)

type DatabaseFileIO struct {
	ID                  int64     `json:"id"`
  DatabaseID          int64     `json:"database_id"`
  FileID              int64     `json:"file_id"`
	Datetime            time.Time `json:"datetime"`
	DatabaseName        string    `json:"database_name"`
	FileName            string    `json:"file_name"`
	FilePhysicalName    string    `json:"file_physical_name"`
	FileSizeOnDiskBytes int64     `json:"file_size_on_disk_bytes"`
	ReadBytesPerSec     float64   `json:"read_bytes_per_sec"`
	WriteBytesPerSec    float64   `json:"write_bytes_per_sec"`
}

type fileState struct {
	ID                  int64
  DatabaseID          int64
  FileID              int64
	DatabaseName        string
	FileName            string
	FilePhysicalName    string
	FileSizeOnDiskBytes int64
	SampleMs            int64
	IoStall             int64
	NumOfReads          int64
	NumOfBytesRead      int64
	IoStallReadMs       int64
	NumOfWrites         int64
	NumOfBytesWritten   int64
	IoStallWriteMs      int64
}

var beforeFileState map[string]fileState

const DatabaseFileIOQuery string = `
	use master
	select
     ROW_NUMBER() OVER(ORDER BY  files.database_id,files.file_id) AS ID
    ,files.database_id as DatabaseID
    ,files.file_id as FileID
		,db.name as DatabaseName
		,files.name as FileName
		,files.physical_name as FilePhysicalName
		,io.size_on_disk_bytes as FileSizeOnDiskBytes
		,io.sample_ms as SampleMs
		,io.io_stall as IoStall
		,io.num_of_reads as NumOfReads
		,io.num_of_bytes_read as NumOfBytesRead
		,io.io_stall_read_ms as IoStallReadMs
		,io.num_of_writes as NumOfWrites
		,io.num_of_bytes_written as NumOfBytesWritten
		,io.io_stall_write_ms as IoStallWriteMs
	from 
		sys.dm_io_virtual_file_stats (null,null ) as io
	inner join 
		sys.databases as db with(nolock)
	on 
		io.database_id = db.database_id
	inner join 
		sys.master_files as files with(nolock)
	on 
		io.database_id = files.database_id
	and io.file_id = files.file_id
  `

// GetDataseFileIOs is getting database file io with query.
func GetDatabaseFileIOs(dbcontext *sql.DB) ([]DatabaseFileIO, error) {
	if beforeFileState == nil {
		beforeFileState = map[string]fileState{}
	}

	ctx := context.Background()
	datetime := time.Now().Truncate(time.Second)

	rows, err := dbcontext.QueryContext(ctx, DatabaseFileIOQuery)
	if err != nil {
		return nil, err
	}

	fileIOs := make([]DatabaseFileIO, 0)
	defer rows.Close()
	for rows.Next() {
		fs := fileState{}
		if err := rows.Scan(
			&fs.ID,
      &fs.DatabaseID,
      &fs.FileID,
			&fs.DatabaseName,
			&fs.FileName,
			&fs.FilePhysicalName,
			&fs.FileSizeOnDiskBytes,
			&fs.SampleMs,
			&fs.IoStall,
			&fs.NumOfReads,
			&fs.NumOfBytesRead,
			&fs.IoStallReadMs,
			&fs.NumOfWrites,
			&fs.NumOfBytesWritten,
			&fs.IoStallWriteMs,
		); err != nil {
			return nil, err
		}
		if bfs, ok := beforeFileState[fs.FilePhysicalName]; ok {
			// 1秒あたりの割合計算用の係数を割り出し
			rangeMs := fs.SampleMs - bfs.SampleMs
			var mag float64 = 0
			if rangeMs > 0 {
				mag = float64(1000.00 / rangeMs)
			}

			fileIO := DatabaseFileIO{}
			fileIO.ID = fs.ID
      fileIO.DatabaseID = fs.DatabaseID
      fileIO.FileID = fs.FileID
			fileIO.Datetime = datetime
			fileIO.DatabaseName = fs.DatabaseName
			fileIO.FileName = fs.FileName
			fileIO.FilePhysicalName = fs.FilePhysicalName
			fileIO.FileSizeOnDiskBytes = fs.FileSizeOnDiskBytes
			fileIO.ReadBytesPerSec = float64(fs.NumOfBytesRead-bfs.NumOfBytesRead) * mag
			fileIO.WriteBytesPerSec = float64(fs.NumOfBytesWritten-bfs.NumOfBytesWritten) * mag
			fileIOs = append(fileIOs, fileIO)
		}
		beforeFileState[fs.FilePhysicalName] = fs
	}
	return fileIOs, nil
}
