package main

import (
	"log"
	"math/rand"
	"time"

	"clinic/db"
)

func main() {
	database, err := db.Open("")
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer database.Close()

	departments := []struct {
		dept   string
		doctor string
	}{
		{"泌尿外科", "张医生"},
		{"泌尿外科", "李医生"},
		{"内科", "王医生"},
		{"内科", "刘医生"},
		{"外科", "陈医生"},
		{"儿科", "赵医生"},
		{"妇科", "孙医生"},
		{"眼科", "周医生"},
		{"口腔科", "吴医生"},
		{"皮肤科", "郑医生"},
	}

	slots := []struct {
		start string
		end   string
	}{
		{"08:00", "12:00"}, // 上午
		{"12:00", "18:00"}, // 下午
		{"18:00", "20:00"}, // 晚上
	}

	// 2026-05-05 to 2026-05-30 inclusive
	startDate := time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 5, 30, 0, 0, 0, 0, time.UTC)

	stmt, err := database.Prepare(`
		INSERT INTO schedules (date, department, doctor_name, start_time, end_time, total_quota, remaining, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, 'available')
	`)
	if err != nil {
		log.Fatalf("prepare: %v", err)
	}
	defer stmt.Close()

	count := 0
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		for _, dept := range departments {
			for _, slot := range slots {
				quota := 10 + rand.Intn(21) // [10, 30]
				_, err := stmt.Exec(dateStr, dept.dept, dept.doctor, slot.start, slot.end, quota, quota)
				if err != nil {
					log.Printf("skip %s %s %s %s-%s: %v", dateStr, dept.dept, dept.doctor, slot.start, slot.end, err)
					continue
				}
				count++
			}
		}
	}

	log.Printf("Seeded %d schedules from %s to %s", count, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
}
