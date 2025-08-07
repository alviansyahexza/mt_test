package config

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func GetDbConnection() *sql.DB {
	connStr := "postgres://postgres:mysecretpassword@localhost:5432/mt_test?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
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
