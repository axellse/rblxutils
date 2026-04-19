package bootstrapper

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"axell.me/rblxutils/common"
	_ "modernc.org/sqlite"
)

func OpenCacheDatabase() {
	dbFile := filepath.Join(common.RobloxAppData, "rbx-storage.db")
	fmt.Println("looking for cache db @ " + dbFile)
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		common.FatalErrorStr("cache db does not exist?")
	}
	fmt.Println("db found ok")
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		common.FatalError(err)
	}

	rows, err := db.Query("SELECT content FROM files")
	if err != nil {
		common.FatalError(err)
	}

	if rows == nil {
		fmt.Println("weird rows edgecase encountered")
	}

	q := 0
	p := 0
	for {
		q++
		ba := &[]byte{}
		rows.Scan(ba)
		if len(*ba) > 35 {
			if string((*ba)[38:41]) == "ggS" {
				//magic := (*ba)[24:37]
				p++
			}
		}
		if !rows.Next() {
			break
		}
	}
	fmt.Println("went through total", q, "cache items")
	fmt.Println("found", p, "pngs")

}
