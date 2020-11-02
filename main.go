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
  "time"
	"os"
	"encoding/json"
	"encoding/csv"
  "reflect"
  "strconv"
  "log"
	"path/filepath"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/gobuffalo/packr/v2"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/urfave/cli/v2"
  "github.com/mattn/go-pubsub"
	db "github.com/atoyr/SQLServerGo/database"
)

// SQL DB Information
var sqlserver string
var instance string
var user string
var password string
var database string

// Web
var webport int
var tickRate int64

// Publisher/Subscriber
var ps = pubsub.New()

// Time Layout
var layout = "2006-01-02 15:04:05"

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
  // config setting
  appPath, err := createConfigDirectoryIfNotExists("sqlServerGo")
  if err != nil {
    log.Fatal(err)
  }
  dbFileName := fmt.Sprintf("%s__%s__data.db", sqlserver, instance)
	dbFile := filepath.Join(appPath, dbFileName)
  bolt, err := openDB(dbFile)
  if err != nil {
    log.Fatal(err)
  }
  defer bolt.Close()

  // Check SQL DB Exists
	err = tryDatabaseConnect()
	if err != nil {
		return err
	}

  // create db connection
	d, err := sql.Open("sqlserver", connectionstring())
	if err != nil {
    log.Fatal(err)
	}
	defer d.Close()

  // create and run websocket hub
	fileIOHub := newHub()
	cpuUsedHub := newHub()
	memoryHub := newHub()
	bufferCacheHub := newHub()
	go fileIOHub.run()
	go cpuUsedHub.run()
	go memoryHub.run()
	go bufferCacheHub.run()

  // register subscribe
  ps.Sub(func(dbFileIO []db.DatabaseFileIO) {
    t := time.Now()
    data, _ := json.Marshal(dbFileIO)
    fileIOHub.broadcast <- data
    // update(bolt, "databaseFileIO", []byte(t.Format(time.RFC3339)), []byte(data))
  })
  ps.Sub(func(cpu db.Cpu) {
    t := time.Now()
    data, _ := json.Marshal(cpu)
    cpuUsedHub.broadcast <- data
    // update(bolt, "cpuUsed", []byte(t.Format(time.RFC3339)), []byte(data))
  })
  ps.Sub(func(memory db.Memory) {
    t := time.Now()
    data, _ := json.Marshal(memory)
    memoryHub.broadcast <- data
    // update(bolt, "memory", []byte(t.Format(time.RFC3339)), []byte(data))
  })
  ps.Sub(func(bufferCache db.BufferCache) {
    t := time.Now()
    data, _ := json.Marshal(bufferCache)
    bufferCacheHub.broadcast <- data
    // update(bolt, "bufferCache", []byte(t.Format(time.RFC3339)), []byte(data))
  })


	back := context.Background()
  //tickerCtx, _ := context.WithCancel(back)
  go func(c context.Context){
    t := time.NewTicker(time.Duration(tickRate) * time.Second)
    for {
      select {
      case <-c.Done():
        return
      case <-t.C:
        publishDatabaseFileIO(d)
        publishCpuUsed(d)
        publishMemory(d)
        publishBufferCache(d)
      }
    }
  }(back)

	ec := echo.New()
	ec.Use(middleware.CORS())

	// frontend
	box := packr.New("webapps", "./public")
	ec.GET("/*", echo.WrapHandler(http.StripPrefix("/", http.FileServer(box))))

	// websocket
	ec.GET("/ws/fileio", func(c echo.Context) error {
		serveWs(fileIOHub, c.Response(), c.Request())
		return nil
	})
	ec.GET("/ws/cpu", func(c echo.Context) error {
		serveWs(cpuUsedHub, c.Response(), c.Request())
		return nil
	})
	ec.GET("/ws/memory", func(c echo.Context) error {
		serveWs(memoryHub, c.Response(), c.Request())
		return nil
	})
	ec.GET("/ws/bufferCache", func(c echo.Context) error {
		serveWs(bufferCacheHub, c.Response(), c.Request())
		return nil
	})

	// webapi
	ec.GET("/api/instance", handleInstance)
	ec.GET("/api/serverStatus", handleServerStatus)
	ec.GET("/api/databaseFiles", handleDatabaseFiles)
	ec.GET("/api/cpuUsed", handleCpuUsed)
	ec.GET("/api/memory", handleMemory)
	ec.GET("/api/bufferCache", handleBufferCache)
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
	err = d.Ping()
	if err != nil {
		fmt.Println(connectionstring())
		return err
	}
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

func writeCsv(path string, target interface{}) error{
  file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
  if err != nil {
    return err
  }
  defer file.Close()

  writer := csv.NewWriter(file)
  writer.Write(createCSV(target))
  writer.Flush()
  return nil
}

func createCSVHeader(target interface{}) []string {
  rv := reflect.ValueOf(target)
  rt := rv.Type()
  ret := make([]string,rt.NumField())

  for i := 0; i < rt.NumField(); i++{
    f := rt.Field(i)
    ret[i] = f.Name
  }
  return ret
}

func createCSV(target interface{}) []string {
  rv := reflect.ValueOf(target)
  rt := rv.Type()
  ret := make([]string,rt.NumField())

  for i := 0; i < rt.NumField(); i++{
    f := rv.Field(i)
    switch f.Interface().(type) {
      case string:
        if v, ok := f.Interface().(string); ok {
          ret[i] = v
        } else {
          ret[i] = ""
        }
      case int:
      case int32:
      case int64:
        if v, ok := f.Interface().(int64); ok {
          ret[i] = strconv.FormatInt(v,10)
        } else {
          ret[i] = ""
        }
      case time.Time:
        if v, ok := f.Interface().(time.Time); ok {
          ret[i] = v.Format(layout)
        } else {
          ret[i] = ""
        }
      default :
        ret[i] = ""
    }
  }
  return ret
}
