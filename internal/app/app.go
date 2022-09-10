package app

import (
	"AlexSarva/tender/models"
	"AlexSarva/tender/storage"
	clickhousestorage "AlexSarva/tender/storage/storageclick"
	"errors"
	"fmt"
)

// Database interface for different types of databases
type Database struct {
	Repo storage.Repo
}

// NewStorage generate new instance of database
func NewStorage(dbName string, cfg models.Config) (*Database, error) {
	if dbName == "CLICK" {
		DB := clickhousestorage.MyClickHouseDB(cfg.DatabaseClick)
		fmt.Println("Using ClickHouse Database")
		return &Database{
			Repo: DB,
		}, nil
	} else {
		return &Database{}, errors.New("u must use database config")
	}

}
