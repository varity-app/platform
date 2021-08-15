package main

import (
	"log"
	"os"

	"github.com/VarityPlatform/scraping/data/tickers"

	"github.com/go-pg/pg/v10"
	"github.com/spf13/viper"

	"github.com/labstack/echo/v4"
)

func initPostgres() *pg.DB {
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

// Entrypoint method
func main() {
	if err := initConfig(); err != nil {
		log.Fatal(err)
	}
	db := initPostgres()

	// Initialize ticker repo
	tickerRepo := tickers.NewIEXTickerRepo(db)
	tickers, err := tickerRepo.List()
	if err != nil {
		log.Fatalf("tickerRepo.List: %v", err)
	}
	tickerRepo.Close()

	// Initialize webserver
	web := echo.New()
	web.HideBanner = true
	setupRoutes(web, tickers)

	// Start webserver
	port := viper.GetString("port")
	web.Logger.Fatal(web.Start(":" + port))
}
