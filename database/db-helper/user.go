package dbhelper

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DBWrapper struct {
	DB *sql.DB
}

func (db *DBWrapper) CreateUser(username, password string) error {

	_, err := db.DB.Exec("INSERT INTO auth (username, password) VALUES ($1, $2)", username, password)
	if err != nil {
		return err
	}
	return nil
}

func (db *DBWrapper) UserExists(username, password string) error {
	var user string
	err := db.DB.QueryRow("SELECT username FROM auth WHERE username = $1 AND password = $2", username, password).Scan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}
	return nil
}
func (db *DBWrapper) LogoutUser(username string) error {
	_, err := db.DB.Exec("DELETE FROM session WHERE username = $1", username)
	if err != nil {
		return err
	}
	return nil
}
