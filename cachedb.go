package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func OpenCacheDatabase() {
	dbFile := filepath.Join(LocalAppData, "Roblox", "rbx-storage.db")
	fmt.Println("looking for cache db @ " + dbFile)
	if  _, err := os.Stat(dbFile); os.IsNotExist(err) {
		FatalErrorStr("cache db does not exist?")
	}
	fmt.Println("db found  ok")
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		FatalError(err)
	}

	rows, err := db.Query("SELECT content FROM files")
	if err != nil {
		FatalError(err)
	}

	q := 0
	p := 0
	for {
		q++
		ba := &[]byte{}
		rows.Scan(ba)
		if len(*ba) > 35 {
			if string((*ba)[38:41]) == "ggS" {
				magic := (*ba)[24:37]
				for _, b := range magic {
					if b == 0 {
						fmt.Print(".")
					} else {
						fmt.Print(string(b))
					}
				}
				fmt.Print(" - ")
				fmt.Printf("%x", magic)


				fmt.Println()
			}
		}
		if !rows.Next() {
			break
		}
	}
	fmt.Println("went through total", q, "cache items")
	fmt.Println("found", p, "pngs")
	
}