package main

import (
	"database/sql"
	"fmt"

	db "github.com/atoyr/SQLServerGo/database"
)

func publishDatabaseFileIO(d *sql.DB){
  ios, err := db.GetDatabaseFileIOs(d)
  if err != nil {
    fmt.Println(err)
  } else {
    ps.Pub(ios)
  }
}

func publishCpuUsed(d *sql.DB){
  cu, err := db.GetCpuUsed(d)
  if err != nil {
    fmt.Println(err)
  } else {
    ps.Pub(cu)
  }
}

func publishMemory(d *sql.DB){
  m, err := db.GetMemory(d)
  if err != nil {
    fmt.Println(err)
  } else {
    ps.Pub(m)
  }
}

func publishBufferCache(d *sql.DB){
  bc, err := db.GetBufferCacheRate(d)
  if err != nil {
    fmt.Println(err)
  } else {
    ps.Pub(bc)
  }
}
