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

func GetServerProperty(dbcontext *sql.DB) (ServerProperty, error) {
	ctx := context.Background()
	sp := new(ServerProperty)
	row := dbcontext.QueryRowContext(ctx, ServerInfoQuery)
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
