package main

import (
	"database/sql"
	"log"
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
			value INTEGER NOT NULL
		);
	`)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
