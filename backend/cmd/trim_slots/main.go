package main

import (
	"log"
	"math/rand"

	"clinic/db"
)

func main() {
	database, err := db.Open("")
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer database.Close()

	type key struct {
		date   string
		doctor string
	}

	rows, err := database.Query(`SELECT id, date, doctor_name FROM schedules ORDER BY id`)
	if err != nil {
		log.Fatalf("query schedules: %v", err)
	}
	groups := map[key][]int64{}
	for rows.Next() {
		var id int64
		var date, doctor string
		if err := rows.Scan(&id, &date, &doctor); err != nil {
			log.Fatalf("scan: %v", err)
		}
		k := key{date, doctor}
		groups[k] = append(groups[k], id)
	}
	rows.Close()

	delStmt, err := database.Prepare(`DELETE FROM schedules WHERE id = ?`)
	if err != nil {
		log.Fatalf("prepare delete: %v", err)
	}
	defer delStmt.Close()

	deleted := 0
	kept := 0
	for _, ids := range groups {
		if len(ids) <= 1 {
			kept++
			continue
		}
		keepIdx := rand.Intn(len(ids))
		for i, id := range ids {
			if i == keepIdx {
				continue
			}
			if _, err := delStmt.Exec(id); err != nil {
				log.Printf("delete %d: %v", id, err)
				continue
			}
			deleted++
		}
		kept++
	}

	log.Printf("Trimmed schedules: kept %d (one per doctor per day), deleted %d", kept, deleted)

	var total int
	_ = database.QueryRow("SELECT COUNT(*) FROM schedules").Scan(&total)
	log.Printf("Schedules remaining: %d", total)
}
