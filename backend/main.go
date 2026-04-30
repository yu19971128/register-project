package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"clinic/db"
	"clinic/router"
)

func main() {
	database, err := db.Open(os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer database.Close()

	if err := runMigrations(database); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	r := router.Setup(database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("start server: %v", err)
	}
}

func runMigrations(database *sql.DB) error {
	migrationsDir := "./migrations"
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		path := filepath.Join(migrationsDir, entry.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", entry.Name(), err)
		}

		// Split statements and execute individually for SQLite compatibility
		statements := splitSQLStatements(string(content))
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" || strings.HasPrefix(stmt, "--") {
				continue
			}
			if _, err := database.Exec(stmt); err != nil {
				// SQLite: ignore duplicate column/table/index errors for idempotency
				msg := strings.ToLower(err.Error())
				if strings.Contains(msg, "duplicate column name") ||
					strings.Contains(msg, "already exists") {
					continue
				}
				return fmt.Errorf("exec migration %s: %w", entry.Name(), err)
			}
		}
		log.Printf("Applied migration: %s", entry.Name())
	}
	return nil
}

func splitSQLStatements(sql string) []string {
	var statements []string
	var current strings.Builder
	lines := strings.Split(sql, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "--") {
			continue
		}
		current.WriteString(line)
		current.WriteString("\n")
		if strings.HasSuffix(trimmed, ";") {
			statements = append(statements, current.String())
			current.Reset()
		}
	}
	if current.Len() > 0 {
		statements = append(statements, current.String())
	}
	return statements
}
