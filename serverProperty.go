package database

import (
	"context"
	"database/sql"
	"fmt"
)

type ServerProperty struct {
	MachineName         string
	InstanceName        string
	ServerName          string
	ProductVersion      string
	ProductMajorVersion string
	Version             string
	Edition             string
	ProductLevel        string
}

func NewServerProperty(dbcontext *sql.DB) (ServerProperty, error) {
	ctx := context.Background()
	sp := new(ServerProperty)
	tsql := `
    SELECT 
      SERVERPROPERTY('MachineName')
    , SERVERPROPERTY('InstanceName')
    , SERVERPROPERTY('ServerName')
    , SERVERPROPERTY('productversion')
    , SERVERPROPERTY('ProductMajorVersion')
    , SERVERPROPERTY('ProductLevel')
    , SERVERPROPERTY('Edition')
  `
	row := dbcontext.QueryRowContext(ctx, tsql)
	if err := row.Scan(
		sp.MachineName,
		sp.InstanceName,
		sp.ServerName,
		sp.ProductVersion,
		sp.ProductMajorVersion,
		sp.ProductLevel,
		sp.Edition); err != nil {
		return *sp, err
	}
	sp.Version = GetVersion(sp.ProductMajorVersion)
	return *sp, nil
}

func (sp *ServerProperty) String() string {
	return fmt.Sprintf("Database  | %s %s %s \t | %s", sp.Version, sp.ProductLevel, sp.Edition, sp.ServerName)
}
