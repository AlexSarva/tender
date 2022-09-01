package storage

import "AlexSarva/tender/models"

// Repo primary interface for all types of databases
type Repo interface {
	Ping() bool
	GetOrgInfo(inn string) (models.Organization, error)
	//NewUser(user *models.User) error
	//GetUser(username string) (*models.User, error)
}
