package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(dsn string) *sqlx.DB {
	database, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	database.SetMaxOpenConns(25)
	database.SetMaxIdleConns(5)

	log.Println("database connection established")
	return database
}
