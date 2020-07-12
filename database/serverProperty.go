package database

import (
	"context"
	"database/sql"
	"fmt"
)

type ServerProperty struct {
	MachineName         string `json:"machineName"`
	InstanceName        string `json:"instanceName"`
	ServerName          string `json:"serverName"`
	ProductVersion      string `json:"productVersion"`
	ProductMajorVersion string `json:"productMajorVersion"`
	Version             string `json:"version"`
	Edition             string `json:"edition"`
	ProductLevel        string `json:"productLevel"`
}

const ServerInfoQuery string = `
    SELECT 
      SERVERPROPERTY('MachineName')
    , SERVERPROPERTY('InstanceName')
    , SERVERPROPERTY('ServerName')
    , SERVERPROPERTY('productversion')
    , SERVERPROPERTY('ProductMajorVersion')
    , SERVERPROPERTY('ProductLevel')
    , SERVERPROPERTY('Edition')
  `

func GetServerProperty(dbcontext *sql.DB) (ServerProperty, error) {
	ctx := context.Background()
	row := dbcontext.QueryRowContext(ctx, ServerInfoQuery)

	sp := new(ServerProperty)
	if err := row.Scan(
		&sp.MachineName,
		&sp.InstanceName,
		&sp.ServerName,
		&sp.ProductVersion,
		&sp.ProductMajorVersion,
		&sp.ProductLevel,
		&sp.Edition); err != nil {
		return *sp, err
	}
	sp.Version = GetVersion(sp.ProductMajorVersion)
	return *sp, nil
}

func (sp *ServerProperty) String() string {
	return fmt.Sprintf("Database  | %s %s %s \t | %s", sp.Version, sp.ProductLevel, sp.Edition, sp.ServerName)
}
