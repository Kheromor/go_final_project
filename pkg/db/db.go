package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

var DB *sql.DB

func Init(dbFile string) error {
	if _, err := os.Stat(dbFile); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	var err error
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if _, err = DB.Exec(schema); err != nil {
		DB.Close()
		return err
	}

	return DB.Ping()
}
