package storagepg

import (
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// ErrDuplicatePK error that occurs when adding exists user or order number
var ErrDuplicatePK = errors.New("duplicate PK")

// ErrNoValues error that occurs when no values selected from database
var ErrNoValues = errors.New("no values from select")

// PostgresDB initializing from PostgreSQL database
type PostgresDB struct {
	database *sqlx.DB
}

// NewPostgresDBConnection initializing from PostgreSQL database connection
func NewPostgresDBConnection(config string) *PostgresDB {
	db, err := sqlx.Connect("postgres", config)
	var schemas = ddl
	db.MustExec(schemas)
	if err != nil {
		log.Println(err)
	}
	return &PostgresDB{
		database: db,
	}
}

// Ping check availability of database
func (d *PostgresDB) Ping() bool {
	//d.database.
	return d.database.Ping() == nil
}

//// NewUser insert new User in Databse
//func (d *PostgresDB) NewUser(user *models.User) error {
//	tx := d.database.MustBegin()
//	resInsert, resErr := tx.NamedExec("INSERT INTO public.users (id, username, passwd, cookie, cookie_expires) VALUES (:id, :username, :passwd, :cookie, :cookie_expires) on conflict (username) do nothing ", &user)
//	if resErr != nil {
//		return resErr
//	}
//	affectedRows, _ := resInsert.RowsAffected()
//	if affectedRows == 0 {
//		return ErrDuplicatePK
//	}
//	commitErr := tx.Commit()
//	if commitErr != nil {
//		return commitErr
//	}
//	return nil
//}
//
//// GetUser get user credentials from database by username
//func (d *PostgresDB) GetUser(username string) (*models.User, error) {
//	var user models.User
//	err := d.database.Get(&user, "SELECT id, username, passwd, cookie, cookie_expires FROM public.users WHERE username=$1", username)
//	if err != nil {
//		log.Println(err)
//		return &models.User{}, err
//	}
//	return &user, err
//}
