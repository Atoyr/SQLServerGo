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
	"net/http"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/atoyr/SQLServerGo/database"
	"github.com/gobuffalo/packr/v2"
	"github.com/labstack/echo"
	"github.com/urfave/cli/v2"
)

var server string
var instance string
var user string
var password string
var db string
var port string

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
	}

	app.Action = action
	err := app.Run(os.Args)
	if err != nil {
		log.Println(err)
	}
}

// func serve(urlPrefix string, fs *packr.Box) echo.MiddlewareFunc {
// 	fileserver := http.FileServer(fs)
// 	if urlPrefix != "" {
// 		fileserver = http.StripPrefix(urlPrefix, fileserver)
// 	}
// 	fmt.Println(fileserver)
// 	return func(before echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			err := before(c)
// 			if err != nil {
// 				if c, ok := err.(*echo.HTTPError); !ok || c.Code != http.StatusNotFound {
// 					return err
// 				}
// 			}
//
// 			w, r := c.Response(), c.Request()
// 			fmt.Println(urlPrefix)
// 			fmt.Println(r.URL.Path)
// 			p := strings.TrimPrefix(r.URL.Path, urlPrefix)
// 			s, err := fs.FindString(p)
// 			if err != nil {
// 				fmt.Println("fuga")
// 				fmt.Println(err)
// 			} else {
// 				fmt.Println("hoge")
// 				fmt.Println(s)
// 			}
// 			if fs.Has(p) {
// 				fileserver.ServeHTTP(w, r)
// 				return nil
// 			}
// 			return err
// 		}
// 	}
// }

// 	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
// 		if _, err := b.Open(p); err != nil {
// 			return false
// 		}
// 		return true
// 	}
// 	return false

func action(c *cli.Context) error {
	hub := newHub()
	go hub.run()
	back := context.Background()
	go getFileIO(back, hub)
	ec := echo.New()

	box := packr.New("webapps", "./public")

	// ec.GET("/", echo.WrapHandler(http.FileServer(box)))
	// ec.GET("/*", echo.WrapHandler(http.FileServer(box)))
	ec.GET("/*", echo.WrapHandler(http.StripPrefix("/", http.FileServer(box))))
	ec.GET("/ws", func(c echo.Context) error {
		serveWs(hub, c.Response(), c.Request())
		return nil
	})
	err := ec.Start(port)
	if err != nil {
		return err
	}
	return nil
}

func getFileIO(ctx context.Context, h *Hub) {
	con := database.NewConn(db, instance, server, user, password)
	d, err := sql.Open("sqlserver", con.Connectionstring())
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	t := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			ios, err := database.GetFileIOs(d)
			if err != nil {
				fmt.Println(err)
			} else {
				data, _ := json.Marshal(ios)
				h.broadcast <- data
			}
		}
	}
}
