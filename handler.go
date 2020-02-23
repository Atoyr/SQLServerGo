package main

import (
	"database/sql"
	"net/http"

	"github.com/atoyr/SQLServerGo/database"
	"github.com/labstack/echo"
)

func handleDatabaseFiles(c echo.Context) error {
	con := database.NewConn(db, instance, server, user, password)
	d, err := sql.Open("sqlserver", con.Connectionstring())
	if err != nil {
    return err
	}
  dbFiles ,err := database.GetDatabaseFiles(d)
  if err != nil {
    return err
  }
  return c.JSON(http.StatusOK, dbFiles)
}
