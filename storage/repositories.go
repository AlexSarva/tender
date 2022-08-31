package storage

// Repo primary interface for all types of databases
type Repo interface {
	Ping() bool
	//NewUser(user *models.User) error
	//GetUser(username string) (*models.User, error)
}
