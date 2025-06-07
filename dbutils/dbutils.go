package dbutils

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func ConnectSqlite(dbPath string) (*sql.DB, func()) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("error opening db: %s \n", err.Error())
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("error connecting to the database: %v \n", err)
	}

	dbClose := func() {
		db.Close()
	}

	return db, dbClose
}

func CleanUpSqliteDbFile(dbPath string) {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Printf("Database file does not exist: %s\n", dbPath)
		return
	}

	if err := os.Remove(dbPath); err != nil {
		log.Printf("failed to delete SQLite database at %s: %w", dbPath, err)
		return
	}

	log.Printf("Successfully deleted SQLite database: %s\n", dbPath)
	return
}
