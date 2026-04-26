package bootstrapper

import (
	"os"
	"path/filepath"

	"axell.me/rblxutils/common"
	_ "modernc.org/sqlite"
)

func DeleteCacheDb() {
	err := os.Remove(filepath.Join(common.RobloxAppData, "rbx-storage.db"))
	if err != nil && !os.IsNotExist(err) {
		common.FatalError(err)
	}

	/*err = os.Remove(filepath.Join(common.RobloxAppData, "rbx-storage.db-shm"))
	if err != nil && !os.IsNotExist(err) {
		common.FatalError(err)
	}

	err = os.Remove(filepath.Join(common.RobloxAppData, "rbx-storage.db-wal"))
	if err != nil && !os.IsNotExist(err) {
		common.FatalError(err)
	}

	err = os.Remove(filepath.Join(common.RobloxAppData, "rbx-storage.id"))
	if err != nil && !os.IsNotExist(err) {
		common.FatalError(err)
	}*/

	err = os.RemoveAll(filepath.Join(common.LocalAppData, "Temp", "Roblox", "http"))
	if err != nil && !os.IsNotExist(err) {
		common.FatalError(err)
	}
}


	/*
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

	_, err = db.Query("DELETE * FROM files")
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
	fmt.Println("found", p, "pngs")*/