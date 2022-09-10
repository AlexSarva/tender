package main

import (
	"AlexSarva/tender/admin"
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
	flag.StringVar(&cfg.DatabasePG, "dbpg", cfg.DatabasePG, "postgresql database config")
	flag.StringVar(&cfg.DatabaseClick, "dbclick", cfg.DatabaseClick, "clickhouse database config")
	flag.Parse()
	log.Printf("%+v\n", cfg)
	log.Printf("ServerAddress: %v", cfg.ServerAddress)
	DBClick, dbErr := app.NewStorage("CLICK", cfg)
	if dbErr != nil {
		log.Fatal(dbErr.Error() + "говно")
	}
	adminPG := admin.NewAdminDBConnection(cfg.DatabasePG)
	ping := DBClick.Repo.Ping()
	log.Println(ping)
	MainApp := server.NewServer(&cfg, DBClick, adminPG)
	if runErr := MainApp.Run(); runErr != nil {
		log.Printf("%s", runErr.Error())
	}
}
