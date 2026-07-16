//go:build ignore

package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT id,date,title,repeat FROM scheduler ORDER BY id DESC LIMIT 8")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var date, title, repeat string
		err = rows.Scan(&id, &date, &title, &repeat)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%d %s %s %s\n", id, date, title, repeat)
	}
}
