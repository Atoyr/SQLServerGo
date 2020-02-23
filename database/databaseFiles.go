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
