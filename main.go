// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/gobuffalo/packr/v2"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/urfave/cli/v2"
)

var sqlserver string
var instance string
var user string
var password string
var database string

var webport int
var tickRate int64

func main() {
	app := new(cli.App)
	app.Name = "SQLServer Tools"
	app.Version = "0.1.0"
	app.Usage = "run apps and access http://localhost:<httpport>"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "server",
			Aliases:     []string{"s"},
			Value:       "",
			Usage:       "SQLServer Server Name",
			EnvVars:     []string{"DBSERVER"},
			Destination: &sqlserver,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "instance",
			Aliases:     []string{"i"},
			Value:       "",
			Usage:       "SQLServer Server Instance Name",
			EnvVars:     []string{"DBINSTANCE"},
			Destination: &instance,
		},
		&cli.StringFlag{
			Name:        "user",
			Aliases:     []string{"u"},
			Value:       "sa",
			Usage:       "SQLServer Server User",
			EnvVars:     []string{"DBUSER"},
			Destination: &user,
		},
		&cli.StringFlag{
			Name:        "password",
			Aliases:     []string{"p"},
			Value:       "",
			Usage:       "SQLServer Server Password",
			EnvVars:     []string{"DBPASS"},
			Destination: &password,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "database",
			Aliases:     []string{"d"},
			Value:       "master",
			Usage:       "SQLServer Server using database",
			Destination: &database,
		},
		&cli.IntFlag{
			Name:        "httpport",
			Aliases:     []string{"hp"},
			Value:       8080,
			Usage:       "http access port no",
			Destination: &webport,
		},
		&cli.Int64Flag{
			Name:        "tickrate",
			Aliases:     []string{"tr"},
			Value:       1,
			Usage:       "performance data tick rate",
			Destination: &tickRate,
		},
	}

	app.Action = action
	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("%s %s", failstring, err)
	}
}

func action(c *cli.Context) error {
	err := tryDatabaseConnect()
	if err != nil {
		return err
	}
	hub := newHub()
	go hub.run()
	back := context.Background()
	go getDatabaseFileIO(back, hub)
	ec := echo.New()
	ec.Use(middleware.CORS())

	// frontend
	box := packr.New("webapps", "./public")
	ec.GET("/*", echo.WrapHandler(http.StripPrefix("/", http.FileServer(box))))

	// websocket
	ec.GET("/ws/fileio", func(c echo.Context) error {
		serveWs(hub, c.Response(), c.Request())
		return nil
	})

	// webapi
	ec.GET("/api/instance", handleInstance)
	ec.GET("/api/databaseFiles", handleDatabaseFiles)
	ec.GET("/api/cpuUsed", handleCpuUsed)
	ec.GET("/api/errorlogs", handleErrorLogs)

	ec.HideBanner = true

	err = ec.Start(fmt.Sprintf(":%d", webport))
	if err != nil {
		return err
	}
	return nil
}

const okstring = "[\x1b[32m OK \x1b[0m]"
const failstring = "[\x1b[31mFAIL\x1b[0m]"
const infostring = "[\x1b[36mINFO\x1b[0m]"

func connectionstring() string {
	var ret = make([]byte, 0, 1024)
	ret = append(ret, "server="...)
	ret = append(ret, sqlserver...)
	if instance != "" {
		ret = append(ret, "\\"...)
		ret = append(ret, instance...)
	}
	ret = append(ret, ";user id="...)
	ret = append(ret, user...)
	ret = append(ret, ";password="...)
	ret = append(ret, password...)
	ret = append(ret, ";database="...)
	ret = append(ret, database...)
	return string(ret)
}

func tryDatabaseConnect() error {
	fmt.Println(connectionstring())
	d, err := sql.Open("sqlserver", connectionstring())
	if err != nil {
		return err
	}
	defer d.Close()
	// err = d.Ping()
	// if err != nil {
	// 	fmt.Println(connectionstring())
	// 	return err
	// }
	return nil
}

func check() (error, bool) {
	// Check DB Status
	d, err := sql.Open("sqlserver", connectionstring())
	if err != nil {
		fmt.Printf("%s SQL Server Connection %s\\%s\n", failstring, sqlserver, instance)
		return err, false
	}
	defer d.Close()
	err = d.Ping()
	if err != nil {
		fmt.Printf("%s SQL Server Connection %s\\%s\n", failstring, sqlserver, instance)
		return err, false
	} else {
		fmt.Printf("%s SQL Server Connection %s\\%s\n", okstring, sqlserver, instance)
	}

	// Check HTTP Listen
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", webport))
	if err == nil {
		fmt.Printf("%s HTTP Listen Port %d\n", okstring, webport)
		defer l.Close()
	} else {
		fmt.Printf("%s HTTP Listen Port %d\n", failstring, webport)
		return err, false
	}
	return nil, true
}
