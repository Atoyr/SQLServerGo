// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/atoyr/SQLServerGo/database"
	"github.com/gobuffalo/packr/v2"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/urfave/cli/v2"
)

var server string
var instance string
var user string
var password string
var db string
var port string
var tickRate int64

func main() {
	app := new(cli.App)
	app.Name = "SQLServer Tools"
	app.Usage = "run apps and access http://localhost:8080"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "server",
			Aliases:     []string{"s"},
			Value:       "",
			Usage:       "SQLServer Server Name",
			EnvVars:     []string{"DBSERVER"},
			Destination: &server,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "instance",
			Aliases:     []string{"i"},
			Value:       "",
			Usage:       "SQLServer Server Instance Name",
			EnvVars:     []string{"DBINSTANCE"},
			Destination: &instance,
			Required:    true,
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
			Destination: &db,
		},
		&cli.StringFlag{
			Name:        "httpport",
			Aliases:     []string{"hp"},
			Value:       ":8080",
			Usage:       "http access port no",
			Destination: &port,
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
		log.Println(err)
	}
}

func action(c *cli.Context) error {
	err, ok := check()
	if ok == false {
		return err
	}
	hub := newHub()
	go hub.run()
	back := context.Background()
	go getDatabaseFileIO(back, hub)
	ec := echo.New()
	ec.Use(middleware.CORS())

	box := packr.New("webapps", "./public")

	ec.GET("/*", echo.WrapHandler(http.StripPrefix("/", http.FileServer(box))))
	ec.GET("/ws", func(c echo.Context) error {
		serveWs(hub, c.Response(), c.Request())
		return nil
	})
	ec.GET("/api/databaseFiles", handleDatabaseFiles)

	ec.HideBanner = true
	err = ec.Start(port)
	if err != nil {
		return err
	}
	return nil
}

var okstring = "[\x1b[32m OK \x1b[0m]"
var failstring = "[\x1b[31mFAIL\x1b[0m]"

func check() (error, bool) {
	// Check DB Status
	con := database.NewConn(db, instance, server, user, password)
	d, err := sql.Open("sqlserver", con.Connectionstring())
	if err != nil {
		fmt.Printf("%s SQL Server Connection %s\\%s\n", failstring, server, instance)
		return err, false
	}
	defer d.Close()
	err = d.Ping()
	if err != nil {
		fmt.Printf("%s SQL Server Connection %s\\%s\n", failstring, server, instance)
		return err, false
	} else {
		fmt.Printf("%s SQL Server Connection %s\\%s\n", okstring, server, instance)
	}

	// Check HTTP Listen
	if port[0] != ':' {
		port = fmt.Sprintf(":%s", port)
	}
	l, err := net.Listen("tcp", port)
	if err == nil {
		fmt.Printf("%s HTTP Listen Port %s\n", okstring, port)
		defer l.Close()
	} else {
		fmt.Printf("%s HTTP Listen Port %s\n", failstring, port)
		return err, false
	}
	return nil, true
}

func getDatabaseFileIO(ctx context.Context, h *Hub) {
	con := database.NewConn(db, instance, server, user, password)
	d, err := sql.Open("sqlserver", con.Connectionstring())
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer d.Close()
	t := time.NewTicker(time.Duration(tickRate) * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			ios, err := database.GetDatabaseFileIOs(d)
			if err != nil {
				fmt.Println(err)
			} else {
				if len(ios) > 1000 {
					ios = ios[:1000]
				}
				data, _ := json.Marshal(ios)
				h.broadcast <- data
			}
		}
	}
}
