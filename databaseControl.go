package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	db "github.com/atoyr/SQLServerGo/database"
)

func getDatabaseFileIO(ctx context.Context, h *Hub) {
	d, err := sql.Open("sqlserver", connectionstring())
	if err != nil {
		log.Println(err)
		// TODO CONTEXT end
		os.Exit(1)
	}
	defer d.Close()

	t := time.NewTicker(time.Duration(tickRate) * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			ios, err := db.GetDatabaseFileIOs(d)
			if err != nil {
				fmt.Println(err)
			} else {
				data, _ := json.Marshal(ios)
				h.broadcast <- data
			}
		}
	}
}

func getCpuUsed(ctx context.Context, h *Hub) {
	d, err := sql.Open("sqlserver", connectionstring())
	if err != nil {
		log.Println(err)
		// TODO CONTEXT end
		os.Exit(1)
	}
	defer d.Close()

	t := time.NewTicker(time.Duration(tickRate) * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			cu, err := db.GetCpuUsed(d)
			if err != nil {
				fmt.Println(err)
			} else {
				data, _ := json.Marshal(cu)
				h.broadcast <- data
			}
		}
	}
}

func getMemory(ctx context.Context, h *Hub) {
	d, err := sql.Open("sqlserver", connectionstring())
	if err != nil {
		log.Println(err)
		// TODO CONTEXT end
		os.Exit(1)
	}
	defer d.Close()

	t := time.NewTicker(time.Duration(tickRate) * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			cu, err := db.GetMemory(d)
			if err != nil {
				fmt.Println(err)
			} else {
				data, _ := json.Marshal(cu)
				h.broadcast <- data
			}
		}
	}
}

func getBufferCache(ctx context.Context, h *Hub) {
	d, err := sql.Open("sqlserver", connectionstring())
	if err != nil {
		log.Println(err)
		// TODO CONTEXT end
		os.Exit(1)
	}
	defer d.Close()

	t := time.NewTicker(time.Duration(tickRate) * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			cu, err := db.GetBufferCacheRate(d)
			if err != nil {
				fmt.Println(err)
			} else {
				data, _ := json.Marshal(cu)
				h.broadcast <- data
			}
		}
	}
}
