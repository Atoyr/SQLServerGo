package main

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/atoyr/SQLServerGo/database"
	"github.com/labstack/echo"
)

func handleDatabaseFiles(c echo.Context) error {
	con := db.NewConn(database, instance, sqlserver, user, password)
	d, err := sql.Open("sqlserver", con.Connectionstring())
	if err != nil {
		return err
	}
	dbFiles, err := db.GetDatabaseFiles(d)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return c.JSON(http.StatusOK, dbFiles)
}

func handleInstance(c echo.Context) error {
	con := db.NewConn(database, instance, sqlserver, user, password)
	d, err := sql.Open("sqlserver", con.Connectionstring())
	if err != nil {
		return err
	}
	instance, err := db.GetServerProperty(d)
	if err != nil {
		fmt.Println(err)
		fmt.Println(instance)
		return err
	}
	return c.JSON(http.StatusOK, instance)
}

func handleCpuUsed(c echo.Context) error {
	con := db.NewConn(database, instance, sqlserver, user, password)
	d, err := sql.Open("sqlserver", con.Connectionstring())
	if err != nil {
		return err
	}
	cpu, err := db.GetCpuUsed(d)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return c.JSON(http.StatusOK, cpu)
}
