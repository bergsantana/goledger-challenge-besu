package database

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("PG_CONN"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS storage (
			id SERIAL PRIMARY KEY,
			address TEXT NOT NULL UNIQUE,
			value INTEGER NOT NULL
		);
	`)

	if err != nil {
		log.Fatalf("Could not connect to dabatase: %v", err)
	}

	return db
}
