package database

import (
	"context"
	"database/sql"
)

type DatabaseFile struct {
	ID               int64  `json:"id"`
	DatabaseName     string `json:"database_name"`
	FileName         string `json:"file_name"`
	FilePhysicalName string `json:"file_physical_name"`
}

const DatabaseFilesQuery string = `
	use master
	select
     ROW_NUMBER() OVER(ORDER BY  files.database_id,files.file_id) AS ID
		,db.name as DatabaseName
		,files.name as FileName
		,files.physical_name as FilePhysicalName
	from 
		sys.master_files as files with(nolock)
	inner join 
		sys.databases as db with(nolock)
	on 
		files.database_id = db.database_id
  `

func GetDatabaseFiles(dbcontext *sql.DB) ([]DatabaseFile, error) {
	ctx := context.Background()
	rows, err := dbcontext.QueryContext(ctx, DatabaseFilesQuery)
	if err != nil {
		return nil, err
	}

	dbFiles := make([]DatabaseFile, 0)
	defer rows.Close()
	for rows.Next() {
		dbf := DatabaseFile{}
		if err := rows.Scan(
			&dbf.ID,
			&dbf.DatabaseName,
			&dbf.FileName,
			&dbf.FilePhysicalName,
		); err != nil {
			return nil, err
		}
		dbFiles = append(dbFiles, dbf)
	}
	return dbFiles, nil
}
