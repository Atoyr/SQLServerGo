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
				if len(ios) > 1000 {
					ios = ios[:1000]
				}
				data, _ := json.Marshal(ios)
				h.broadcast <- data
			}
		}
	}
}
