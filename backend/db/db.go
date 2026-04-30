package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func Open(dsn string) (*sql.DB, error) {
	if dsn == "" {
		dsn = "./data/clinic.db"
	}
	_ = os.MkdirAll("./data", 0755)
	db, err := sql.Open("sqlite", dsn+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return db, nil
}

func ExecMigration(db *sql.DB, sql string) error {
	_, err := db.Exec(sql)
	return err
}
