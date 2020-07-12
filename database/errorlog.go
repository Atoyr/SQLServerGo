package database

import (
	"context"
	"database/sql"
)

// Database is SQLServer instrance database
type ErrorLog struct {
	LogDate     string `json:"log_date"`
	ProcessInfo string `json:"process_info"`
	Text        string `json:"text"`
}

const ErrorLogQuery string = `
  EXEC [sys].[xp_readerrorlog]
`

func GetErrorLogs(dbcontext *sql.DB) ([]ErrorLog, error) {
	ctx := context.Background()
	els := make([]ErrorLog, 0, 0)
	rows, err := dbcontext.QueryContext(ctx, ErrorLogQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		el := ErrorLog{}
		if err := rows.Scan(
			&el.LogDate,
			&el.ProcessInfo,
			&el.Text,
		); err != nil {
			return nil, err
		}
		els = append(els, el)
	}
	return els, nil
}
