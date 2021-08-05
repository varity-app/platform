package main

import (
	"os"

	"github.com/go-pg/pg/v10"
)

// Initialize postgres
func initPostgres() *pg.DB {
	hostname := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	database := os.Getenv("POSTGRES_DB")
	username := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")

	db := pg.Connect(&pg.Options{
		Addr:     hostname + ":" + port,
		User:     username,
		Password: password,
		Database: database,
	})

	return db
}
