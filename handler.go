package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/atoyr/SQLServerGo/database"
	"github.com/labstack/echo"
)

func handleDatabaseFiles(c echo.Context) error {
	con := database.NewConn(db, instance, server, user, password)
	d, err := sql.Open("sqlserver", con.Connectionstring())
	if err != nil {
		fmt.Println(err)
		fmt.Println(con)
		return err
	}
	dbFiles, err := database.GetDatabaseFiles(d)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return c.JSON(http.StatusOK, dbFiles)
}
