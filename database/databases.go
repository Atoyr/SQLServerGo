package database

import (
	"context"
	"database/sql"
	"fmt"
)

// Database is SQLServer instrance database
type Database struct {
	Name          string `json:"name"`
	Id            int    `json:"id"`
	Status        string `json:"status"`
	RecoveryModel string `json:"recovery_model"`
}

func NewDatabases(dbcontext *sql.DB) ([]Database, error) {
	ctx := context.Background()
	ds := make([]Database, 0, 10)
	rows, err := dbcontext.QueryContext(ctx, DatabasesQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		d := Database{}
		if err := rows.Scan(
			&d.Name,
			&d.Id,
			&d.Status,
			&d.RecoveryModel,
		); err != nil {
			return nil, err
		}
		ds = append(ds, d)
	}
	return ds, nil
}

func (d *Database) String() string {
	return fmt.Sprintf("%s %d %s %s", d.Name, d.Id, d.Status, d.RecoveryModel)
}
