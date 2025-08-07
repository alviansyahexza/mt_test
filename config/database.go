package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func GetDbConnection() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(5 * 60)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Database connection established successfully")
	return db
}
