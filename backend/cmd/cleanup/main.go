package main

import (
	"database/sql"
	"log"

	"clinic/db"
)

func main() {
	database, err := db.Open("")
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer database.Close()

	tables := []string{"orders", "patients", "schedules"}

	log.Println("=== Before cleanup ===")
	printCounts(database, tables)

	// Delete in FK-safe order: orders depend on patients & schedules.
	for _, t := range tables {
		res, err := database.Exec("DELETE FROM " + t)
		if err != nil {
			log.Fatalf("delete from %s: %v", t, err)
		}
		n, _ := res.RowsAffected()
		log.Printf("DELETE FROM %s: %d rows", t, n)
	}

	// Reset AUTOINCREMENT counters so new inserts start from id=1.
	if _, err := database.Exec(
		`DELETE FROM sqlite_sequence WHERE name IN ('orders','patients','schedules')`,
	); err != nil {
		log.Printf("reset sqlite_sequence: %v", err)
	}

	log.Println("=== After cleanup ===")
	printCounts(database, tables)
}

func printCounts(database *sql.DB, tables []string) {
	for _, t := range tables {
		var n int
		if err := database.QueryRow("SELECT COUNT(*) FROM " + t).Scan(&n); err != nil {
			log.Printf("count %s: %v", t, err)
			continue
		}
		log.Printf("  %s: %d", t, n)
	}
}
