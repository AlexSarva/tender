package main

import (
	"AlexSarva/tender/internal/app"
	"AlexSarva/tender/models"
	"AlexSarva/tender/server"
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

func main() {
	var cfg models.Config
	// Priority on flags
	// Load config from env
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	// Rewrite from start parameters
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "host:port to listen on")
	flag.StringVar(&cfg.Database, "d", cfg.Database, "database config")
	flag.Parse()
	log.Printf("%+v\n", cfg)
	log.Printf("ServerAddress: %v", cfg.ServerAddress)
	DB, dbErr := app.NewStorage(cfg.Database)
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	ping := DB.Repo.Ping()
	log.Println(ping)
	MainApp := server.NewServer(&cfg, DB)
	if runErr := MainApp.Run(); runErr != nil {
		log.Printf("%s", runErr.Error())
	}
}
