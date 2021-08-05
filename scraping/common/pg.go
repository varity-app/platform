package common

import (
	"os"

	"github.com/go-pg/pg/v10"
)

// Initialize postgres
func InitPostgres() *pg.DB {
	address := os.Getenv("POSTGRES_ADDRESS")
	database := os.Getenv("POSTGRES_DB")
	username := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	network := os.Getenv("POSTGRES_NETWORK") // Can only be `unix` or `tcp`
	if network != "unix" {
		network = "tcp"
	}

	db := pg.Connect(&pg.Options{
		Addr:     address,
		User:     username,
		Password: password,
		Database: database,
		Network:  network,
	})

	return db
}
